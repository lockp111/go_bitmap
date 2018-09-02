# README

## Usage
### NewBitmap
```golang
bm, err := NewBitmap(<size> int)
```

### NewBitmapMax
```golang
bm, err := NewBitmapMax()
```

### Add
```golang
bm.Add((<num> uint64)
```

### Del
```golang
bm.Del((<num> uint64)
```

### Has
```golang
if bm.Has(123) {
    fmt.Println("has: ", 123)
}
```

### Get Maxpos
```golang
maxpos := bm.Maxpos()
```

### Show all has (top 100)
```golang
bm.String()
```

### Has next pos
```golang
next, has := bm.Next(<pos> uint64)
if !has{
    fmt.Printf("pos is max")
}
```

### Has prev pos
```golang
prev, has := bm.Prev(<pos> uint64)
if !has{
    fmt.Printf("pos is min")
}
```