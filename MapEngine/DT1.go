package MapEngine

import (
	"log"

	"github.com/essial/OpenDiablo2/Common"
)

// https://d2mods.info/forum/viewtopic.php?t=65163

type Block struct {
	X          int16
	Y          int16
	GridX      byte
	GridY      byte
	Format     int16
	Length     int32
	FileOffset int32
}

type Tile struct {
	Direction          int32
	RoofHeight         int16
	SoundIndex         byte
	Animated           byte
	Height             int32
	Width              int32
	Orientation        int32
	MainIndex          int32
	SubIndex           int32
	RarityFrameIndex   int32
	SubTileFlags       [25]byte
	blockHeaderPointer int32
	blockHeaderSize    int32
	blocks             []Block
}

type DT1 struct {
	Tiles []Tile
}

func LoadDT1(path string, fileProvider Common.FileProvider) *DT1 {
	result := &DT1{}
	fileData := fileProvider.LoadFile(path)
	br := Common.CreateStreamReader(fileData)
	ver1 := br.GetInt32()
	ver2 := br.GetInt32()
	if ver1 != 7 || ver2 != 6 {
		log.Panicf("Expected %s to have a version of 7.6, but got %d.%d instead", path, ver1, ver2)
	}
	br.SkipBytes(260)
	numberOfTiles := br.GetInt32()
	br.SkipBytes(4)
	result.Tiles = make([]Tile, numberOfTiles)
	for tileIdx := range result.Tiles {
		newTile := Tile{}
		newTile.Direction = br.GetInt32()
		newTile.RoofHeight = br.GetInt16()
		newTile.SoundIndex = br.GetByte()
		newTile.Animated = br.GetByte()
		newTile.Height = br.GetInt32()
		newTile.Width = br.GetInt32()
		br.SkipBytes(4)
		newTile.Orientation = br.GetInt32()
		newTile.MainIndex = br.GetInt32()
		newTile.SubIndex = br.GetInt32()
		newTile.RarityFrameIndex = br.GetInt32()
		br.SkipBytes(4)
		for i := range newTile.SubTileFlags {
			newTile.SubTileFlags[i] = br.GetByte()
		}
		br.SkipBytes(7)
		newTile.blockHeaderPointer = br.GetInt32()
		newTile.blockHeaderSize = br.GetInt32()
		newTile.blocks = make([]Block, br.GetInt32())
		br.SkipBytes(12)
		result.Tiles[tileIdx] = newTile
	}
	for tileIdx, tile := range result.Tiles {
		br.SetPosition(uint64(tile.blockHeaderPointer))
		for blockIdx := range tile.blocks {
			result.Tiles[tileIdx].blocks[blockIdx].X = br.GetInt16()
			result.Tiles[tileIdx].blocks[blockIdx].Y = br.GetInt16()
			br.SkipBytes(2)
			result.Tiles[tileIdx].blocks[blockIdx].GridX = br.GetByte()
			result.Tiles[tileIdx].blocks[blockIdx].GridY = br.GetByte()
			result.Tiles[tileIdx].blocks[blockIdx].Format = br.GetInt16()
			result.Tiles[tileIdx].blocks[blockIdx].Length = br.GetInt32()
			br.SkipBytes(2)
			result.Tiles[tileIdx].blocks[blockIdx].FileOffset = br.GetInt32()
		}
		/*
			for blockIdx, block := range tile.blocks {
				br.SetPosition(uint64(tile.blockHeaderPointer + block.FileOffset))
				encodedData, _ := br.ReadBytes(block.Length)
				bs := Common.CreateBitStream(encodedData)
			}
		*/
	}
	return result
}
