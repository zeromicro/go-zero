package api // ApiParser
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/zeromicro/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 28, 391,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 3, 2, 7, 2, 86, 10, 2, 12,
	2, 14, 2, 89, 11, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 96, 10, 3, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 5, 5, 106, 10, 5, 3, 6, 3,
	6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3, 7, 6, 7, 116, 10, 7, 13, 7, 14, 7,
	117, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 3,
	10, 6, 10, 131, 10, 10, 13, 10, 14, 10, 132, 3, 10, 3, 10, 3, 11, 3, 11,
	5, 11, 139, 10, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 13, 3, 13, 3, 13, 3,
	13, 7, 13, 149, 10, 13, 12, 13, 14, 13, 152, 11, 13, 3, 13, 3, 13, 3, 14,
	3, 14, 5, 14, 158, 10, 14, 3, 15, 3, 15, 5, 15, 162, 10, 15, 3, 16, 3,
	16, 3, 16, 5, 16, 167, 10, 16, 3, 16, 3, 16, 7, 16, 171, 10, 16, 12, 16,
	14, 16, 174, 11, 16, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17, 5, 17, 181, 10,
	17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 5, 18, 188, 10, 18, 3, 18, 3, 18,
	7, 18, 192, 10, 18, 12, 18, 14, 18, 195, 11, 18, 3, 18, 3, 18, 3, 19, 3,
	19, 3, 19, 5, 19, 202, 10, 19, 3, 19, 3, 19, 3, 20, 3, 20, 3, 20, 5, 20,
	209, 10, 20, 3, 21, 3, 21, 3, 21, 3, 21, 5, 21, 215, 10, 21, 3, 22, 5,
	22, 218, 10, 22, 3, 22, 3, 22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23,
	3, 23, 3, 23, 5, 23, 230, 10, 23, 3, 24, 3, 24, 3, 24, 3, 24, 3, 25, 3,
	25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 26, 3, 26, 3, 26, 3, 26,
	3, 27, 5, 27, 249, 10, 27, 3, 27, 3, 27, 3, 28, 3, 28, 3, 28, 6, 28, 256,
	10, 28, 13, 28, 14, 28, 257, 3, 28, 3, 28, 3, 29, 3, 29, 3, 29, 3, 29,
	3, 29, 7, 29, 267, 10, 29, 12, 29, 14, 29, 270, 11, 29, 3, 29, 3, 29, 3,
	30, 5, 30, 275, 10, 30, 3, 30, 3, 30, 5, 30, 279, 10, 30, 3, 30, 7, 30,
	282, 10, 30, 12, 30, 14, 30, 285, 11, 30, 3, 30, 3, 30, 3, 31, 3, 31, 5,
	31, 291, 10, 31, 3, 31, 6, 31, 294, 10, 31, 13, 31, 14, 31, 295, 3, 31,
	5, 31, 299, 10, 31, 3, 31, 5, 31, 302, 10, 31, 3, 32, 3, 32, 3, 32, 3,
	33, 3, 33, 3, 33, 3, 33, 3, 33, 6, 33, 312, 10, 33, 13, 33, 14, 33, 313,
	3, 33, 5, 33, 317, 10, 33, 3, 33, 3, 33, 3, 34, 6, 34, 322, 10, 34, 13,
	34, 14, 34, 323, 3, 35, 3, 35, 3, 35, 3, 35, 5, 35, 330, 10, 35, 3, 35,
	5, 35, 333, 10, 35, 3, 36, 3, 36, 5, 36, 337, 10, 36, 3, 36, 3, 36, 3,
	37, 3, 37, 3, 37, 5, 37, 344, 10, 37, 3, 37, 3, 37, 3, 38, 3, 38, 3, 38,
	3, 38, 3, 39, 3, 39, 3, 39, 3, 39, 3, 40, 3, 40, 5, 40, 358, 10, 40, 6,
	40, 360, 10, 40, 13, 40, 14, 40, 361, 3, 41, 3, 41, 3, 41, 3, 41, 7, 41,
	368, 10, 41, 12, 41, 14, 41, 371, 11, 41, 3, 41, 3, 41, 3, 41, 3, 41, 5,
	41, 377, 10, 41, 6, 41, 379, 10, 41, 13, 41, 14, 41, 380, 3, 41, 5, 41,
	384, 10, 41, 3, 42, 6, 42, 387, 10, 42, 13, 42, 14, 42, 388, 3, 42, 2,
	2, 43, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34,
	36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70,
	72, 74, 76, 78, 80, 82, 2, 3, 3, 2, 27, 28, 2, 401, 2, 87, 3, 2, 2, 2,
	4, 95, 3, 2, 2, 2, 6, 97, 3, 2, 2, 2, 8, 105, 3, 2, 2, 2, 10, 107, 3, 2,
	2, 2, 12, 111, 3, 2, 2, 2, 14, 121, 3, 2, 2, 2, 16, 123, 3, 2, 2, 2, 18,
	126, 3, 2, 2, 2, 20, 138, 3, 2, 2, 2, 22, 140, 3, 2, 2, 2, 24, 144, 3,
	2, 2, 2, 26, 157, 3, 2, 2, 2, 28, 161, 3, 2, 2, 2, 30, 163, 3, 2, 2, 2,
	32, 177, 3, 2, 2, 2, 34, 184, 3, 2, 2, 2, 36, 198, 3, 2, 2, 2, 38, 208,
	3, 2, 2, 2, 40, 210, 3, 2, 2, 2, 42, 217, 3, 2, 2, 2, 44, 229, 3, 2, 2,
	2, 46, 231, 3, 2, 2, 2, 48, 235, 3, 2, 2, 2, 50, 243, 3, 2, 2, 2, 52, 248,
	3, 2, 2, 2, 54, 252, 3, 2, 2, 2, 56, 261, 3, 2, 2, 2, 58, 274, 3, 2, 2,
	2, 60, 288, 3, 2, 2, 2, 62, 303, 3, 2, 2, 2, 64, 306, 3, 2, 2, 2, 66, 321,
	3, 2, 2, 2, 68, 325, 3, 2, 2, 2, 70, 334, 3, 2, 2, 2, 72, 340, 3, 2, 2,
	2, 74, 347, 3, 2, 2, 2, 76, 351, 3, 2, 2, 2, 78, 359, 3, 2, 2, 2, 80, 383,
	3, 2, 2, 2, 82, 386, 3, 2, 2, 2, 84, 86, 5, 4, 3, 2, 85, 84, 3, 2, 2, 2,
	86, 89, 3, 2, 2, 2, 87, 85, 3, 2, 2, 2, 87, 88, 3, 2, 2, 2, 88, 3, 3, 2,
	2, 2, 89, 87, 3, 2, 2, 2, 90, 96, 5, 6, 4, 2, 91, 96, 5, 8, 5, 2, 92, 96,
	5, 18, 10, 2, 93, 96, 5, 20, 11, 2, 94, 96, 5, 52, 27, 2, 95, 90, 3, 2,
	2, 2, 95, 91, 3, 2, 2, 2, 95, 92, 3, 2, 2, 2, 95, 93, 3, 2, 2, 2, 95, 94,
	3, 2, 2, 2, 96, 5, 3, 2, 2, 2, 97, 98, 8, 4, 1, 2, 98, 99, 7, 27, 2, 2,
	99, 100, 7, 3, 2, 2, 100, 101, 8, 4, 1, 2, 101, 102, 7, 24, 2, 2, 102,
	7, 3, 2, 2, 2, 103, 106, 5, 10, 6, 2, 104, 106, 5, 12, 7, 2, 105, 103,
	3, 2, 2, 2, 105, 104, 3, 2, 2, 2, 106, 9, 3, 2, 2, 2, 107, 108, 8, 6, 1,
	2, 108, 109, 7, 27, 2, 2, 109, 110, 5, 16, 9, 2, 110, 11, 3, 2, 2, 2, 111,
	112, 8, 7, 1, 2, 112, 113, 7, 27, 2, 2, 113, 115, 7, 4, 2, 2, 114, 116,
	5, 14, 8, 2, 115, 114, 3, 2, 2, 2, 116, 117, 3, 2, 2, 2, 117, 115, 3, 2,
	2, 2, 117, 118, 3, 2, 2, 2, 118, 119, 3, 2, 2, 2, 119, 120, 7, 5, 2, 2,
	120, 13, 3, 2, 2, 2, 121, 122, 5, 16, 9, 2, 122, 15, 3, 2, 2, 2, 123, 124,
	8, 9, 1, 2, 124, 125, 7, 24, 2, 2, 125, 17, 3, 2, 2, 2, 126, 127, 8, 10,
	1, 2, 127, 128, 7, 27, 2, 2, 128, 130, 7, 4, 2, 2, 129, 131, 5, 74, 38,
	2, 130, 129, 3, 2, 2, 2, 131, 132, 3, 2, 2, 2, 132, 130, 3, 2, 2, 2, 132,
	133, 3, 2, 2, 2, 133, 134, 3, 2, 2, 2, 134, 135, 7, 5, 2, 2, 135, 19, 3,
	2, 2, 2, 136, 139, 5, 22, 12, 2, 137, 139, 5, 24, 13, 2, 138, 136, 3, 2,
	2, 2, 138, 137, 3, 2, 2, 2, 139, 21, 3, 2, 2, 2, 140, 141, 8, 12, 1, 2,
	141, 142, 7, 27, 2, 2, 142, 143, 5, 26, 14, 2, 143, 23, 3, 2, 2, 2, 144,
	145, 8, 13, 1, 2, 145, 146, 7, 27, 2, 2, 146, 150, 7, 4, 2, 2, 147, 149,
	5, 28, 15, 2, 148, 147, 3, 2, 2, 2, 149, 152, 3, 2, 2, 2, 150, 148, 3,
	2, 2, 2, 150, 151, 3, 2, 2, 2, 151, 153, 3, 2, 2, 2, 152, 150, 3, 2, 2,
	2, 153, 154, 7, 5, 2, 2, 154, 25, 3, 2, 2, 2, 155, 158, 5, 30, 16, 2, 156,
	158, 5, 32, 17, 2, 157, 155, 3, 2, 2, 2, 157, 156, 3, 2, 2, 2, 158, 27,
	3, 2, 2, 2, 159, 162, 5, 34, 18, 2, 160, 162, 5, 36, 19, 2, 161, 159, 3,
	2, 2, 2, 161, 160, 3, 2, 2, 2, 162, 29, 3, 2, 2, 2, 163, 164, 8, 16, 1,
	2, 164, 166, 7, 27, 2, 2, 165, 167, 7, 27, 2, 2, 166, 165, 3, 2, 2, 2,
	166, 167, 3, 2, 2, 2, 167, 168, 3, 2, 2, 2, 168, 172, 7, 6, 2, 2, 169,
	171, 5, 38, 20, 2, 170, 169, 3, 2, 2, 2, 171, 174, 3, 2, 2, 2, 172, 170,
	3, 2, 2, 2, 172, 173, 3, 2, 2, 2, 173, 175, 3, 2, 2, 2, 174, 172, 3, 2,
	2, 2, 175, 176, 7, 7, 2, 2, 176, 31, 3, 2, 2, 2, 177, 178, 8, 17, 1, 2,
	178, 180, 7, 27, 2, 2, 179, 181, 7, 3, 2, 2, 180, 179, 3, 2, 2, 2, 180,
	181, 3, 2, 2, 2, 181, 182, 3, 2, 2, 2, 182, 183, 5, 44, 23, 2, 183, 33,
	3, 2, 2, 2, 184, 185, 8, 18, 1, 2, 185, 187, 7, 27, 2, 2, 186, 188, 7,
	27, 2, 2, 187, 186, 3, 2, 2, 2, 187, 188, 3, 2, 2, 2, 188, 189, 3, 2, 2,
	2, 189, 193, 7, 6, 2, 2, 190, 192, 5, 38, 20, 2, 191, 190, 3, 2, 2, 2,
	192, 195, 3, 2, 2, 2, 193, 191, 3, 2, 2, 2, 193, 194, 3, 2, 2, 2, 194,
	196, 3, 2, 2, 2, 195, 193, 3, 2, 2, 2, 196, 197, 7, 7, 2, 2, 197, 35, 3,
	2, 2, 2, 198, 199, 8, 19, 1, 2, 199, 201, 7, 27, 2, 2, 200, 202, 7, 3,
	2, 2, 201, 200, 3, 2, 2, 2, 201, 202, 3, 2, 2, 2, 202, 203, 3, 2, 2, 2,
	203, 204, 5, 44, 23, 2, 204, 37, 3, 2, 2, 2, 205, 206, 6, 20, 2, 2, 206,
	209, 5, 40, 21, 2, 207, 209, 5, 42, 22, 2, 208, 205, 3, 2, 2, 2, 208, 207,
	3, 2, 2, 2, 209, 39, 3, 2, 2, 2, 210, 211, 8, 21, 1, 2, 211, 212, 7, 27,
	2, 2, 212, 214, 5, 44, 23, 2, 213, 215, 7, 25, 2, 2, 214, 213, 3, 2, 2,
	2, 214, 215, 3, 2, 2, 2, 215, 41, 3, 2, 2, 2, 216, 218, 7, 8, 2, 2, 217,
	216, 3, 2, 2, 2, 217, 218, 3, 2, 2, 2, 218, 219, 3, 2, 2, 2, 219, 220,
	7, 27, 2, 2, 220, 43, 3, 2, 2, 2, 221, 222, 8, 23, 1, 2, 222, 230, 7, 27,
	2, 2, 223, 230, 5, 48, 25, 2, 224, 230, 5, 50, 26, 2, 225, 230, 7, 18,
	2, 2, 226, 230, 7, 9, 2, 2, 227, 230, 5, 46, 24, 2, 228, 230, 5, 30, 16,
	2, 229, 221, 3, 2, 2, 2, 229, 223, 3, 2, 2, 2, 229, 224, 3, 2, 2, 2, 229,
	225, 3, 2, 2, 2, 229, 226, 3, 2, 2, 2, 229, 227, 3, 2, 2, 2, 229, 228,
	3, 2, 2, 2, 230, 45, 3, 2, 2, 2, 231, 232, 7, 8, 2, 2, 232, 233, 8, 24,
	1, 2, 233, 234, 7, 27, 2, 2, 234, 47, 3, 2, 2, 2, 235, 236, 8, 25, 1, 2,
	236, 237, 7, 27, 2, 2, 237, 238, 7, 10, 2, 2, 238, 239, 8, 25, 1, 2, 239,
	240, 7, 27, 2, 2, 240, 241, 7, 11, 2, 2, 241, 242, 5, 44, 23, 2, 242, 49,
	3, 2, 2, 2, 243, 244, 7, 10, 2, 2, 244, 245, 7, 11, 2, 2, 245, 246, 5,
	44, 23, 2, 246, 51, 3, 2, 2, 2, 247, 249, 5, 54, 28, 2, 248, 247, 3, 2,
	2, 2, 248, 249, 3, 2, 2, 2, 249, 250, 3, 2, 2, 2, 250, 251, 5, 56, 29,
	2, 251, 53, 3, 2, 2, 2, 252, 253, 7, 19, 2, 2, 253, 255, 7, 4, 2, 2, 254,
	256, 5, 74, 38, 2, 255, 254, 3, 2, 2, 2, 256, 257, 3, 2, 2, 2, 257, 255,
	3, 2, 2, 2, 257, 258, 3, 2, 2, 2, 258, 259, 3, 2, 2, 2, 259, 260, 7, 5,
	2, 2, 260, 55, 3, 2, 2, 2, 261, 262, 8, 29, 1, 2, 262, 263, 7, 27, 2, 2,
	263, 264, 5, 78, 40, 2, 264, 268, 7, 6, 2, 2, 265, 267, 5, 58, 30, 2, 266,
	265, 3, 2, 2, 2, 267, 270, 3, 2, 2, 2, 268, 266, 3, 2, 2, 2, 268, 269,
	3, 2, 2, 2, 269, 271, 3, 2, 2, 2, 270, 268, 3, 2, 2, 2, 271, 272, 7, 7,
	2, 2, 272, 57, 3, 2, 2, 2, 273, 275, 5, 60, 31, 2, 274, 273, 3, 2, 2, 2,
	274, 275, 3, 2, 2, 2, 275, 278, 3, 2, 2, 2, 276, 279, 5, 54, 28, 2, 277,
	279, 5, 62, 32, 2, 278, 276, 3, 2, 2, 2, 278, 277, 3, 2, 2, 2, 279, 283,
	3, 2, 2, 2, 280, 282, 5, 64, 33, 2, 281, 280, 3, 2, 2, 2, 282, 285, 3,
	2, 2, 2, 283, 281, 3, 2, 2, 2, 283, 284, 3, 2, 2, 2, 284, 286, 3, 2, 2,
	2, 285, 283, 3, 2, 2, 2, 286, 287, 5, 68, 35, 2, 287, 59, 3, 2, 2, 2, 288,
	290, 7, 16, 2, 2, 289, 291, 7, 4, 2, 2, 290, 289, 3, 2, 2, 2, 290, 291,
	3, 2, 2, 2, 291, 298, 3, 2, 2, 2, 292, 294, 5, 74, 38, 2, 293, 292, 3,
	2, 2, 2, 294, 295, 3, 2, 2, 2, 295, 293, 3, 2, 2, 2, 295, 296, 3, 2, 2,
	2, 296, 299, 3, 2, 2, 2, 297, 299, 7, 24, 2, 2, 298, 293, 3, 2, 2, 2, 298,
	297, 3, 2, 2, 2, 299, 301, 3, 2, 2, 2, 300, 302, 7, 5, 2, 2, 301, 300,
	3, 2, 2, 2, 301, 302, 3, 2, 2, 2, 302, 61, 3, 2, 2, 2, 303, 304, 7, 17,
	2, 2, 304, 305, 7, 27, 2, 2, 305, 63, 3, 2, 2, 2, 306, 307, 7, 20, 2, 2,
	307, 308, 7, 12, 2, 2, 308, 309, 5, 66, 34, 2, 309, 316, 7, 4, 2, 2, 310,
	312, 5, 76, 39, 2, 311, 310, 3, 2, 2, 2, 312, 313, 3, 2, 2, 2, 313, 311,
	3, 2, 2, 2, 313, 314, 3, 2, 2, 2, 314, 317, 3, 2, 2, 2, 315, 317, 5, 44,
	23, 2, 316, 311, 3, 2, 2, 2, 316, 315, 3, 2, 2, 2, 317, 318, 3, 2, 2, 2,
	318, 319, 7, 5, 2, 2, 319, 65, 3, 2, 2, 2, 320, 322, 9, 2, 2, 2, 321, 320,
	3, 2, 2, 2, 322, 323, 3, 2, 2, 2, 323, 321, 3, 2, 2, 2, 323, 324, 3, 2,
	2, 2, 324, 67, 3, 2, 2, 2, 325, 326, 8, 35, 1, 2, 326, 327, 7, 27, 2, 2,
	327, 329, 5, 80, 41, 2, 328, 330, 5, 70, 36, 2, 329, 328, 3, 2, 2, 2, 329,
	330, 3, 2, 2, 2, 330, 332, 3, 2, 2, 2, 331, 333, 5, 72, 37, 2, 332, 331,
	3, 2, 2, 2, 332, 333, 3, 2, 2, 2, 333, 69, 3, 2, 2, 2, 334, 336, 7, 4,
	2, 2, 335, 337, 7, 27, 2, 2, 336, 335, 3, 2, 2, 2, 336, 337, 3, 2, 2, 2,
	337, 338, 3, 2, 2, 2, 338, 339, 7, 5, 2, 2, 339, 71, 3, 2, 2, 2, 340, 341,
	7, 13, 2, 2, 341, 343, 7, 4, 2, 2, 342, 344, 5, 44, 23, 2, 343, 342, 3,
	2, 2, 2, 343, 344, 3, 2, 2, 2, 344, 345, 3, 2, 2, 2, 345, 346, 7, 5, 2,
	2, 346, 73, 3, 2, 2, 2, 347, 348, 7, 27, 2, 2, 348, 349, 8, 38, 1, 2, 349,
	350, 7, 26, 2, 2, 350, 75, 3, 2, 2, 2, 351, 352, 5, 66, 34, 2, 352, 353,
	8, 39, 1, 2, 353, 354, 7, 26, 2, 2, 354, 77, 3, 2, 2, 2, 355, 357, 7, 27,
	2, 2, 356, 358, 7, 12, 2, 2, 357, 356, 3, 2, 2, 2, 357, 358, 3, 2, 2, 2,
	358, 360, 3, 2, 2, 2, 359, 355, 3, 2, 2, 2, 360, 361, 3, 2, 2, 2, 361,
	359, 3, 2, 2, 2, 361, 362, 3, 2, 2, 2, 362, 79, 3, 2, 2, 2, 363, 364, 7,
	14, 2, 2, 364, 369, 5, 82, 42, 2, 365, 366, 7, 12, 2, 2, 366, 368, 5, 82,
	42, 2, 367, 365, 3, 2, 2, 2, 368, 371, 3, 2, 2, 2, 369, 367, 3, 2, 2, 2,
	369, 370, 3, 2, 2, 2, 370, 379, 3, 2, 2, 2, 371, 369, 3, 2, 2, 2, 372,
	373, 7, 15, 2, 2, 373, 376, 5, 82, 42, 2, 374, 375, 7, 12, 2, 2, 375, 377,
	5, 82, 42, 2, 376, 374, 3, 2, 2, 2, 376, 377, 3, 2, 2, 2, 377, 379, 3,
	2, 2, 2, 378, 363, 3, 2, 2, 2, 378, 372, 3, 2, 2, 2, 379, 380, 3, 2, 2,
	2, 380, 378, 3, 2, 2, 2, 380, 381, 3, 2, 2, 2, 381, 384, 3, 2, 2, 2, 382,
	384, 7, 14, 2, 2, 383, 378, 3, 2, 2, 2, 383, 382, 3, 2, 2, 2, 384, 81,
	3, 2, 2, 2, 385, 387, 9, 2, 2, 2, 386, 385, 3, 2, 2, 2, 387, 388, 3, 2,
	2, 2, 388, 386, 3, 2, 2, 2, 388, 389, 3, 2, 2, 2, 389, 83, 3, 2, 2, 2,
	46, 87, 95, 105, 117, 132, 138, 150, 157, 161, 166, 172, 180, 187, 193,
	201, 208, 214, 217, 229, 248, 257, 268, 274, 278, 283, 290, 295, 298, 301,
	313, 316, 323, 329, 332, 336, 343, 357, 361, 369, 376, 378, 380, 383, 388,
}
var literalNames = []string{
	"", "'='", "'('", "')'", "'{'", "'}'", "'*'", "'time.Time'", "'['", "']'",
	"'-'", "'returns'", "'/'", "'/:'", "'@doc'", "'@handler'", "'interface{}'",
	"'@server'", "'@respdoc'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "ATDOC", "ATHANDLER",
	"INTERFACE", "ATSERVER", "ATRESPDOC", "WS", "COMMENT", "LINE_COMMENT",
	"STRING", "RAW_STRING", "LINE_VALUE", "ID", "LetterOrDigit",
}

var ruleNames = []string{
	"api", "spec", "syntaxLit", "importSpec", "importLit", "importBlock", "importBlockValue",
	"importValue", "infoSpec", "typeSpec", "typeLit", "typeBlock", "typeLitBody",
	"typeBlockBody", "typeStruct", "typeAlias", "typeBlockStruct", "typeBlockAlias",
	"field", "normalField", "anonymousFiled", "dataType", "pointerType", "mapType",
	"arrayType", "serviceSpec", "atServer", "serviceApi", "serviceRoute", "atDoc",
	"atHandler", "atRespDoc", "respDocCode", "route", "body", "replybody",
	"kvLit", "respDocKvLit", "serviceName", "path", "pathItem",
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
	ApiParserParserEOF           = antlr.TokenEOF
	ApiParserParserT__0          = 1
	ApiParserParserT__1          = 2
	ApiParserParserT__2          = 3
	ApiParserParserT__3          = 4
	ApiParserParserT__4          = 5
	ApiParserParserT__5          = 6
	ApiParserParserT__6          = 7
	ApiParserParserT__7          = 8
	ApiParserParserT__8          = 9
	ApiParserParserT__9          = 10
	ApiParserParserT__10         = 11
	ApiParserParserT__11         = 12
	ApiParserParserT__12         = 13
	ApiParserParserATDOC         = 14
	ApiParserParserATHANDLER     = 15
	ApiParserParserINTERFACE     = 16
	ApiParserParserATSERVER      = 17
	ApiParserParserATRESPDOC     = 18
	ApiParserParserWS            = 19
	ApiParserParserCOMMENT       = 20
	ApiParserParserLINE_COMMENT  = 21
	ApiParserParserSTRING        = 22
	ApiParserParserRAW_STRING    = 23
	ApiParserParserLINE_VALUE    = 24
	ApiParserParserID            = 25
	ApiParserParserLetterOrDigit = 26
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
	ApiParserParserRULE_atRespDoc        = 31
	ApiParserParserRULE_respDocCode      = 32
	ApiParserParserRULE_route            = 33
	ApiParserParserRULE_body             = 34
	ApiParserParserRULE_replybody        = 35
	ApiParserParserRULE_kvLit            = 36
	ApiParserParserRULE_respDocKvLit     = 37
	ApiParserParserRULE_serviceName      = 38
	ApiParserParserRULE_path             = 39
	ApiParserParserRULE_pathItem         = 40
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
	var p = new(ApiContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_api
	return p
}

func (*ApiContext) IsApiContext() {}

func NewApiContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ApiContext {
	var p = new(ApiContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_api

	return p
}

func (s *ApiContext) GetParser() antlr.Parser { return s.parser }

func (s *ApiContext) AllSpec() []ISpecContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISpecContext)(nil)).Elem())
	var tst = make([]ISpecContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISpecContext)
		}
	}

	return tst
}

func (s *ApiContext) Spec(i int) ISpecContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISpecContext)(nil)).Elem(), i)

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
	p.SetState(85)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == ApiParserParserATSERVER || _la == ApiParserParserID {
		{
			p.SetState(82)
			p.Spec()
		}

		p.SetState(87)
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
	var p = new(SpecContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_spec
	return p
}

func (*SpecContext) IsSpecContext() {}

func NewSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SpecContext {
	var p = new(SpecContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_spec

	return p
}

func (s *SpecContext) GetParser() antlr.Parser { return s.parser }

func (s *SpecContext) SyntaxLit() ISyntaxLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISyntaxLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISyntaxLitContext)
}

func (s *SpecContext) ImportSpec() IImportSpecContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImportSpecContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportSpecContext)
}

func (s *SpecContext) InfoSpec() IInfoSpecContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInfoSpecContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInfoSpecContext)
}

func (s *SpecContext) TypeSpec() ITypeSpecContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeSpecContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeSpecContext)
}

func (s *SpecContext) ServiceSpec() IServiceSpecContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IServiceSpecContext)(nil)).Elem(), 0)

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

	p.SetState(93)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(88)
			p.SyntaxLit()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(89)
			p.ImportSpec()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(90)
			p.InfoSpec()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(91)
			p.TypeSpec()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(92)
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
	var p = new(SyntaxLitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_syntaxLit
	return p
}

func (*SyntaxLitContext) IsSyntaxLitContext() {}

func NewSyntaxLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SyntaxLitContext {
	var p = new(SyntaxLitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_syntaxLit

	return p
}

func (s *SyntaxLitContext) GetParser() antlr.Parser { return s.parser }

func (s *SyntaxLitContext) GetSyntaxToken() antlr.Token { return s.syntaxToken }

func (s *SyntaxLitContext) GetAssign() antlr.Token { return s.assign }

func (s *SyntaxLitContext) GetVersion() antlr.Token { return s.version }

func (s *SyntaxLitContext) SetSyntaxToken(v antlr.Token) { s.syntaxToken = v }

func (s *SyntaxLitContext) SetAssign(v antlr.Token) { s.assign = v }

func (s *SyntaxLitContext) SetVersion(v antlr.Token) { s.version = v }

func (s *SyntaxLitContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}
