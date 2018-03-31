package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var le = binary.LittleEndian

func decoder_xor(b []byte) {
	var key uint32 = 0xa340f119
	for i := 0; i < 0x400; i += 4 {
		old_value := le.Uint32(b[i:])
		le.PutUint32(b[i:], old_value^key)
		key = old_value
	}
}

func decoder_swap(b []byte) {
	index1 := []int{0x18, 0xa8, 0xc0, 0xe4, 0x16c, 0x190, 0x198, 0x1c8}
	index2 := []int{0x278, 0x258, 0x25c, 0x338, 0x318, 0x33c, 0x3fc, 0x394}
	for i, _ := range index1 {
		temp := le.Uint32(b[index1[i]:])
		le.PutUint32(b[index1[i]:], le.Uint32(b[index2[i]:]))
		le.PutUint32(b[index2[i]:], temp)
	}
}

func decoder_reverse(b []byte) {
	for i, j := 0, 255; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func main() {
	// process command line
	if len(os.Args) == 1 || len(os.Args) > 3 {
		fmt.Println("Usage:", os.Args[0], "inputfile [outputfile]")
		os.Exit(1)
	}
	inputfn := os.Args[1]
	var outputfn string
	if len(os.Args) == 3 {
		outputfn = os.Args[2]
	}

	// read relevant parts of input file
	f, err := os.Open(inputfn)
	if err != nil {
		fmt.Printf("Error opening %q: %v\n", inputfn, err)
		os.Exit(2)
	}
	filesize, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		fmt.Printf("Seek error: %v\n", err)
		os.Exit(2)
	}
	if filesize < 0x300400 {
		fmt.Printf("File is too small!\n")
		os.Exit(2)
	}
	b := make([]byte, 17)
	f.ReadAt(b, filesize-17)
	if string(b) != "dkaghghkehlsvkdlf" {
		fmt.Printf("No magic bytes at the end of the file!\n")
		os.Exit(2)
	}
	b = make([]byte, 0x300000)
	_, err = f.ReadAt(b, filesize-0x300400)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		os.Exit(2)
	}
	f.Close()

	// decode encoded area
	decoders := []func(b []byte){decoder_xor, decoder_swap, decoder_reverse}
	for i := 0; i < 0x300000/0x400; i++ {
		decoders[i%3](b[i*0x400 : (i+1)*0x400])
	}

	// write back file
	if outputfn != "" {
		f, err = os.Open(inputfn)
		if err != nil {
			fmt.Printf("Error opening %q: %v\n", inputfn, err)
			os.Exit(2)
		}
		g, err := os.Create(outputfn)
		if err != nil {
			fmt.Printf("Error creating %q: %v\n", outputfn, err)
			os.Exit(2)
		}
		_, err = io.CopyN(g, f, filesize-0x300400)
		if err != nil {
			fmt.Printf("Copy error: %v\n", err)
			os.Exit(2)
		}
		f.Close()
		_, err = g.Write(b)
		if err != nil {
			fmt.Printf("Write error: %v\n", err)
			os.Exit(2)
		}
		g.Close()
	} else { // in-place fix
		f, err := os.OpenFile(inputfn, os.O_RDWR, 0644)
		if err != nil {
			fmt.Printf("Error opening %q: %v\n", inputfn, err)
			os.Exit(2)
		}
		_, err = f.WriteAt(b, filesize-0x300400)
		if err != nil {
			fmt.Printf("Write error: %v\n", err)
			os.Exit(2)
		}
		err = f.Truncate(filesize - 0x400)
		if err != nil {
			fmt.Printf("Truncate error: %v\n", err)
			os.Exit(2)
		}
		f.Close()
	}
}
