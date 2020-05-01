/*
 *          Copyright 2020, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
	"strconv"
)

func init() {
    runtime.LockOSThread()
}

func main() {
	glfw.SetErrorCallback(onError)

	err := glfw.Init()

	if err == nil {
		defer glfw.Terminate()

		var window *glfw.Window
		window, err = glfw.CreateWindow(480, 360, "OpenGL Example", nil, nil)

		if err == nil {
			defer window.Destroy()

			window.SetKeyCallback(onKey)
			window.SetSizeCallback(onResize)
			window.MakeContextCurrent()

			for !window.ShouldClose() {
				display()
				window.SwapBuffers()
				glfw.PollEvents()
			}
		} else {
			printError(2, err)
		}
	} else {
		printError(1, err)
	}
}

func onError(err glfw.ErrorCode, desc string) {
	errStr := strconv.Itoa(int(err))
	printError(3, errStr + " " + desc)
}

func onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func display() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	shaderProgram.Use()
	defer gl.ProgramUnuse()

	posBuffer.Bind(gl.ARRAY_BUFFER)

	positionAttrib := gl.AttribLocation(shaderProgram.GetAttribLocation("position"))
	positionAttrib.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))
	positionAttrib.EnableArray()
	defer positionAttrib.DisableArray()

	colorAttrib := gl.AttribLocation(shaderProgram.GetAttribLocation("color"))
	colorAttrib.AttribPointer(4, gl.FLOAT, false, 0, uintptr((len(vertices)*float32_size)/2))
	colorAttrib.EnableArray()
	defer colorAttrib.DisableArray()

	gl.DrawArrays(gl.TRIANGLES, 0, len(vertices)/2/float32_size)
}
