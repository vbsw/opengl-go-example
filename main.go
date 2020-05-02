/*
 *          Copyright 2020, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/vbsw/plainshader"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()

	if err == nil {
		var window *glfw.Window
		defer glfw.Terminate()
		window, err = glfw.CreateWindow(300, 300, "OpenGL Example", nil, nil)

		if err == nil {
			defer window.Destroy()
			window.SetKeyCallback(onKey)
			window.SetSizeCallback(onResize)
			window.MakeContextCurrent()
			err = gl.Init()

			if err == nil {
				var vertexShader uint32
				vertexShader, err = loadShader(gl.VERTEX_SHADER, plainshader.VertexShader)

				if err == nil {
					var fragmentShader uint32
					fragmentShader, err = loadShader(gl.FRAGMENT_SHADER, plainshader.FragmentShader)

					if err == nil {
						var program uint32
						program, err = loadProgram(vertexShader, fragmentShader)

						if err == nil {
							defer gl.DeleteProgram(program)
							triangleVBOs := createTriangle()
							defer gl.DeleteBuffers(2, &triangleVBOs[0])
							triangleVAO := createVAO(triangleVBOs)
							defer gl.DeleteVertexArrays(1, &triangleVAO)
							gl.UseProgram(program)
							// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

							for !window.ShouldClose() {
								gl.ClearColor(0, 0, 0, 0)
								gl.Clear(gl.COLOR_BUFFER_BIT)
								gl.DrawArrays(gl.TRIANGLES, 0, 3)

								window.SwapBuffers()
								glfw.PollEvents()
							}
						}
					}
				}
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func onResize(w *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func loadShader(shaderType uint32, shaderSource **uint8) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, shaderSource, nil)
	gl.CompileShader(shader)
	err := checkShader(shader, gl.COMPILE_STATUS)

	if err != nil {
		gl.DeleteShader(shader)
	}
	return shader, err
}

func loadProgram(vShader, fShader uint32) (uint32, error) {
	program := gl.CreateProgram()
	gl.AttachShader(program, vShader)
	gl.AttachShader(program, fShader)
	gl.LinkProgram(program)
	err := checkProgram(program, gl.LINK_STATUS)

	if err == nil {
		gl.ValidateProgram(program)
		err = checkProgram(program, gl.VALIDATE_STATUS)
	}
	if err != nil {
		gl.DeleteProgram(program)
	}
	return program, err
}

func checkShader(shader, statusType uint32) error {
	var status int32
	var err error

	gl.GetShaderiv(shader, statusType, &status)

	if status == gl.FALSE {
		var length int32
		var infoLog string

		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)

		if length > 0 {
			infoLogBytes := make([]byte, length)
			gl.GetShaderInfoLog(shader, length, nil, &infoLogBytes[0])
			infoLog = string(infoLogBytes)
		}
		err = errors.New("shader " + infoLog)
	}
	return err
}

func checkProgram(program, statusType uint32) error {
	var status int32
	var err error

	gl.GetProgramiv(program, statusType, &status)

	if status == gl.FALSE {
		var length int32
		var infoLog string

		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)

		if length > 0 {
			infoLogBytes := make([]byte, length)
			gl.GetProgramInfoLog(program, length, nil, &infoLogBytes[0])
			infoLog = string(infoLogBytes)
		}
		err = errors.New("program " + infoLog)
	}
	return err
}

func createTriangle() []uint32 {
	points := []float32{
		0.0, 1.0, 0.0,
		1.0, -1.0, 0.0,
		-1.0, -1.0, 0.0,
	}
	colours := []float32{
		1.0, 0.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 1.0,
	}

	buffers := make([]uint32, 2)
	gl.GenBuffers(2, &buffers[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(colours)*4, gl.Ptr(colours), gl.STATIC_DRAW)

	return buffers
}

func createVAO(buffers []uint32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[0])
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[1])
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)

	return vao
}
