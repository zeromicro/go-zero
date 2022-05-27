package api

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/zeromicro/antlr"
)

// Suppress unused import errors
var (
	_ = fmt.Printf
	_ = reflect.Copy
	_ = strconv.Itoa
)

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 26, 349,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 3, 2, 7,
	2, 78, 10, 2, 12, 2, 14, 2, 81, 11, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5,
	3, 88, 10, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 5, 5, 98,
	10, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3, 7, 6, 7, 108, 10, 7,
	13, 7, 14, 7, 109, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 10, 3,
	10, 3, 10, 3, 10, 6, 10, 123, 10, 10, 13, 10, 14, 10, 124, 3, 10, 3, 10,
	3, 11, 3, 11, 5, 11, 131, 10, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 13, 3,
	13, 3, 13, 3, 13, 7, 13, 141, 10, 13, 12, 13, 14, 13, 144, 11, 13, 3, 13,
	3, 13, 3, 14, 3, 14, 5, 14, 150, 10, 14, 3, 15, 3, 15, 5, 15, 154, 10,
	15, 3, 16, 3, 16, 3, 16, 5, 16, 159, 10, 16, 3, 16, 3, 16, 7, 16, 163,
	10, 16, 12, 16, 14, 16, 166, 11, 16, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17,
	5, 17, 173, 10, 17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 5, 18, 180, 10,
	18, 3, 18, 3, 18, 7, 18, 184, 10, 18, 12, 18, 14, 18, 187, 11, 18, 3, 18,
	3, 18, 3, 19, 3, 19, 3, 19, 5, 19, 194, 10, 19, 3, 19, 3, 19, 3, 20, 3,
	20, 3, 20, 5, 20, 201, 10, 20, 3, 21, 3, 21, 3, 21, 3, 21, 5, 21, 207,
	10, 21, 3, 22, 5, 22, 210, 10, 22, 3, 22, 3, 22, 3, 23, 3, 23, 3, 23, 3,
	23, 3, 23, 3, 23, 3, 23, 3, 23, 5, 23, 222, 10, 23, 3, 24, 3, 24, 3, 24,
	3, 24, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 26, 3,
	26, 3, 26, 3, 26, 3, 27, 5, 27, 241, 10, 27, 3, 27, 3, 27, 3, 28, 3, 28,
	3, 28, 6, 28, 248, 10, 28, 13, 28, 14, 28, 249, 3, 28, 3, 28, 3, 29, 3,
	29, 3, 29, 3, 29, 3, 29, 7, 29, 259, 10, 29, 12, 29, 14, 29, 262, 11, 29,
	3, 29, 3, 29, 3, 30, 5, 30, 267, 10, 30, 3, 30, 3, 30, 5, 30, 271, 10,
	30, 3, 30, 3, 30, 3, 31, 3, 31, 5, 31, 277, 10, 31, 3, 31, 6, 31, 280,
	10, 31, 13, 31, 14, 31, 281, 3, 31, 5, 31, 285, 10, 31, 3, 31, 5, 31, 288,
	10, 31, 3, 32, 3, 32, 3, 32, 3, 33, 3, 33, 3, 33, 3, 33, 5, 33, 297, 10,
	33, 3, 33, 5, 33, 300, 10, 33, 3, 34, 3, 34, 5, 34, 304, 10, 34, 3, 34,
	3, 34, 3, 35, 3, 35, 3, 35, 5, 35, 311, 10, 35, 3, 35, 3, 35, 3, 36, 3,
	36, 3, 36, 3, 36, 3, 37, 3, 37, 5, 37, 321, 10, 37, 6, 37, 323, 10, 37,
	13, 37, 14, 37, 324, 3, 38, 3, 38, 3, 38, 3, 38, 7, 38, 331, 10, 38, 12,
	38, 14, 38, 334, 11, 38, 3, 38, 3, 38, 3, 38, 3, 38, 5, 38, 340, 10, 38,
	6, 38, 342, 10, 38, 13, 38, 14, 38, 343, 3, 38, 5, 38, 347, 10, 38, 3,
	38, 2, 2, 39, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32,
	34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68,
	70, 72, 74, 2, 2, 2, 358, 2, 79, 3, 2, 2, 2, 4, 87, 3, 2, 2, 2, 6, 89,
	3, 2, 2, 2, 8, 97, 3, 2, 2, 2, 10, 99, 3, 2, 2, 2, 12, 103, 3, 2, 2, 2,
	14, 113, 3, 2, 2, 2, 16, 115, 3, 2, 2, 2, 18, 118, 3, 2, 2, 2, 20, 130,
	3, 2, 2, 2, 22, 132, 3, 2, 2, 2, 24, 136, 3, 2, 2, 2, 26, 149, 3, 2, 2,
	2, 28, 153, 3, 2, 2, 2, 30, 155, 3, 2, 2, 2, 32, 169, 3, 2, 2, 2, 34, 176,
	3, 2, 2, 2, 36, 190, 3, 2, 2, 2, 38, 200, 3, 2, 2, 2, 40, 202, 3, 2, 2,
	2, 42, 209, 3, 2, 2, 2, 44, 221, 3, 2, 2, 2, 46, 223, 3, 2, 2, 2, 48, 227,
	3, 2, 2, 2, 50, 235, 3, 2, 2, 2, 52, 240, 3, 2, 2, 2, 54, 244, 3, 2, 2,
	2, 56, 253, 3, 2, 2, 2, 58, 266, 3, 2, 2, 2, 60, 274, 3, 2, 2, 2, 62, 289,
	3, 2, 2, 2, 64, 292, 3, 2, 2, 2, 66, 301, 3, 2, 2, 2, 68, 307, 3, 2, 2,
	2, 70, 314, 3, 2, 2, 2, 72, 322, 3, 2, 2, 2, 74, 346, 3, 2, 2, 2, 76, 78,
	5, 4, 3, 2, 77, 76, 3, 2, 2, 2, 78, 81, 3, 2, 2, 2, 79, 77, 3, 2, 2, 2,
	79, 80, 3, 2, 2, 2, 80, 3, 3, 2, 2, 2, 81, 79, 3, 2, 2, 2, 82, 88, 5, 6,
	4, 2, 83, 88, 5, 8, 5, 2, 84, 88, 5, 18, 10, 2, 85, 88, 5, 20, 11, 2, 86,
	88, 5, 52, 27, 2, 87, 82, 3, 2, 2, 2, 87, 83, 3, 2, 2, 2, 87, 84, 3, 2,
	2, 2, 87, 85, 3, 2, 2, 2, 87, 86, 3, 2, 2, 2, 88, 5, 3, 2, 2, 2, 89, 90,
	8, 4, 1, 2, 90, 91, 7, 26, 2, 2, 91, 92, 7, 3, 2, 2, 92, 93, 8, 4, 1, 2,
	93, 94, 7, 23, 2, 2, 94, 7, 3, 2, 2, 2, 95, 98, 5, 10, 6, 2, 96, 98, 5,
	12, 7, 2, 97, 95, 3, 2, 2, 2, 97, 96, 3, 2, 2, 2, 98, 9, 3, 2, 2, 2, 99,
	100, 8, 6, 1, 2, 100, 101, 7, 26, 2, 2, 101, 102, 5, 16, 9, 2, 102, 11,
	3, 2, 2, 2, 103, 104, 8, 7, 1, 2, 104, 105, 7, 26, 2, 2, 105, 107, 7, 4,
	2, 2, 106, 108, 5, 14, 8, 2, 107, 106, 3, 2, 2, 2, 108, 109, 3, 2, 2, 2,
	109, 107, 3, 2, 2, 2, 109, 110, 3, 2, 2, 2, 110, 111, 3, 2, 2, 2, 111,
	112, 7, 5, 2, 2, 112, 13, 3, 2, 2, 2, 113, 114, 5, 16, 9, 2, 114, 15, 3,
	2, 2, 2, 115, 116, 8, 9, 1, 2, 116, 117, 7, 23, 2, 2, 117, 17, 3, 2, 2,
	2, 118, 119, 8, 10, 1, 2, 119, 120, 7, 26, 2, 2, 120, 122, 7, 4, 2, 2,
	121, 123, 5, 70, 36, 2, 122, 121, 3, 2, 2, 2, 123, 124, 3, 2, 2, 2, 124,
	122, 3, 2, 2, 2, 124, 125, 3, 2, 2, 2, 125, 126, 3, 2, 2, 2, 126, 127,
	7, 5, 2, 2, 127, 19, 3, 2, 2, 2, 128, 131, 5, 22, 12, 2, 129, 131, 5, 24,
	13, 2, 130, 128, 3, 2, 2, 2, 130, 129, 3, 2, 2, 2, 131, 21, 3, 2, 2, 2,
	132, 133, 8, 12, 1, 2, 133, 134, 7, 26, 2, 2, 134, 135, 5, 26, 14, 2, 135,
	23, 3, 2, 2, 2, 136, 137, 8, 13, 1, 2, 137, 138, 7, 26, 2, 2, 138, 142,
	7, 4, 2, 2, 139, 141, 5, 28, 15, 2, 140, 139, 3, 2, 2, 2, 141, 144, 3,
	2, 2, 2, 142, 140, 3, 2, 2, 2, 142, 143, 3, 2, 2, 2, 143, 145, 3, 2, 2,
	2, 144, 142, 3, 2, 2, 2, 145, 146, 7, 5, 2, 2, 146, 25, 3, 2, 2, 2, 147,
	150, 5, 30, 16, 2, 148, 150, 5, 32, 17, 2, 149, 147, 3, 2, 2, 2, 149, 148,
	3, 2, 2, 2, 150, 27, 3, 2, 2, 2, 151, 154, 5, 34, 18, 2, 152, 154, 5, 36,
	19, 2, 153, 151, 3, 2, 2, 2, 153, 152, 3, 2, 2, 2, 154, 29, 3, 2, 2, 2,
	155, 156, 8, 16, 1, 2, 156, 158, 7, 26, 2, 2, 157, 159, 7, 26, 2, 2, 158,
	157, 3, 2, 2, 2, 158, 159, 3, 2, 2, 2, 159, 160, 3, 2, 2, 2, 160, 164,
	7, 6, 2, 2, 161, 163, 5, 38, 20, 2, 162, 161, 3, 2, 2, 2, 163, 166, 3,
	2, 2, 2, 164, 162, 3, 2, 2, 2, 164, 165, 3, 2, 2, 2, 165, 167, 3, 2, 2,
	2, 166, 164, 3, 2, 2, 2, 167, 168, 7, 7, 2, 2, 168, 31, 3, 2, 2, 2, 169,
	170, 8, 17, 1, 2, 170, 172, 7, 26, 2, 2, 171, 173, 7, 3, 2, 2, 172, 171,
	3, 2, 2, 2, 172, 173, 3, 2, 2, 2, 173, 174, 3, 2, 2, 2, 174, 175, 5, 44,
	23, 2, 175, 33, 3, 2, 2, 2, 176, 177, 8, 18, 1, 2, 177, 179, 7, 26, 2,
	2, 178, 180, 7, 26, 2, 2, 179, 178, 3, 2, 2, 2, 179, 180, 3, 2, 2, 2, 180,
	181, 3, 2, 2, 2, 181, 185, 7, 6, 2, 2, 182, 184, 5, 38, 20, 2, 183, 182,
	3, 2, 2, 2, 184, 187, 3, 2, 2, 2, 185, 183, 3, 2, 2, 2, 185, 186, 3, 2,
	2, 2, 186, 188, 3, 2, 2, 2, 187, 185, 3, 2, 2, 2, 188, 189, 7, 7, 2, 2,
	189, 35, 3, 2, 2, 2, 190, 191, 8, 19, 1, 2, 191, 193, 7, 26, 2, 2, 192,
	194, 7, 3, 2, 2, 193, 192, 3, 2, 2, 2, 193, 194, 3, 2, 2, 2, 194, 195,
	3, 2, 2, 2, 195, 196, 5, 44, 23, 2, 196, 37, 3, 2, 2, 2, 197, 198, 6, 20,
	2, 2, 198, 201, 5, 40, 21, 2, 199, 201, 5, 42, 22, 2, 200, 197, 3, 2, 2,
	2, 200, 199, 3, 2, 2, 2, 201, 39, 3, 2, 2, 2, 202, 203, 8, 21, 1, 2, 203,
	204, 7, 26, 2, 2, 204, 206, 5, 44, 23, 2, 205, 207, 7, 24, 2, 2, 206, 205,
	3, 2, 2, 2, 206, 207, 3, 2, 2, 2, 207, 41, 3, 2, 2, 2, 208, 210, 7, 8,
	2, 2, 209, 208, 3, 2, 2, 2, 209, 210, 3, 2, 2, 2, 210, 211, 3, 2, 2, 2,
	211, 212, 7, 26, 2, 2, 212, 43, 3, 2, 2, 2, 213, 214, 8, 23, 1, 2, 214,
	222, 7, 26, 2, 2, 215, 222, 5, 48, 25, 2, 216, 222, 5, 50, 26, 2, 217,
	222, 7, 18, 2, 2, 218, 222, 7, 9, 2, 2, 219, 222, 5, 46, 24, 2, 220, 222,
	5, 30, 16, 2, 221, 213, 3, 2, 2, 2, 221, 215, 3, 2, 2, 2, 221, 216, 3,
	2, 2, 2, 221, 217, 3, 2, 2, 2, 221, 218, 3, 2, 2, 2, 221, 219, 3, 2, 2,
	2, 221, 220, 3, 2, 2, 2, 222, 45, 3, 2, 2, 2, 223, 224, 7, 8, 2, 2, 224,
	225, 8, 24, 1, 2, 225, 226, 7, 26, 2, 2, 226, 47, 3, 2, 2, 2, 227, 228,
	8, 25, 1, 2, 228, 229, 7, 26, 2, 2, 229, 230, 7, 10, 2, 2, 230, 231, 8,
	25, 1, 2, 231, 232, 7, 26, 2, 2, 232, 233, 7, 11, 2, 2, 233, 234, 5, 44,
	23, 2, 234, 49, 3, 2, 2, 2, 235, 236, 7, 10, 2, 2, 236, 237, 7, 11, 2,
	2, 237, 238, 5, 44, 23, 2, 238, 51, 3, 2, 2, 2, 239, 241, 5, 54, 28, 2,
	240, 239, 3, 2, 2, 2, 240, 241, 3, 2, 2, 2, 241, 242, 3, 2, 2, 2, 242,
	243, 5, 56, 29, 2, 243, 53, 3, 2, 2, 2, 244, 245, 7, 19, 2, 2, 245, 247,
	7, 4, 2, 2, 246, 248, 5, 70, 36, 2, 247, 246, 3, 2, 2, 2, 248, 249, 3,
	2, 2, 2, 249, 247, 3, 2, 2, 2, 249, 250, 3, 2, 2, 2, 250, 251, 3, 2, 2,
	2, 251, 252, 7, 5, 2, 2, 252, 55, 3, 2, 2, 2, 253, 254, 8, 29, 1, 2, 254,
	255, 7, 26, 2, 2, 255, 256, 5, 72, 37, 2, 256, 260, 7, 6, 2, 2, 257, 259,
	5, 58, 30, 2, 258, 257, 3, 2, 2, 2, 259, 262, 3, 2, 2, 2, 260, 258, 3,
	2, 2, 2, 260, 261, 3, 2, 2, 2, 261, 263, 3, 2, 2, 2, 262, 260, 3, 2, 2,
	2, 263, 264, 7, 7, 2, 2, 264, 57, 3, 2, 2, 2, 265, 267, 5, 60, 31, 2, 266,
	265, 3, 2, 2, 2, 266, 267, 3, 2, 2, 2, 267, 270, 3, 2, 2, 2, 268, 271,
	5, 54, 28, 2, 269, 271, 5, 62, 32, 2, 270, 268, 3, 2, 2, 2, 270, 269, 3,
	2, 2, 2, 271, 272, 3, 2, 2, 2, 272, 273, 5, 64, 33, 2, 273, 59, 3, 2, 2,
	2, 274, 276, 7, 16, 2, 2, 275, 277, 7, 4, 2, 2, 276, 275, 3, 2, 2, 2, 276,
	277, 3, 2, 2, 2, 277, 284, 3, 2, 2, 2, 278, 280, 5, 70, 36, 2, 279, 278,
	3, 2, 2, 2, 280, 281, 3, 2, 2, 2, 281, 279, 3, 2, 2, 2, 281, 282, 3, 2,
	2, 2, 282, 285, 3, 2, 2, 2, 283, 285, 7, 23, 2, 2, 284, 279, 3, 2, 2, 2,
	284, 283, 3, 2, 2, 2, 285, 287, 3, 2, 2, 2, 286, 288, 7, 5, 2, 2, 287,
	286, 3, 2, 2, 2, 287, 288, 3, 2, 2, 2, 288, 61, 3, 2, 2, 2, 289, 290, 7,
	17, 2, 2, 290, 291, 7, 26, 2, 2, 291, 63, 3, 2, 2, 2, 292, 293, 8, 33,
	1, 2, 293, 294, 7, 26, 2, 2, 294, 296, 5, 74, 38, 2, 295, 297, 5, 66, 34,
	2, 296, 295, 3, 2, 2, 2, 296, 297, 3, 2, 2, 2, 297, 299, 3, 2, 2, 2, 298,
	300, 5, 68, 35, 2, 299, 298, 3, 2, 2, 2, 299, 300, 3, 2, 2, 2, 300, 65,
	3, 2, 2, 2, 301, 303, 7, 4, 2, 2, 302, 304, 7, 26, 2, 2, 303, 302, 3, 2,
	2, 2, 303, 304, 3, 2, 2, 2, 304, 305, 3, 2, 2, 2, 305, 306, 7, 5, 2, 2,
	306, 67, 3, 2, 2, 2, 307, 308, 7, 12, 2, 2, 308, 310, 7, 4, 2, 2, 309,
	311, 5, 44, 23, 2, 310, 309, 3, 2, 2, 2, 310, 311, 3, 2, 2, 2, 311, 312,
	3, 2, 2, 2, 312, 313, 7, 5, 2, 2, 313, 69, 3, 2, 2, 2, 314, 315, 7, 26,
	2, 2, 315, 316, 8, 36, 1, 2, 316, 317, 7, 25, 2, 2, 317, 71, 3, 2, 2, 2,
	318, 320, 7, 26, 2, 2, 319, 321, 7, 13, 2, 2, 320, 319, 3, 2, 2, 2, 320,
	321, 3, 2, 2, 2, 321, 323, 3, 2, 2, 2, 322, 318, 3, 2, 2, 2, 323, 324,
	3, 2, 2, 2, 324, 322, 3, 2, 2, 2, 324, 325, 3, 2, 2, 2, 325, 73, 3, 2,
	2, 2, 326, 327, 7, 14, 2, 2, 327, 332, 7, 26, 2, 2, 328, 329, 7, 13, 2,
	2, 329, 331, 7, 26, 2, 2, 330, 328, 3, 2, 2, 2, 331, 334, 3, 2, 2, 2, 332,
	330, 3, 2, 2, 2, 332, 333, 3, 2, 2, 2, 333, 342, 3, 2, 2, 2, 334, 332,
	3, 2, 2, 2, 335, 336, 7, 15, 2, 2, 336, 339, 7, 26, 2, 2, 337, 338, 7,
	13, 2, 2, 338, 340, 7, 26, 2, 2, 339, 337, 3, 2, 2, 2, 339, 340, 3, 2,
	2, 2, 340, 342, 3, 2, 2, 2, 341, 326, 3, 2, 2, 2, 341, 335, 3, 2, 2, 2,
	342, 343, 3, 2, 2, 2, 343, 341, 3, 2, 2, 2, 343, 344, 3, 2, 2, 2, 344,
	347, 3, 2, 2, 2, 345, 347, 7, 14, 2, 2, 346, 341, 3, 2, 2, 2, 346, 345,
	3, 2, 2, 2, 347, 75, 3, 2, 2, 2, 41, 79, 87, 97, 109, 124, 130, 142, 149,
	153, 158, 164, 172, 179, 185, 193, 200, 206, 209, 221, 240, 249, 260, 266,
	270, 276, 281, 284, 287, 296, 299, 303, 310, 320, 324, 332, 339, 341, 343,
	346,
}

var literalNames = []string{
	"", "'='", "'('", "')'", "'{'", "'}'", "'*'", "'time.Time'", "'['", "']'",
	"'returns'", "'-'", "'/'", "'/:'", "'@doc'", "'@handler'", "'interface{}'",
	"'@server'",
}

var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "ATDOC", "ATHANDLER",
	"INTERFACE", "ATSERVER", "WS", "COMMENT", "LINE_COMMENT", "STRING", "RAW_STRING",
	"LINE_VALUE", "ID",
}

var ruleNames = []string{
	"api", "spec", "syntaxLit", "importSpec", "importLit", "importBlock", "importBlockValue",
	"importValue", "infoSpec", "typeSpec", "typeLit", "typeBlock", "typeLitBody",
	"typeBlockBody", "typeStruct", "typeAlias", "typeBlockStruct", "typeBlockAlias",
	"field", "normalField", "anonymousFiled", "dataType", "pointerType", "mapType",
	"arrayType", "serviceSpec", "atServer", "serviceApi", "serviceRoute", "atDoc",
	"atHandler", "route", "body", "replybody", "kvLit", "serviceName", "path",
}

type ApiParserParser struct {
	*antlr.BaseParser
}

// NewApiParserParser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *ApiParserParser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewApiParserParser(input antlr.TokenStream) *ApiParserParser {
	this := new(ApiParserParser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "ApiParser.g4"

	return this
}

// ApiParserParser tokens.
const (
	ApiParserParserEOF          = antlr.TokenEOF
	ApiParserParserT__0         = 1
	ApiParserParserT__1         = 2
	ApiParserParserT__2         = 3
	ApiParserParserT__3         = 4
	ApiParserParserT__4         = 5
	ApiParserParserT__5         = 6
	ApiParserParserT__6         = 7
	ApiParserParserT__7         = 8
	ApiParserParserT__8         = 9
	ApiParserParserT__9         = 10
	ApiParserParserT__10        = 11
	ApiParserParserT__11        = 12
	ApiParserParserT__12        = 13
	ApiParserParserATDOC        = 14
	ApiParserParserATHANDLER    = 15
	ApiParserParserINTERFACE    = 16
	ApiParserParserATSERVER     = 17
	ApiParserParserWS           = 18
	ApiParserParserCOMMENT      = 19
	ApiParserParserLINE_COMMENT = 20
	ApiParserParserSTRING       = 21
	ApiParserParserRAW_STRING   = 22
	ApiParserParserLINE_VALUE   = 23
	ApiParserParserID           = 24
)

// ApiParserParser rules.
const (
	ApiParserParserRULE_api              = 0
	ApiParserParserRULE_spec             = 1
	ApiParserParserRULE_syntaxLit        = 2
	ApiParserParserRULE_importSpec       = 3
	ApiParserParserRULE_importLit        = 4
	ApiParserParserRULE_importBlock      = 5
	ApiParserParserRULE_importBlockValue = 6
	ApiParserParserRULE_importValue      = 7
	ApiParserParserRULE_infoSpec         = 8
	ApiParserParserRULE_typeSpec         = 9
	ApiParserParserRULE_typeLit          = 10
	ApiParserParserRULE_typeBlock        = 11
	ApiParserParserRULE_typeLitBody      = 12
	ApiParserParserRULE_typeBlockBody    = 13
	ApiParserParserRULE_typeStruct       = 14
	ApiParserParserRULE_typeAlias        = 15
	ApiParserParserRULE_typeBlockStruct  = 16
	ApiParserParserRULE_typeBlockAlias   = 17
	ApiParserParserRULE_field            = 18
	ApiParserParserRULE_normalField      = 19
	ApiParserParserRULE_anonymousFiled   = 20
	ApiParserParserRULE_dataType         = 21
	ApiParserParserRULE_pointerType      = 22
	ApiParserParserRULE_mapType          = 23
	ApiParserParserRULE_arrayType        = 24
	ApiParserParserRULE_serviceSpec      = 25
	ApiParserParserRULE_atServer         = 26
	ApiParserParserRULE_serviceApi       = 27
	ApiParserParserRULE_serviceRoute     = 28
	ApiParserParserRULE_atDoc            = 29
	ApiParserParserRULE_atHandler        = 30
	ApiParserParserRULE_route            = 31
	ApiParserParserRULE_body             = 32
	ApiParserParserRULE_replybody        = 33
	ApiParserParserRULE_kvLit            = 34
	ApiParserParserRULE_serviceName      = 35
	ApiParserParserRULE_path             = 36
)

// IApiContext is an interface to support dynamic dispatch.
type IApiContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsApiContext differentiates from other interfaces.
	IsApiContext()
}

type ApiContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyApiContext() *ApiContext {
	p := new(ApiContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_api
	return p
}

func (*ApiContext) IsApiContext() {}

func NewApiContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ApiContext {
	p := new(ApiContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_api

	return p
}

func (s *ApiContext) GetParser() antlr.Parser { return s.parser }

func (s *ApiContext) AllSpec() []ISpecContext {
	ts := s.GetTypedRuleContexts(reflect.TypeOf((*ISpecContext)(nil)).Elem())
	tst := make([]ISpecContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISpecContext)
		}
	}

	return tst
}

func (s *ApiContext) Spec(i int) ISpecContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*ISpecContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISpecContext)
}

func (s *ApiContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ApiContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ApiContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitApi(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) Api() (localctx IApiContext) {
	localctx = NewApiContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, ApiParserParserRULE_api)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(77)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == ApiParserParserATSERVER || _la == ApiParserParserID {
		{
			p.SetState(74)
			p.Spec()
		}

		p.SetState(79)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISpecContext is an interface to support dynamic dispatch.
type ISpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSpecContext differentiates from other interfaces.
	IsSpecContext()
}

type SpecContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySpecContext() *SpecContext {
	p := new(SpecContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_spec
	return p
}

func (*SpecContext) IsSpecContext() {}

func NewSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SpecContext {
	p := new(SpecContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_spec

	return p
}

func (s *SpecContext) GetParser() antlr.Parser { return s.parser }

func (s *SpecContext) SyntaxLit() ISyntaxLitContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*ISyntaxLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISyntaxLitContext)
}

func (s *SpecContext) ImportSpec() IImportSpecContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IImportSpecContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportSpecContext)
}

func (s *SpecContext) InfoSpec() IInfoSpecContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IInfoSpecContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInfoSpecContext)
}

func (s *SpecContext) TypeSpec() ITypeSpecContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*ITypeSpecContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeSpecContext)
}

func (s *SpecContext) ServiceSpec() IServiceSpecContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IServiceSpecContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IServiceSpecContext)
}

func (s *SpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SpecContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitSpec(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) Spec() (localctx ISpecContext) {
	localctx = NewSpecContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, ApiParserParserRULE_spec)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(85)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(80)
			p.SyntaxLit()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(81)
			p.ImportSpec()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(82)
			p.InfoSpec()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(83)
			p.TypeSpec()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(84)
			p.ServiceSpec()
		}

	}

	return localctx
}

// ISyntaxLitContext is an interface to support dynamic dispatch.
type ISyntaxLitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetSyntaxToken returns the syntaxToken token.
	GetSyntaxToken() antlr.Token

	// GetAssign returns the assign token.
	GetAssign() antlr.Token

	// GetVersion returns the version token.
	GetVersion() antlr.Token

	// SetSyntaxToken sets the syntaxToken token.
	SetSyntaxToken(antlr.Token)

	// SetAssign sets the assign token.
	SetAssign(antlr.Token)

	// SetVersion sets the version token.
	SetVersion(antlr.Token)

	// IsSyntaxLitContext differentiates from other interfaces.
	IsSyntaxLitContext()
}

type SyntaxLitContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	syntaxToken antlr.Token
	assign      antlr.Token
	version     antlr.Token
}

func NewEmptySyntaxLitContext() *SyntaxLitContext {
	p := new(SyntaxLitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_syntaxLit
	return p
}

func (*SyntaxLitContext) IsSyntaxLitContext() {}

func NewSyntaxLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SyntaxLitContext {
	p := new(SyntaxLitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_syntaxLit

	return p
}
