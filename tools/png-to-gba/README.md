# png-to-gba

Convert PNG to GameBoy Advance images

Create a paletted PNG using Gimp:
 - Create a new image 240 x 160 pixels (to fill the whole screen) or smaller
 - Go to Image > Mode > Indexed
   - Select "Generate optimum palette"
   - Set "Maximum number of colours" to 256
   - Click "Convert"
 - When using the colour selection dialog, you will be able to choose colours
   and their index. Make sure the background is index 0

There is an example image ([PNG](img.png), [XCF](img.xcf) (Gimp)) that
is 5x1 with black background and each pixel is one of the colours:
Black, Red, Blue, Green and White.
