package main

import (
    "io/ioutil"
    "fmt"
    "math/rand"
    // "encoding/hex"
)

func operation(opcode int, cpu *cpu, memory *[4096]byte, gfx *[32*64]byte) {
  switch opcode & 0xF000 {
  case 0xA000:
    // Move NNN to I
    fmt.Println("Move NNN to I")
    cpu.I = (opcode & 0x0FFF)
    cpu.pc += 2

  case 0xC000:
    fmt.Println("set NN into VX")
    index := (opcode & 0x0F00) >> 8
    nn := (opcode & 0x00FF) & rand.Intn(100)

    fmt.Println(nn)
    cpu.V[index] = byte(nn)
    cpu.pc += 2

  case 0x2000:
    fmt.Println("JUMPS to SR NNN")
    cpu.stack_pointers += 1
    // cpu.stack[cpu.stack_pointers] = int(pc >> 8)
    new_pc := (opcode & 0x0FFF) >> 8
    cpu.pc = new_pc

  case 0x7000:
    fmt.Println("adds NN into VX")
    index := (opcode & 0x0F00) >> 8
    cpu.V[index] += byte((opcode & 0x00FF))
    cpu.pc += 2

  case 0x6000:
    fmt.Println("sets NN into VX")
    index := (opcode & 0x0F00) >> 8
    nn := (opcode & 0x00FF)
    cpu.V[index] = byte(nn)
    cpu.pc += 2

  case 0x8000:
    fmt.Println("Dentro de 8xxx")
    num := (opcode & 0x000F)

    switch num {
    case 0:
      fmt.Println("En caso 0")

      cpu.V[(opcode & 0x0F00) >> 8] = cpu.V[(opcode & 0x00F0) >> 4]
    }

    cpu.pc += 2

  case 0x1000:
    fmt.Println("JUMP")
    cpu.pc = opcode & 0x0FFF

  case 0x3000:
    fmt.Println("if NN == VX skip next")
    index := (opcode & 0x0F00) >> 8
    nn := (opcode & 0x00FF)
    if cpu.V[index] == byte(nn) {
      cpu.pc += 4
    }else{
      cpu.pc += 2
    }
  case 0xD000:
    fmt.Println("drawing")
    x := cpu.V[((opcode & 0x0F00) >> 8)]
    y := cpu.V[((opcode & 0x00F0) >> 4)]
    h := (opcode & 0x000F)

    cpu.V[0xF] = 0

    for yline := 0; yline < h; yline++ {
      pixel := memory[cpu.I + yline]

      for xline := 0; xline < 8; xline++ {
        if((pixel & (0x80 >> uint(xline))) != 0) {
          lo := int(x) + xline + ((int(y)+yline)*64)
          if gfx[lo] == 1 {
            cpu.V[0xF] = 1
          }else{
            gfx[lo] ^= 1
          }
        }
      }
    }

    cpu.pc += 2
  default:
    fmt.Printf("No implemented %x" , opcode)
    cpu.pc += 2
  }
}

type cpu struct {
  I int
  pc int
  V [16]byte
  stack [16]byte
  stack_pointers int
}

func NewCpu() *cpu{
  return &cpu{ pc: 0x200 }
}

func debugRender(gfx *[32*64]byte) {
  for y := 1; y < 32; y++ {
    for x := 1; x < 64; x++ {
      if gfx[(y*64)+x] == 0{
        fmt.Printf(" ")
      }else{
        fmt.Printf("0")
      }
      fmt.Printf("")
    }
      fmt.Println("")
  }
}

func main() {
  var gfx [32*64]byte
  var memory [4096]byte

  b, err := ioutil.ReadFile("MAZE")

  if err != nil {
    panic(err)
  }

  cpu := NewCpu()

  for i, nim := range b {
    memory[cpu.pc+i] = nim
  }

  for i := 0; i < 1000; i++ {
    opcode := int(memory[cpu.pc]) << 8 | int(memory[cpu.pc + 1]);
    operation(opcode, cpu, &memory, &gfx)
  }

  debugRender(&gfx)
  // fmt.Printf("%b", gfx)
}
