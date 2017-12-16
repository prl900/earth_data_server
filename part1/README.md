# earth_data_server (part 1)

This code is referenced on the Medium blog post: https://medium.com/@p.rozas.larraondo/divide-compress-and-conquer-building-an-earth-data-server-in-go-part-1-d82eee2eceb1

To run this code you'll need a Go 1.0+ environment set in your system.
In Linux/Mac OSX do:

1.- Download the Blue Marble Next Generation file on this folder:

`$ curl -O "https://eoimages.gsfc.nasa.gov/images/imagerecords/73000/73909/world.topo.bathy.200412.3x21600x10800.png"`

2.- Generate the tiles for this file:

`$ go run generate_tiles.go`

3.- Request a region providing the coordinates of any place in the world. The result will be saved on "output.png"

`$ time go run get_region.go -lat 42 -lon -1`

4.- Same operation but using the generated tiles

`$ time go run get_region_tiles.go -lat 42 -lon -1`
