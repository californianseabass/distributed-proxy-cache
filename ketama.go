package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"

	"github.com/ghetzel/shmtool/shm"
)

var nIntervalSplits int = 40
var replicationFactor int = 4

type Server struct {
	Weight int
	addr   string
}

type Ring []uint32

func (r Ring) Len() int           { return len(r) }
func (r Ring) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Ring) Less(i, j int) bool { return r[i] < r[j] }

type Ketama struct {
	Ring             Ring
	hashToServerAddr map[uint32]string
}

func (k Ketama) getServerForString(key string) string {
	hash := ketamaHash(key)
	i := sort.Search(len(k.Ring), func(i int) bool { return k.Ring[i] >= hash })
	serverHash := k.Ring[i]
	return k.hashToServerAddr[serverHash]
}

func md5Sum(in string) [16]byte {
	return md5.Sum([]byte(in))
}

func ketamaHash(in string) uint32 {
	hash := md5.Sum([]byte(in))

	firstFourBytes := [4]byte{}
	firstFourBytes[0] = hash[0]
	firstFourBytes[1] = hash[1]
	firstFourBytes[2] = hash[2]
	firstFourBytes[3] = hash[3]

	return binary.BigEndian.Uint32(firstFourBytes[:])
}

func createRingFromServers(servers []Server) Ketama {

	totalMemory := 0
	for _, s := range servers {
		totalMemory += s.Weight
	}

	nSlices := nIntervalSplits * replicationFactor * totalMemory

	ring := make([]uint32, nSlices)
	hashToServerAddr := map[uint32]string{}

	count := 0
	for _, s := range servers {
		ks := nIntervalSplits * replicationFactor * s.Weight

		for k := 0; k < ks; k++ {
			name := fmt.Sprintf("%s-%d", s.addr, k)
			hash := ketamaHash(name)
			ring = append(ring, hash)
			if val, ok := hashToServerAddr[hash]; ok {
				fmt.Println("hash collision initializing ring: ", val, " ", name)
			}
			hashToServerAddr[hash] = s.addr
			count++
		}
	}

	sort.Slice(ring, func(i, j int) bool { return ring[i] < ring[j] })

	return Ketama{ring, hashToServerAddr}
}

func addRingToSharedMemory(ring Ring) int {
	size := len(ring) * int(reflect.TypeOf(uint32(0)).Size())
	segment, err := shm.Create(size)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, size)
	for i, v := range ring {
		binary.BigEndian.PutUint32(buf[i*4:], v)
	}
	_, err = segment.Write(buf)
	if err != nil {
		panic(err)
	}
	return segment.Id
}

func persist(filename string, ketama Ketama) int {
	sharedMemId := addRingToSharedMemory(ketama.Ring)
	jsonString, err := json.Marshal(ketama.hashToServerAddr)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(filename, []byte(jsonString), 0644)
	return sharedMemId
}

func funcLoadRingFromSharedMemory(shmKey int) Ring {
	segment, err := shm.Open(shmKey)
	segment.Reset()
	if err != nil {
		panic(err)
	}
	var ring []uint32 = make([]uint32, int(segment.Size)/int(reflect.TypeOf(uint32(0)).Size()))
	data, err := ioutil.ReadAll(segment)
	if err != nil {
		panic(err)
	}
	binary.Read(bytes.NewReader(data), binary.BigEndian, ring)
	return ring
}

func recreateKetama(filename string, shmKey int) Ketama {
	ring := funcLoadRingFromSharedMemory(shmKey)
	var hashToServerAddr map[uint32]string
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &hashToServerAddr)
	if err != nil {
		panic(err)
	}
	return Ketama{ring, hashToServerAddr}
}

func main() {
	serverOne := Server{
		addr:   "10.0.1.1:11211",
		Weight: 300,
	}
	serverTwo := Server{
		addr:   "10.0.1.2:11211",
		Weight: 300,
	}
	serverThree := Server{
		addr:   "10.0.1.3:11211",
		Weight: 300,
	}
	serverFour := Server{
		addr:   "10.0.1.4:11211",
		Weight: 300,
	}
	serverFive := Server{
		addr:   "10.0.1.5:11211",
		Weight: 300,
	}
	serverSix := Server{
		addr:   "10.0.1.6:11211",
		Weight: 300,
	}
	serverSeven := Server{
		addr:   "10.0.1.7:11211",
		Weight: 300,
	}
	serverEight := Server{
		addr:   "10.0.1.8:11211",
		Weight: 300,
	}
	servers := []Server{serverOne, serverTwo, serverThree, serverFour, serverFive, serverSix, serverSeven, serverEight}
	ketama := createRingFromServers(servers)
	server := ketama.getServerForString("www.espn.com")
	fmt.Println("server: ", server)
	segmentId := persist("test.json", ketama)
	ketama = recreateKetama("test.json", segmentId)
	server = ketama.getServerForString("www.espn.com")
	fmt.Println("server: ", server)
}

func getServers() []Server {
	serverOne := Server{
		addr:   "10.0.1.1:11211",
		Weight: 300,
	}
	serverTwo := Server{
		addr:   "10.0.1.2:11211",
		Weight: 300,
	}
	serverThree := Server{
		addr:   "10.0.1.3:11211",
		Weight: 300,
	}
	serverFour := Server{
		addr:   "10.0.1.4:11211",
		Weight: 300,
	}
	serverFive := Server{
		addr:   "10.0.1.5:11211",
		Weight: 300,
	}
	serverSix := Server{
		addr:   "10.0.1.6:11211",
		Weight: 300,
	}
	serverSeven := Server{
		addr:   "10.0.1.7:11211",
		Weight: 300,
	}
	serverEight := Server{
		addr:   "10.0.1.8:11211",
		Weight: 300,
	}
	return []Server{serverOne, serverTwo, serverThree, serverFour, serverFive, serverSix, serverSeven, serverEight}
}
