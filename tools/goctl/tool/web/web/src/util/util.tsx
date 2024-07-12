
const hash = (s: string) => {
    let h = 1315423911, i, ch
    for (i = s.length - 1; i >= 0; i--) {
        ch = s.charCodeAt(i)
        h ^= ((h << 5) + ch + (h >> 2))
    }
    return (h & 0x7FFFFFFF)
}

export const getColor = (s: string) => {
    return "#" + hash(s).toString().substring(0, 6)
}


export const formatKey = (s1: string, s2: string) => {
    return s1 + '_' + s2
}

