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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 27, 356,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 3, 2, 7, 2, 80, 10, 2, 12, 2, 14, 2, 83, 11, 2, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 5, 3, 90, 10, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3,
	5, 5, 5, 100, 10, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3, 7, 6,
	7, 110, 10, 7, 13, 7, 14, 7, 111, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3,
	9, 3, 10, 3, 10, 3, 10, 3, 10, 6, 10, 125, 10, 10, 13, 10, 14, 10, 126,
	3, 10, 3, 10, 3, 11, 3, 11, 5, 11, 133, 10, 11, 3, 12, 3, 12, 3, 12, 3,
	12, 3, 13, 3, 13, 3, 13, 3, 13, 7, 13, 143, 10, 13, 12, 13, 14, 13, 146,
	11, 13, 3, 13, 3, 13, 3, 14, 3, 14, 5, 14, 152, 10, 14, 3, 15, 3, 15, 5,
	15, 156, 10, 15, 3, 16, 3, 16, 3, 16, 5, 16, 161, 10, 16, 3, 16, 3, 16,
	7, 16, 165, 10, 16, 12, 16, 14, 16, 168, 11, 16, 3, 16, 3, 16, 3, 17, 3,
	17, 3, 17, 5, 17, 175, 10, 17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 5, 18,
	182, 10, 18, 3, 18, 3, 18, 7, 18, 186, 10, 18, 12, 18, 14, 18, 189, 11,
	18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 5, 19, 196, 10, 19, 3, 19, 3, 19,
	3, 20, 3, 20, 3, 20, 5, 20, 203, 10, 20, 3, 21, 3, 21, 3, 21, 3, 21, 5,
	21, 209, 10, 21, 3, 22, 5, 22, 212, 10, 22, 3, 22, 3, 22, 3, 23, 3, 23,
	3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 5, 23, 224, 10, 23, 3, 24, 3,
	24, 3, 24, 3, 24, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25,
	3, 26, 3, 26, 3, 26, 3, 26, 3, 27, 5, 27, 243, 10, 27, 3, 27, 3, 27, 3,
	28, 3, 28, 3, 28, 6, 28, 250, 10, 28, 13, 28, 14, 28, 251, 3, 28, 3, 28,
	3, 29, 3, 29, 3, 29, 3, 29, 3, 29, 7, 29, 261, 10, 29, 12, 29, 14, 29,
	264, 11, 29, 3, 29, 3, 29, 3, 30, 5, 30, 269, 10, 30, 3, 30, 3, 30, 5,
	30, 273, 10, 30, 3, 30, 3, 30, 3, 31, 3, 31, 5, 31, 279, 10, 31, 3, 31,
	6, 31, 282, 10, 31, 13, 31, 14, 31, 283, 3, 31, 5, 31, 287, 10, 31, 3,
	31, 5, 31, 290, 10, 31, 3, 32, 3, 32, 3, 32, 3, 33, 3, 33, 3, 33, 3, 33,
	5, 33, 299, 10, 33, 3, 33, 5, 33, 302, 10, 33, 3, 34, 3, 34, 5, 34, 306,
	10, 34, 3, 34, 3, 34, 3, 35, 3, 35, 3, 35, 5, 35, 313, 10, 35, 3, 35, 3,
	35, 3, 36, 3, 36, 3, 36, 3, 36, 3, 37, 3, 37, 5, 37, 323, 10, 37, 6, 37,
	325, 10, 37, 13, 37, 14, 37, 326, 3, 38, 3, 38, 3, 38, 3, 38, 7, 38, 333,
	10, 38, 12, 38, 14, 38, 336, 11, 38, 3, 38, 3, 38, 3, 38, 3, 38, 5, 38,
	342, 10, 38, 6, 38, 344, 10, 38, 13, 38, 14, 38, 345, 3, 38, 5, 38, 349,
	10, 38, 3, 39, 6, 39, 352, 10, 39, 13, 39, 14, 39, 353, 3, 39, 2, 2, 40,
	2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38,
	40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74,
	76, 2, 3, 3, 2, 26, 27, 2, 365, 2, 81, 3, 2, 2, 2, 4, 89, 3, 2, 2, 2, 6,
	91, 3, 2, 2, 2, 8, 99, 3, 2, 2, 2, 10, 101, 3, 2, 2, 2, 12, 105, 3, 2,
	2, 2, 14, 115, 3, 2, 2, 2, 16, 117, 3, 2, 2, 2, 18, 120, 3, 2, 2, 2, 20,
	132, 3, 2, 2, 2, 22, 134, 3, 2, 2, 2, 24, 138, 3, 2, 2, 2, 26, 151, 3,
	2, 2, 2, 28, 155, 3, 2, 2, 2, 30, 157, 3, 2, 2, 2, 32, 171, 3, 2, 2, 2,
	34, 178, 3, 2, 2, 2, 36, 192, 3, 2, 2, 2, 38, 202, 3, 2, 2, 2, 40, 204,
	3, 2, 2, 2, 42, 211, 3, 2, 2, 2, 44, 223, 3, 2, 2, 2, 46, 225, 3, 2, 2,
	2, 48, 229, 3, 2, 2, 2, 50, 237, 3, 2, 2, 2, 52, 242, 3, 2, 2, 2, 54, 246,
	3, 2, 2, 2, 56, 255, 3, 2, 2, 2, 58, 268, 3, 2, 2, 2, 60, 276, 3, 2, 2,
	2, 62, 291, 3, 2, 2, 2, 64, 294, 3, 2, 2, 2, 66, 303, 3, 2, 2, 2, 68, 309,
	3, 2, 2, 2, 70, 316, 3, 2, 2, 2, 72, 324, 3, 2, 2, 2, 74, 348, 3, 2, 2,
	2, 76, 351, 3, 2, 2, 2, 78, 80, 5, 4, 3, 2, 79, 78, 3, 2, 2, 2, 80, 83,
	3, 2, 2, 2, 81, 79, 3, 2, 2, 2, 81, 82, 3, 2, 2, 2, 82, 3, 3, 2, 2, 2,
	83, 81, 3, 2, 2, 2, 84, 90, 5, 6, 4, 2, 85, 90, 5, 8, 5, 2, 86, 90, 5,
	18, 10, 2, 87, 90, 5, 20, 11, 2, 88, 90, 5, 52, 27, 2, 89, 84, 3, 2, 2,
	2, 89, 85, 3, 2, 2, 2, 89, 86, 3, 2, 2, 2, 89, 87, 3, 2, 2, 2, 89, 88,
	3, 2, 2, 2, 90, 5, 3, 2, 2, 2, 91, 92, 8, 4, 1, 2, 92, 93, 7, 26, 2, 2,
	93, 94, 7, 3, 2, 2, 94, 95, 8, 4, 1, 2, 95, 96, 7, 23, 2, 2, 96, 7, 3,
	2, 2, 2, 97, 100, 5, 10, 6, 2, 98, 100, 5, 12, 7, 2, 99, 97, 3, 2, 2, 2,
	99, 98, 3, 2, 2, 2, 100, 9, 3, 2, 2, 2, 101, 102, 8, 6, 1, 2, 102, 103,
	7, 26, 2, 2, 103, 104, 5, 16, 9, 2, 104, 11, 3, 2, 2, 2, 105, 106, 8, 7,
	1, 2, 106, 107, 7, 26, 2, 2, 107, 109, 7, 4, 2, 2, 108, 110, 5, 14, 8,
	2, 109, 108, 3, 2, 2, 2, 110, 111, 3, 2, 2, 2, 111, 109, 3, 2, 2, 2, 111,
	112, 3, 2, 2, 2, 112, 113, 3, 2, 2, 2, 113, 114, 7, 5, 2, 2, 114, 13, 3,
	2, 2, 2, 115, 116, 5, 16, 9, 2, 116, 15, 3, 2, 2, 2, 117, 118, 8, 9, 1,
	2, 118, 119, 7, 23, 2, 2, 119, 17, 3, 2, 2, 2, 120, 121, 8, 10, 1, 2, 121,
	122, 7, 26, 2, 2, 122, 124, 7, 4, 2, 2, 123, 125, 5, 70, 36, 2, 124, 123,
	3, 2, 2, 2, 125, 126, 3, 2, 2, 2, 126, 124, 3, 2, 2, 2, 126, 127, 3, 2,
	2, 2, 127, 128, 3, 2, 2, 2, 128, 129, 7, 5, 2, 2, 129, 19, 3, 2, 2, 2,
	130, 133, 5, 22, 12, 2, 131, 133, 5, 24, 13, 2, 132, 130, 3, 2, 2, 2, 132,
	131, 3, 2, 2, 2, 133, 21, 3, 2, 2, 2, 134, 135, 8, 12, 1, 2, 135, 136,
	7, 26, 2, 2, 136, 137, 5, 26, 14, 2, 137, 23, 3, 2, 2, 2, 138, 139, 8,
	13, 1, 2, 139, 140, 7, 26, 2, 2, 140, 144, 7, 4, 2, 2, 141, 143, 5, 28,
	15, 2, 142, 141, 3, 2, 2, 2, 143, 146, 3, 2, 2, 2, 144, 142, 3, 2, 2, 2,
	144, 145, 3, 2, 2, 2, 145, 147, 3, 2, 2, 2, 146, 144, 3, 2, 2, 2, 147,
	148, 7, 5, 2, 2, 148, 25, 3, 2, 2, 2, 149, 152, 5, 30, 16, 2, 150, 152,
	5, 32, 17, 2, 151, 149, 3, 2, 2, 2, 151, 150, 3, 2, 2, 2, 152, 27, 3, 2,
	2, 2, 153, 156, 5, 34, 18, 2, 154, 156, 5, 36, 19, 2, 155, 153, 3, 2, 2,
	2, 155, 154, 3, 2, 2, 2, 156, 29, 3, 2, 2, 2, 157, 158, 8, 16, 1, 2, 158,
	160, 7, 26, 2, 2, 159, 161, 7, 26, 2, 2, 160, 159, 3, 2, 2, 2, 160, 161,
	3, 2, 2, 2, 161, 162, 3, 2, 2, 2, 162, 166, 7, 6, 2, 2, 163, 165, 5, 38,
	20, 2, 164, 163, 3, 2, 2, 2, 165, 168, 3, 2, 2, 2, 166, 164, 3, 2, 2, 2,
	166, 167, 3, 2, 2, 2, 167, 169, 3, 2, 2, 2, 168, 166, 3, 2, 2, 2, 169,
	170, 7, 7, 2, 2, 170, 31, 3, 2, 2, 2, 171, 172, 8, 17, 1, 2, 172, 174,
	7, 26, 2, 2, 173, 175, 7, 3, 2, 2, 174, 173, 3, 2, 2, 2, 174, 175, 3, 2,
	2, 2, 175, 176, 3, 2, 2, 2, 176, 177, 5, 44, 23, 2, 177, 33, 3, 2, 2, 2,
	178, 179, 8, 18, 1, 2, 179, 181, 7, 26, 2, 2, 180, 182, 7, 26, 2, 2, 181,
	180, 3, 2, 2, 2, 181, 182, 3, 2, 2, 2, 182, 183, 3, 2, 2, 2, 183, 187,
	7, 6, 2, 2, 184, 186, 5, 38, 20, 2, 185, 184, 3, 2, 2, 2, 186, 189, 3,
	2, 2, 2, 187, 185, 3, 2, 2, 2, 187, 188, 3, 2, 2, 2, 188, 190, 3, 2, 2,
	2, 189, 187, 3, 2, 2, 2, 190, 191, 7, 7, 2, 2, 191, 35, 3, 2, 2, 2, 192,
	193, 8, 19, 1, 2, 193, 195, 7, 26, 2, 2, 194, 196, 7, 3, 2, 2, 195, 194,
	3, 2, 2, 2, 195, 196, 3, 2, 2, 2, 196, 197, 3, 2, 2, 2, 197, 198, 5, 44,
	23, 2, 198, 37, 3, 2, 2, 2, 199, 200, 6, 20, 2, 2, 200, 203, 5, 40, 21,
	2, 201, 203, 5, 42, 22, 2, 202, 199, 3, 2, 2, 2, 202, 201, 3, 2, 2, 2,
	203, 39, 3, 2, 2, 2, 204, 205, 8, 21, 1, 2, 205, 206, 7, 26, 2, 2, 206,
	208, 5, 44, 23, 2, 207, 209, 7, 24, 2, 2, 208, 207, 3, 2, 2, 2, 208, 209,
	3, 2, 2, 2, 209, 41, 3, 2, 2, 2, 210, 212, 7, 8, 2, 2, 211, 210, 3, 2,
	2, 2, 211, 212, 3, 2, 2, 2, 212, 213, 3, 2, 2, 2, 213, 214, 7, 26, 2, 2,
	214, 43, 3, 2, 2, 2, 215, 216, 8, 23, 1, 2, 216, 224, 7, 26, 2, 2, 217,
	224, 5, 48, 25, 2, 218, 224, 5, 50, 26, 2, 219, 224, 7, 18, 2, 2, 220,
	224, 7, 9, 2, 2, 221, 224, 5, 46, 24, 2, 222, 224, 5, 30, 16, 2, 223, 215,
	3, 2, 2, 2, 223, 217, 3, 2, 2, 2, 223, 218, 3, 2, 2, 2, 223, 219, 3, 2,
	2, 2, 223, 220, 3, 2, 2, 2, 223, 221, 3, 2, 2, 2, 223, 222, 3, 2, 2, 2,
	224, 45, 3, 2, 2, 2, 225, 226, 7, 8, 2, 2, 226, 227, 8, 24, 1, 2, 227,
	228, 7, 26, 2, 2, 228, 47, 3, 2, 2, 2, 229, 230, 8, 25, 1, 2, 230, 231,
	7, 26, 2, 2, 231, 232, 7, 10, 2, 2, 232, 233, 8, 25, 1, 2, 233, 234, 7,
	26, 2, 2, 234, 235, 7, 11, 2, 2, 235, 236, 5, 44, 23, 2, 236, 49, 3, 2,
	2, 2, 237, 238, 7, 10, 2, 2, 238, 239, 7, 11, 2, 2, 239, 240, 5, 44, 23,
	2, 240, 51, 3, 2, 2, 2, 241, 243, 5, 54, 28, 2, 242, 241, 3, 2, 2, 2, 242,
	243, 3, 2, 2, 2, 243, 244, 3, 2, 2, 2, 244, 245, 5, 56, 29, 2, 245, 53,
	3, 2, 2, 2, 246, 247, 7, 19, 2, 2, 247, 249, 7, 4, 2, 2, 248, 250, 5, 70,
	36, 2, 249, 248, 3, 2, 2, 2, 250, 251, 3, 2, 2, 2, 251, 249, 3, 2, 2, 2,
	251, 252, 3, 2, 2, 2, 252, 253, 3, 2, 2, 2, 253, 254, 7, 5, 2, 2, 254,
	55, 3, 2, 2, 2, 255, 256, 8, 29, 1, 2, 256, 257, 7, 26, 2, 2, 257, 258,
	5, 72, 37, 2, 258, 262, 7, 6, 2, 2, 259, 261, 5, 58, 30, 2, 260, 259, 3,
	2, 2, 2, 261, 264, 3, 2, 2, 2, 262, 260, 3, 2, 2, 2, 262, 263, 3, 2, 2,
	2, 263, 265, 3, 2, 2, 2, 264, 262, 3, 2, 2, 2, 265, 266, 7, 7, 2, 2, 266,
	57, 3, 2, 2, 2, 267, 269, 5, 60, 31, 2, 268, 267, 3, 2, 2, 2, 268, 269,
	3, 2, 2, 2, 269, 272, 3, 2, 2, 2, 270, 273, 5, 54, 28, 2, 271, 273, 5,
	62, 32, 2, 272, 270, 3, 2, 2, 2, 272, 271, 3, 2, 2, 2, 273, 274, 3, 2,
	2, 2, 274, 275, 5, 64, 33, 2, 275, 59, 3, 2, 2, 2, 276, 278, 7, 16, 2,
	2, 277, 279, 7, 4, 2, 2, 278, 277, 3, 2, 2, 2, 278, 279, 3, 2, 2, 2, 279,
	286, 3, 2, 2, 2, 280, 282, 5, 70, 36, 2, 281, 280, 3, 2, 2, 2, 282, 283,
	3, 2, 2, 2, 283, 281, 3, 2, 2, 2, 283, 284, 3, 2, 2, 2, 284, 287, 3, 2,
	2, 2, 285, 287, 7, 23, 2, 2, 286, 281, 3, 2, 2, 2, 286, 285, 3, 2, 2, 2,
	287, 289, 3, 2, 2, 2, 288, 290, 7, 5, 2, 2, 289, 288, 3, 2, 2, 2, 289,
	290, 3, 2, 2, 2, 290, 61, 3, 2, 2, 2, 291, 292, 7, 17, 2, 2, 292, 293,
	7, 26, 2, 2, 293, 63, 3, 2, 2, 2, 294, 295, 8, 33, 1, 2, 295, 296, 7, 26,
	2, 2, 296, 298, 5, 74, 38, 2, 297, 299, 5, 66, 34, 2, 298, 297, 3, 2, 2,
	2, 298, 299, 3, 2, 2, 2, 299, 301, 3, 2, 2, 2, 300, 302, 5, 68, 35, 2,
	301, 300, 3, 2, 2, 2, 301, 302, 3, 2, 2, 2, 302, 65, 3, 2, 2, 2, 303, 305,
	7, 4, 2, 2, 304, 306, 7, 26, 2, 2, 305, 304, 3, 2, 2, 2, 305, 306, 3, 2,
	2, 2, 306, 307, 3, 2, 2, 2, 307, 308, 7, 5, 2, 2, 308, 67, 3, 2, 2, 2,
	309, 310, 7, 12, 2, 2, 310, 312, 7, 4, 2, 2, 311, 313, 5, 44, 23, 2, 312,
	311, 3, 2, 2, 2, 312, 313, 3, 2, 2, 2, 313, 314, 3, 2, 2, 2, 314, 315,
	7, 5, 2, 2, 315, 69, 3, 2, 2, 2, 316, 317, 7, 26, 2, 2, 317, 318, 8, 36,
	1, 2, 318, 319, 7, 25, 2, 2, 319, 71, 3, 2, 2, 2, 320, 322, 7, 26, 2, 2,
	321, 323, 7, 13, 2, 2, 322, 321, 3, 2, 2, 2, 322, 323, 3, 2, 2, 2, 323,
	325, 3, 2, 2, 2, 324, 320, 3, 2, 2, 2, 325, 326, 3, 2, 2, 2, 326, 324,
	3, 2, 2, 2, 326, 327, 3, 2, 2, 2, 327, 73, 3, 2, 2, 2, 328, 329, 7, 14,
	2, 2, 329, 334, 5, 76, 39, 2, 330, 331, 7, 13, 2, 2, 331, 333, 5, 76, 39,
	2, 332, 330, 3, 2, 2, 2, 333, 336, 3, 2, 2, 2, 334, 332, 3, 2, 2, 2, 334,
	335, 3, 2, 2, 2, 335, 344, 3, 2, 2, 2, 336, 334, 3, 2, 2, 2, 337, 338,
	7, 15, 2, 2, 338, 341, 5, 76, 39, 2, 339, 340, 7, 13, 2, 2, 340, 342, 5,
	76, 39, 2, 341, 339, 3, 2, 2, 2, 341, 342, 3, 2, 2, 2, 342, 344, 3, 2,
	2, 2, 343, 328, 3, 2, 2, 2, 343, 337, 3, 2, 2, 2, 344, 345, 3, 2, 2, 2,
	345, 343, 3, 2, 2, 2, 345, 346, 3, 2, 2, 2, 346, 349, 3, 2, 2, 2, 347,
	349, 7, 14, 2, 2, 348, 343, 3, 2, 2, 2, 348, 347, 3, 2, 2, 2, 349, 75,
	3, 2, 2, 2, 350, 352, 9, 2, 2, 2, 351, 350, 3, 2, 2, 2, 352, 353, 3, 2,
	2, 2, 353, 351, 3, 2, 2, 2, 353, 354, 3, 2, 2, 2, 354, 77, 3, 2, 2, 2,
	42, 81, 89, 99, 111, 126, 132, 144, 151, 155, 160, 166, 174, 181, 187,
	195, 202, 208, 211, 223, 242, 251, 262, 268, 272, 278, 283, 286, 289, 298,
	301, 305, 312, 322, 326, 334, 341, 343, 345, 348, 353,
}
var literalNames = []string{
	"", "'='", "'('", "')'", "'{'", "'}'", "'*'", "'time.Time'", "'['", "']'",
	"'returns'", "'-'", "'/'", "'/:'", "'@doc'", "'@handler'", "'interface{}'",
	"'@server'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "ATDOC", "ATHANDLER",
	"INTERFACE", "ATSERVER", "WS", "COMMENT", "LINE_COMMENT", "STRING", "RAW_STRING",
	"LINE_VALUE", "ID", "LetterOrDigit",
}

var ruleNames = []string{
	"api", "spec", "syntaxLit", "importSpec", "importLit", "importBlock", "importBlockValue",
	"importValue", "infoSpec", "typeSpec", "typeLit", "typeBlock", "typeLitBody",
	"typeBlockBody", "typeStruct", "typeAlias", "typeBlockStruct", "typeBlockAlias",
	"field", "normalField", "anonymousFiled", "dataType", "pointerType", "mapType",
	"arrayType", "serviceSpec", "atServer", "serviceApi", "serviceRoute", "atDoc",
	"atHandler", "route", "body", "replybody", "kvLit", "serviceName", "path",
	"pathItem",
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
	ApiParserParserWS            = 18
	ApiParserParserCOMMENT       = 19
	ApiParserParserLINE_COMMENT  = 20
	ApiParserParserSTRING        = 21
	ApiParserParserRAW_STRING    = 22
	ApiParserParserLINE_VALUE    = 23
	ApiParserParserID            = 24
	ApiParserParserLetterOrDigit = 25
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
	ApiParserParserRULE_pathItem         = 37
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

func (s *ApiContext) Accept(visitor antlr.ParseTreeVisitor) any {
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
	p.SetState(79)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == ApiParserParserATSERVER || _la == ApiParserParserID {
		{
			p.SetState(76)
			p.Spec()
		}

		p.SetState(81)
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

func (s *SpecContext) Accept(visitor antlr.ParseTreeVisitor) any {
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

	p.SetState(87)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(82)
			p.SyntaxLit()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(83)
			p.ImportSpec()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(84)
			p.InfoSpec()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(85)
			p.TypeSpec()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(86)
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
