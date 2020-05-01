/*
 *          Copyright 2020, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import "fmt"

func printError(id int, err error) {
	fmt.Printf("error (%d): %s", id, err.Error())
}
