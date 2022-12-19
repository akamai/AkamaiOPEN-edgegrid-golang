package imaging

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// TransformationType is implemented by Transformation types
	TransformationType interface {
		transformationType() string
	}

	// TransformationTypePost is implemented by PostBreakpointTransformations types
	TransformationTypePost interface {
		transformationTypePost() string
	}

	// Transformations represents an array of Transformations
	Transformations []TransformationType

	// PostBreakpointTransformations represents an array of PostBreakPointTransformations
	PostBreakpointTransformations []TransformationTypePost

	// ImageType is implemented by ImageType types
	ImageType interface {
		imageType() string
	}

	// ImageTypePost is implemented by ImageTypePost types
	ImageTypePost interface {
		imageTypePost() string
	}

	// ShapeType is implemented by ImageType types
	ShapeType interface {
		shapeType() string
	}

	// InlineVariable represents an inline variable
	InlineVariable struct {
		Var string `json:"var"`
	}

	/*-----------------------------------------------*/
	///////////////// Generated types /////////////////
	/*-----------------------------------------------*/

	// Append Places a specified `image` beside the source image. The API places the `image` on a major dimension, then aligns it on the minor dimension. Transparent pixels fill any area not covered by either image.
	Append struct {
		// Gravity Specifies where to place the `image` relative to the source image. The available values represent the eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and a `Center` by default.
		Gravity *GravityVariableInline `json:"gravity,omitempty"`
		// GravityPriority Determines the exact placement of the `image` when `gravity` is `Center` or a diagonal. The value is either `horizontal` or `vertical`. Use `horizontal` to append an `image` east or west of the source image. This aligns the `image` on the vertical gravity component, placing `Center` gravity east. Use `vertical` to append an `image` north or south of the source image. This aligns the `image` on the horizontal gravity component, placing `Center` gravity south.
		GravityPriority *AppendGravityPriorityVariableInline `json:"gravityPriority,omitempty"`
		Image           ImageType                            `json:"image"`
		// PreserveMinorDimension Whether to preserve the source image's minor dimension, `false` by default. The minor dimension is the dimension opposite the dimension that the appending `image` is placed. For example, the dimensions of the source image are 100 &times; 100 pixels. The dimensions of the appending `image` are 50 &times; 150 pixels. The `gravity` is set to `East`. This makes the major dimension horizontal and the source image's minor dimension vertical. To preserve the source image's minor dimension at 100 pixels, the `preserveMinorDimension` is set to `true`. As a result of the append, the major dimension expanded with the appended image to 150 pixels. The source image's minor dimension was maintained at 100 pixels. The total combined dimension of the image is 150 &times; 100 pixels.
		PreserveMinorDimension *BooleanVariableInline `json:"preserveMinorDimension,omitempty"`
		// Transformation Identifies this type of transformation, `Append` in this case.
		Transformation AppendTransformation `json:"transformation"`
	}

	// AppendGravityPriority ...
	AppendGravityPriority string

	// AppendGravityPriorityVariableInline represents a type which stores either a value or a variable name
	AppendGravityPriorityVariableInline struct {
		Name  *string
		Value *AppendGravityPriority
	}

	// AppendTransformation ...
	AppendTransformation string

	// AspectCrop Lets you change the height or width of an image (either by cropping or expanding the area) to an aspect ratio of your choosing.
	AspectCrop struct {
		// AllowExpansion Increases the size of the image canvas to achieve the requested aspect ratio instead of cropping the image. Use the Horizontal Offset and Vertical Offset settings to determine where to add the fully transparent pixels on the expanded image canvas.
		AllowExpansion *BooleanVariableInline `json:"allowExpansion,omitempty"`
		// Height The height term of the aspect ratio to crop.
		Height *NumberVariableInline `json:"height,omitempty"`
		// Transformation Identifies this type of transformation, `AspectCrop` in this case.
		Transformation AspectCropTransformation `json:"transformation"`
		// Width The width term of the aspect ratio to crop.
		Width *NumberVariableInline `json:"width,omitempty"`
		// XPosition Specifies the horizontal portion of the image you want to keep when the aspect ratio cropping is applied. When using Allow Expansion this setting defines the horizontal position of the image on the new expanded image canvas.
		XPosition *NumberVariableInline `json:"xPosition,omitempty"`
		// YPosition Specifies the horizontal portion of the image you want to keep when the aspect ratio cropping is applied. When using Allow Expansion this setting defines the horizontal position of the image on the new expanded image canvas.
		YPosition *NumberVariableInline `json:"yPosition,omitempty"`
	}

	// AspectCropTransformation ...
	AspectCropTransformation string

	// BackgroundColor Places a transparent image on a set background color. Color is specified in the typical CSS hexadecimal format.
	BackgroundColor struct {
		// Color The hexadecimal CSS color value for the background.
		Color *StringVariableInline `json:"color"`
		// Transformation Identifies this type of transformation, `BackgroundColor` in this case.
		Transformation BackgroundColorTransformation `json:"transformation"`
	}

	// BackgroundColorTransformation ...
	BackgroundColorTransformation string

	// Blur Applies a Gaussian blur to the image.
	Blur struct {
		// Sigma The number of pixels to scatter the original pixel by to create the blur effect. Resulting images may be larger than the original as pixels at the edge of the image might scatter outside the image's original dimensions.
		Sigma *NumberVariableInline `json:"sigma,omitempty"`
		// Transformation Identifies this type of transformation, `Blur` in this case.
		Transformation BlurTransformation `json:"transformation"`
	}

	// BlurTransformation ...
	BlurTransformation string

	// BooleanVariableInline represents a type which stores either a value or a variable name
	BooleanVariableInline struct {
		Name  *string
		Value *bool
	}

	// BoxImageType A rectangular box, with a specified color and applied transformation.
	BoxImageType struct {
		// Color The fill color of the box, not the edge of the box. The API supports hexadecimal representation and CSS hexadecimal color values.
		Color *StringVariableInline `json:"color,omitempty"`
		// Height The height of the box in pixels.
		Height         *IntegerVariableInline `json:"height,omitempty"`
		Transformation TransformationType     `json:"transformation,omitempty"`
		// Type Identifies the type of image, `Box` in this case.
		Type BoxImageTypeType `json:"type"`
		// Width The width of the box in pixels.
		Width *IntegerVariableInline `json:"width,omitempty"`
	}

	// BoxImageTypePost A rectangular box, with a specified color and applied transformation.
	BoxImageTypePost struct {
		// Color The fill color of the box, not the edge of the box. The API supports hexadecimal representation and CSS hexadecimal color values.
		Color *StringVariableInline `json:"color,omitempty"`
		// Height The height of the box in pixels.
		Height         *IntegerVariableInline `json:"height,omitempty"`
		Transformation TransformationTypePost `json:"transformation,omitempty"`
		// Type Identifies the type of image, `Box` in this case.
		Type BoxImageTypePostType `json:"type"`
		// Width The width of the box in pixels.
		Width *IntegerVariableInline `json:"width,omitempty"`
	}

	// BoxImageTypePostType ...
	BoxImageTypePostType string

	// BoxImageTypeType ...
	BoxImageTypeType string

	// Breakpoints The breakpoint widths (in pixels) to use to create derivative images/videos.
	Breakpoints struct {
		Widths []int `json:"widths,omitempty"`
	}

	// ChromaKey Changes any color in an image within the specified volume of the HSL colorspace to transparent or semitransparent. This transformation applies a 'green screen' technique commonly used to isolate and remove background colors.
	ChromaKey struct {
		// Hue The hue to remove. Enter the degree of rotation between 0 and 360 degrees around the color wheel. By default Chroma Key removes a green hue, 120° on the color wheel.
		Hue *NumberVariableInline `json:"hue,omitempty"`
		// HueFeather How much additional hue to make semi-transparent beyond the Hue Tolerance. By default Hue Feather is 0.083 which applies semi-transparency to hues 30° around the Hue Tolerance.
		HueFeather *NumberVariableInline `json:"hueFeather,omitempty"`
		// HueTolerance How close a color's hue needs to be to the selected hue for it to be changed to fully transparent. If you enter the maximum value of 1.0 the entire image is made transparent. By default Hue Tolerance is approximately 0.083 or 8.3% of the color wheel. This value corresponds to 30° around the specified hue.
		HueTolerance *NumberVariableInline `json:"hueTolerance,omitempty"`
		// LightnessFeather How much additional lightness to make semi-transparent beyond the Lightness Tolerance. The default value of 0.1 corresponds to 10% away from the tolerated lightness towards full black or full white.
		LightnessFeather *NumberVariableInline `json:"lightnessFeather,omitempty"`
		// LightnessTolerance How much of the lightest part and darkest part of a color to preserve. For example, you can space this value out from the middle (i.e. 0.5 lightness or full color) to help preserve the splash lighting impact in the image. You can define how close the color needs to be to the full color to remove it from your image. The default value of 0.75 means that a colour must be within 75% of the full color to full white or full black for full removal.
		LightnessTolerance *NumberVariableInline `json:"lightnessTolerance,omitempty"`
		// SaturationFeather How much additional saturation to make semi-transparent beyond the Saturation Tolerance. By default Saturation Feather is 0.1 which applies semi-transparency to hues 10% below the saturationTolerance.
		SaturationFeather *NumberVariableInline `json:"saturationFeather,omitempty"`
		// SaturationTolerance How close a color's saturation needs to be to full saturation for it to be changed to fully transparent. For example, you can define how green the color needs to be to remove it from your image. The default value of 0.75 means that a color must be within 75% of full saturation in order to be made fully transparent.
		SaturationTolerance *NumberVariableInline `json:"saturationTolerance,omitempty"`
		// Transformation Identifies this type of transformation, `ChromaKey` in this case.
		Transformation ChromaKeyTransformation `json:"transformation"`
	}

	// ChromaKeyTransformation ...
	ChromaKeyTransformation string

	// CircleImageType A rectangular box, with a specified color and applied transformation.
	CircleImageType struct {
		// Color The fill color of the circle. The API supports hexadecimal representation and CSS hexadecimal color values.
		Color *StringVariableInline `json:"color,omitempty"`
		// Diameter The diameter of the circle. The diameter will be the width and the height of the image in pixels.
		Diameter       *IntegerVariableInline `json:"diameter,omitempty"`
		Transformation TransformationType     `json:"transformation,omitempty"`
		// Type Identifies the type of image, `Circle` in this case.
		Type CircleImageTypeType `json:"type"`
		// Width The width of the box in pixels.
		Width *IntegerVariableInline `json:"width,omitempty"`
	}

	// CircleImageTypePost A rectangular box, with a specified color and applied transformation.
	CircleImageTypePost struct {
		// Color The fill color of the circle. The API supports hexadecimal representation and CSS hexadecimal color values.
		Color *StringVariableInline `json:"color,omitempty"`
		// Diameter The diameter of the circle. The diameter will be the width and the height of the image in pixels.
		Diameter       *IntegerVariableInline `json:"diameter,omitempty"`
		Transformation TransformationTypePost `json:"transformation,omitempty"`
		// Type Identifies the type of image, `Circle` in this case.
		Type CircleImageTypePostType `json:"type"`
		// Width The width of the box in pixels.
		Width *IntegerVariableInline `json:"width,omitempty"`
	}

	// CircleImageTypePostType ...
	CircleImageTypePostType string

	// CircleImageTypeType ...
	CircleImageTypeType string

	// CircleShapeType Defines a circle with a specified `radius` from its `center` point.
	CircleShapeType struct {
		// Center Defines coordinates for a single point, to help define polygons and rectangles. Each point may be an object with `x`and `y` members, or a two-element array.
		Center *PointShapeType `json:"center"`
		// Radius The radius of the circle measured in pixels.
		Radius *NumberVariableInline `json:"radius"`
	}

	// Composite Applies another image to the source image, either as an overlay or an underlay. The image that's underneath is visible in areas that are beyond the edges of the top image or that are less than 100% opaque. A common use of an overlay composite is to add a watermark.
	Composite struct {
		// Gravity Compass direction indicating the corner or edge of the base image to place the applied image. Use Horizontal and Vertical Offset to adjust the applied image's gravity position
		Gravity *GravityVariableInline `json:"gravity,omitempty"`
		Image   ImageType              `json:"image"`
		// Placement Place applied image on top of or underneath the base image. Watermarks are usually applied over. Backgrounds are usually applied under.
		Placement *CompositePlacementVariableInline `json:"placement,omitempty"`
		// Scale A multiplier to resize the applied image relative to the source image and preserve aspect ratio, 1 by default. Set the `scaleDimension` to calculate the `scale` from the source image's width or height.
		Scale *NumberVariableInline `json:"scale,omitempty"`
		// ScaleDimension The dimension, either `width` or `height`, of the source image to scale.
		ScaleDimension *CompositeScaleDimensionVariableInline `json:"scaleDimension,omitempty"`
		// Transformation Identifies this type of transformation, `Composite` in this case.
		Transformation CompositeTransformation `json:"transformation"`
		// XPosition The x-axis position of the image to apply.
		XPosition *IntegerVariableInline `json:"xPosition,omitempty"`
		// YPosition The y-axis position of the image to apply.
		YPosition *IntegerVariableInline `json:"yPosition,omitempty"`
	}

	// CompositePlacement ...
	CompositePlacement string

	// CompositePlacementVariableInline represents a type which stores either a value or a variable name
	CompositePlacementVariableInline struct {
		Name  *string
		Value *CompositePlacement
	}

	// CompositePost Applies another image to the source image, either as an overlay or an underlay. The image that's underneath is visible in areas that are beyond the edges of the top image or that are less than 100% opaque. A common use of an overlay composite is to add a watermark.
	CompositePost struct {
		// Gravity Compass direction indicating the corner or edge of the base image to place the applied image. Use Horizontal and Vertical Offset to adjust the applied image's gravity position
		Gravity *GravityPostVariableInline `json:"gravity,omitempty"`
		Image   ImageTypePost              `json:"image"`
		// Placement Place applied image on top of or underneath the base image. Watermarks are usually applied over. Backgrounds are usually applied under.
		Placement *CompositePostPlacementVariableInline `json:"placement,omitempty"`
		// Scale A multiplier to resize the applied image relative to the source image and preserve aspect ratio, 1 by default. Set the `scaleDimension` to calculate the `scale` from the source image's width or height.
		Scale *NumberVariableInline `json:"scale,omitempty"`
		// ScaleDimension The dimension, either `width` or `height`, of the source image to scale.
		ScaleDimension *CompositePostScaleDimensionVariableInline `json:"scaleDimension,omitempty"`
		// Transformation Identifies this type of transformation, `Composite` in this case.
		Transformation CompositePostTransformation `json:"transformation"`
		// XPosition The x-axis position of the image to apply.
		XPosition *IntegerVariableInline `json:"xPosition,omitempty"`
		// YPosition The y-axis position of the image to apply.
		YPosition *IntegerVariableInline `json:"yPosition,omitempty"`
	}

	// CompositePostPlacement ...
	CompositePostPlacement string

	// CompositePostPlacementVariableInline represents a type which stores either a value or a variable name
	CompositePostPlacementVariableInline struct {
		Name  *string
		Value *CompositePostPlacement
	}

	// CompositePostScaleDimension ...
	CompositePostScaleDimension string

	// CompositePostScaleDimensionVariableInline represents a type which stores either a value or a variable name
	CompositePostScaleDimensionVariableInline struct {
		Name  *string
		Value *CompositePostScaleDimension
	}

	// CompositePostTransformation ...
	CompositePostTransformation string

	// CompositeScaleDimension ...
	CompositeScaleDimension string

	// CompositeScaleDimensionVariableInline represents a type which stores either a value or a variable name
	CompositeScaleDimensionVariableInline struct {
		Name  *string
		Value *CompositeScaleDimension
	}

	// CompositeTransformation ...
	CompositeTransformation string

	// Compound ...
	Compound struct {
		// Transformation Identifies this type of transformation, `Compound` in this case.
		Transformation  CompoundTransformation `json:"transformation"`
		Transformations Transformations        `json:"transformations,omitempty"`
	}

	// CompoundPost ...
	CompoundPost struct {
		// Transformation Identifies this type of transformation, `Compound` in this case.
		Transformation  CompoundPostTransformation    `json:"transformation"`
		Transformations PostBreakpointTransformations `json:"transformations,omitempty"`
	}

	// CompoundPostTransformation ...
	CompoundPostTransformation string

	// CompoundTransformation ...
	CompoundTransformation string

	// Contrast Adjusts both the contrast and brightness of an image.
	Contrast struct {
		// Brightness Adjusts the brightness of the image. Positive values increase brightness and negative values decrease brightness. A value of  `1` produces a white image. A value of  `-1` produces a black image. The default value is `0`, which leaves the image unchanged. The acceptable value range is `-1.0` to `1.0`. Values outside of the acceptable range clamp to this range.
		Brightness *NumberVariableInline `json:"brightness,omitempty"`
		// Contrast Adjusts the contrast of the image. Expressed as a range from `-1` to `1`, positive values increase contrast, negative values decrease it, while `0` leaves the image unchanged. Values outside of the `-1` to `1` range clamp to this range.
		Contrast *NumberVariableInline `json:"contrast,omitempty"`
		// Transformation Identifies this type of transformation, `Contrast` in this case.
		Transformation ContrastTransformation `json:"transformation"`
	}

	// ContrastTransformation ...
	ContrastTransformation string

	// Crop Crops an image.
	Crop struct {
		// AllowExpansion If cropping an area outside of the existing canvas, expands the image canvas.
		AllowExpansion *BooleanVariableInline `json:"allowExpansion,omitempty"`
		// Gravity Frame of reference for X and Y Positions.
		Gravity *GravityVariableInline `json:"gravity,omitempty"`
		// Height The number of pixels to crop along the y-axis.
		Height *IntegerVariableInline `json:"height"`
		// Transformation Identifies this type of transformation, `Crop` in this case.
		Transformation CropTransformation `json:"transformation"`
		// Width The number of pixels to crop along the x-axis.
		Width *IntegerVariableInline `json:"width"`
		// XPosition The x-axis position of the image to crop from.
		XPosition *IntegerVariableInline `json:"xPosition,omitempty"`
		// YPosition The y-axis position of the image to crop from.
		YPosition *IntegerVariableInline `json:"yPosition,omitempty"`
	}

	// CropTransformation ...
	CropTransformation string

	// EnumOptions Optionally limits the set of possible values for a variable. References to an enum `id` insert a corresponding `value`.
	EnumOptions struct {
		// ID The unique identifier for each enum value, up to 50 alphanumeric characters.
		ID string `json:"id"`
		// Value The value of the variable when the `id` is provided.
		Value string `json:"value"`
	}

	// FaceCrop Applies a method to detect faces in the source image and applies the rectangular crop on either the `biggest` face or `all` of the faces detected. Image and Video Manager tries to preserve faces in the image instead of using specified crop coordinates.
	FaceCrop struct {
		// Algorithm Specifies the type of algorithm used to detect faces in the image, either `cascade` for the cascade classifier algorithm or `dnn` for the deep neural network algorithm, `cascade` by default.
		Algorithm *FaceCropAlgorithmVariableInline `json:"algorithm,omitempty"`
		// Confidence With `algorithm` set to `dnn`, specifies the minimum confidence needed to detect faces in the image. Values range from `0` to `1` for increased confidence, and possibly fewer faces detected.
		Confidence *NumberVariableInline `json:"confidence,omitempty"`
		// FailGravity Controls placement of the crop if Image and Video Manager does not detect any faces in the image. Directions are relative to the edges of the image being transformed.
		FailGravity *GravityVariableInline `json:"failGravity,omitempty"`
		// Focus Distinguishes the faces detected, either `biggestFace` or `allFaces` to place the crop rectangle around the full set of faces, `all` by default.
		Focus *FaceCropFocusVariableInline `json:"focus,omitempty"`
		// Gravity Controls placement of the crop. Directions are relative to the face(s) plus padding.
		Gravity *GravityVariableInline `json:"gravity,omitempty"`
		// Height The height of the output image in pixels relative to the specified `style` value.
		Height *IntegerVariableInline `json:"height"`
		// Padding The padding ratio based on the dimensions of the biggest face detected, `0.5` by default. Larger values increase padding.
		Padding *NumberVariableInline `json:"padding,omitempty"`
		// Style Specifies how to crop or scale a crop area for the faces detected in the source image, `zoom` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` places a raw crop around the faces, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and faces as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the faces, relative to the specified `width` and `height` values. Allows Variable substitution.
		Style *FaceCropStyleVariableInline `json:"style,omitempty"`
		// Transformation Identifies this type of transformation, `FaceCrop` in this case.
		Transformation FaceCropTransformation `json:"transformation"`
		// Width The width of the output image in pixels relative to the specified `style` value.
		Width *IntegerVariableInline `json:"width"`
	}

	// FaceCropAlgorithm ...
	FaceCropAlgorithm string

	// FaceCropAlgorithmVariableInline represents a type which stores either a value or a variable name
	FaceCropAlgorithmVariableInline struct {
		Name  *string
		Value *FaceCropAlgorithm
	}

	// FaceCropFocus ...
	FaceCropFocus string

	// FaceCropFocusVariableInline represents a type which stores either a value or a variable name
	FaceCropFocusVariableInline struct {
		Name  *string
		Value *FaceCropFocus
	}

	// FaceCropStyle ...
	FaceCropStyle string

	// FaceCropStyleVariableInline represents a type which stores either a value or a variable name
	FaceCropStyleVariableInline struct {
		Name  *string
		Value *FaceCropStyle
	}

	// FaceCropTransformation ...
	FaceCropTransformation string

	// FeatureCrop Identifies prominent features of the source image, then crops around as many of these features as possible relative to the specified `width` and `height` values.
	FeatureCrop struct {
		// FailGravity Controls placement of the crop if Image and Video Manager does not detect any features in the image. Directions are relative to the edges of the image being transformed.
		FailGravity *GravityVariableInline `json:"failGravity,omitempty"`
		// FeatureRadius The size in pixels of the important features to search for. If identified, two features never appear closer together than this value, `8.0` by default.
		FeatureRadius *NumberVariableInline `json:"featureRadius,omitempty"`
		// Gravity Controls placement of the crop. Directions are relative to the region of interest plus padding.
		Gravity *GravityVariableInline `json:"gravity,omitempty"`
		// Height The height in pixels of the output image relative to the specified `style` value.
		Height *IntegerVariableInline `json:"height"`
		// MaxFeatures The maximum number of features to identify as important features, `32` by default. The strongest features are always chosen.
		MaxFeatures *IntegerVariableInline `json:"maxFeatures,omitempty"`
		// MinFeatureQuality Determines the minimum quality level of the feature identified. To consider a feature important, the feature needs to surpass this value.  Image and Video Manager measures quality on a scale from `0` for the lowest quality to `1` for the highest quality, `.1` by default.
		MinFeatureQuality *NumberVariableInline `json:"minFeatureQuality,omitempty"`
		// Padding Adds space around the region of interest. The amount of padding added is directly related to the size of the bounding box of the selected features. Specifically, the region of interest is expanded in all directions by the largest dimension of the bounding box of the selected features multiplied by this value.
		Padding *NumberVariableInline `json:"padding,omitempty"`
		// Style Specifies how to crop or scale a crop area for the features identified in the source image, `fill` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` performs a raw crop around the features, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and features as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the features, relative to the specified `width` and `height` values. Allows Variable substitution.
		Style *FeatureCropStyleVariableInline `json:"style,omitempty"`
		// Transformation Identifies this type of transformation, `FeatureCrop` in this case.
		Transformation FeatureCropTransformation `json:"transformation"`
		// Width The width in pixels of the output image relative to the specified `style` value.
		Width *IntegerVariableInline `json:"width"`
	}

	// FeatureCropStyle ...
	FeatureCropStyle string

	// FeatureCropStyleVariableInline represents a type which stores either a value or a variable name
	FeatureCropStyleVariableInline struct {
		Name  *string
		Value *FeatureCropStyle
	}

	// FeatureCropTransformation ...
	FeatureCropTransformation string

	// FitAndFill Resizes an image to fit within a specific size box and then uses a fill of that same image to cover any transparent space at the edges. By default the fill image has a Blur transformation with a sigma value of 8 applied, but the transformation can be customized using the fillTransformation parameter.
	FitAndFill struct {
		FillTransformation TransformationType `json:"fillTransformation,omitempty"`
		// Height The height value of the resized image.
		Height *IntegerVariableInline `json:"height"`
		// Transformation Identifies this type of transformation, `FitAndFill` in this case.
		Transformation FitAndFillTransformation `json:"transformation"`
		// Width The width value of the resized image.
		Width *IntegerVariableInline `json:"width"`
	}

	// FitAndFillTransformation ...
	FitAndFillTransformation string

	// Goop Distorts an image by randomly repositioning a set of control points along a specified grid. The transformed image appears _goopy_. Adjust the density of the grid and the degree of randomity. You can use this transformation to create watermarks for use in security.
	Goop struct {
		// Chaos Specifies the greatest distance control points may move from their original position. A value of `1.0` shifts control points over as far as the next one in the original grid. A value of `0.0` leaves the image unchanged. Values under `0.5` work better for subtle distortions, otherwise control points may pass each other and cause a twisting effect.
		Chaos *NumberVariableInline `json:"chaos,omitempty"`
		// Density Controls the density of control points used to distort the image. The largest dimension of the input image is divided up to fit this number of control points. A grid of points is extended on the smaller dimension such that each row and column of control points is equidistant from each adjacent row or column. This parameter strongly affects transformation performance. Be careful choosing values above the default if you expect to transform medium to large size images.
		Density *IntegerVariableInline `json:"density,omitempty"`
		// Power By default, the distortion algorithm relies on inverse squares to calculate distance but this allows you to change the exponent. You shouldnt need to vary the default value of `2.0`.
		Power *NumberVariableInline `json:"power,omitempty"`
		// Seed Specifies your own `seed` value as an alternative to the default, which is subject to variability. This allows for reproducible and deterministic distortions. If all parameters are kept equal and a constant seed is used, `Goop` distorts an input image consistently over many transformations. By default, this value is set to the current Epoch Time measured in milliseconds, which provides inconsistent transformation output.
		Seed *IntegerVariableInline `json:"seed,omitempty"`
		// Transformation Identifies this type of transformation, `Goop` in this case.
		Transformation GoopTransformation `json:"transformation"`
	}

	// GoopTransformation ...
	GoopTransformation string

	// Gravity ...
	Gravity string

	// GravityPost ...
	GravityPost string

	// GravityPostVariableInline represents a type which stores either a value or a variable name
	GravityPostVariableInline struct {
		Name  *string
		Value *GravityPost
	}

	// GravityVariableInline represents a type which stores either a value or a variable name
	GravityVariableInline struct {
		Name  *string
		Value *Gravity
	}

	// Grayscale Restricts image color to shades of gray only.
	Grayscale struct {
		// Transformation Identifies this type of transformation, `Grayscale` in this case.
		Transformation GrayscaleTransformation `json:"transformation"`
		// Type The algorithm used to transform colors to grays, either `Brightness`, `Lightness`, `Rec601`, or the default `Rec709`.
		Type *GrayscaleTypeVariableInline `json:"type,omitempty"`
	}

	// GrayscaleTransformation ...
	GrayscaleTransformation string

	// GrayscaleType ...
	GrayscaleType string

	// GrayscaleTypeVariableInline represents a type which stores either a value or a variable name
	GrayscaleTypeVariableInline struct {
		Name  *string
		Value *GrayscaleType
	}

	// HSL Adjusts the hue, saturation, and lightness (HSL) of an image. Hue is the number of degrees that colors rotate around the color wheel. Saturation is a multiplier to increase or decrease color saturation. Lightness is a multiplier to increase or decrease the lightness of an image. Other transformations can also affect color, such as `Grayscale` and `MaxColors`. If youre using more than one, consider the order to apply them for the desired results.
	HSL struct {
		// Hue The number of degrees to rotate colors around the color wheel, `0` by default.
		Hue *NumberVariableInline `json:"hue,omitempty"`
		// Lightness A multiplier to adjust the lightness of colors in the image. Note that lightness is distinct from brightness. For example, reducing the lightness of a light green might give you a lime green whereas reducing the brightness of a light green might give you a darker shade of the same green. Values less than `1.0` decrease the lightness of colors in the image. Values greater than `1.0` increase the lightness of colors in the image.
		Lightness *NumberVariableInline `json:"lightness,omitempty"`
		// Saturation A multiplier to adjust the saturation of colors in the image. Values less than `1.0` decrease saturation and values greater than `1.0` increase the saturation. A value of `0.0` removes all color from the image.
		Saturation *NumberVariableInline `json:"saturation,omitempty"`
		// Transformation Identifies this type of transformation, `HSL` in this case.
		Transformation HSLTransformation `json:"transformation"`
	}

	// HSLTransformation ...
	HSLTransformation string

	// HSV Identical to HSL except it replaces `lightness` with `value`. For example, if you reduce the `lightness` of a light green, almost white, image, the color turns a vibrant green. Reducing the `value` turns the image a darker color, close to grey. This happens because the original image color is very close to white.
	HSV struct {
		// Hue The number of degrees to rotate colors around the color wheel, `0.0` by default.
		Hue *NumberVariableInline `json:"hue,omitempty"`
		// Saturation A multiplier to adjust the saturation of colors in the image. Values less than `1.0` decrease saturation and values greater than `1.0` increase the saturation. A value of `0.0` removes all color from the image.
		Saturation *NumberVariableInline `json:"saturation,omitempty"`
		// Transformation Identifies this type of transformation, `HSV` in this case.
		Transformation HSVTransformation `json:"transformation"`
		// Value A multiplier to adjust the lightness or darkness of the images base color. Values less than 1.0 decrease the base colors in the image, making them appear darker. Values greater than 1.0 increase the base colors in the image, making them appear lighter.
		Value *NumberVariableInline `json:"value,omitempty"`
	}

	// HSVTransformation ...
	HSVTransformation string

	// IfDimension ...
	IfDimension struct {
		Default TransformationType `json:"default,omitempty"`
		// Dimension The dimension to use to select the transformation, either `height`, `width`, or `both`.
		Dimension   *IfDimensionDimensionVariableInline `json:"dimension,omitempty"`
		Equal       TransformationType                  `json:"equal,omitempty"`
		GreaterThan TransformationType                  `json:"greaterThan,omitempty"`
		LessThan    TransformationType                  `json:"lessThan,omitempty"`
		// Transformation Identifies this type of transformation, `IfDimension` in this case.
		Transformation IfDimensionTransformation `json:"transformation"`
		// Value The value to compare against the source image dimension. For example, if the image dimension is less than the value the lessThan transformation is applied.
		Value *IntegerVariableInline `json:"value"`
	}

	// IfDimensionDimension ...
	IfDimensionDimension string

	// IfDimensionDimensionVariableInline represents a type which stores either a value or a variable name
	IfDimensionDimensionVariableInline struct {
		Name  *string
		Value *IfDimensionDimension
	}

	// IfDimensionPost ...
	IfDimensionPost struct {
		Default TransformationTypePost `json:"default,omitempty"`
		// Dimension The dimension to use to select the transformation, either `height`, `width`, or `both`.
		Dimension   *IfDimensionPostDimensionVariableInline `json:"dimension,omitempty"`
		Equal       TransformationTypePost                  `json:"equal,omitempty"`
		GreaterThan TransformationTypePost                  `json:"greaterThan,omitempty"`
		LessThan    TransformationTypePost                  `json:"lessThan,omitempty"`
		// Transformation Identifies this type of transformation, `IfDimension` in this case.
		Transformation IfDimensionPostTransformation `json:"transformation"`
		// Value The value to compare against the source image dimension. For example, if the image dimension is less than the value the lessThan transformation is applied.
		Value *IntegerVariableInline `json:"value"`
	}

	// IfDimensionPostDimension ...
	IfDimensionPostDimension string

	// IfDimensionPostDimensionVariableInline represents a type which stores either a value or a variable name
	IfDimensionPostDimensionVariableInline struct {
		Name  *string
		Value *IfDimensionPostDimension
	}

	// IfDimensionPostTransformation ...
	IfDimensionPostTransformation string

	// IfDimensionTransformation ...
	IfDimensionTransformation string

	// IfOrientation ...
	IfOrientation struct {
		Default   TransformationType `json:"default,omitempty"`
		Landscape TransformationType `json:"landscape,omitempty"`
		Portrait  TransformationType `json:"portrait,omitempty"`
		Square    TransformationType `json:"square,omitempty"`
		// Transformation Identifies this type of transformation, `IfOrientation` in this case.
		Transformation IfOrientationTransformation `json:"transformation"`
	}

	// IfOrientationPost ...
	IfOrientationPost struct {
		Default   TransformationTypePost `json:"default,omitempty"`
		Landscape TransformationTypePost `json:"landscape,omitempty"`
		Portrait  TransformationTypePost `json:"portrait,omitempty"`
		Square    TransformationTypePost `json:"square,omitempty"`
		// Transformation Identifies this type of transformation, `IfOrientation` in this case.
		Transformation IfOrientationPostTransformation `json:"transformation"`
	}

	// IfOrientationPostTransformation ...
	IfOrientationPostTransformation string

	// IfOrientationTransformation ...
	IfOrientationTransformation string

	// ImQuery Apply artistic transformations to images quickly and dynamically by specifying transformations with a query string appendedto the image URL.
	ImQuery struct {
		// AllowedTransformations Specifies the transformations that can be applied using the query string parameter.
		AllowedTransformations []ImQueryAllowedTransformations `json:"allowedTransformations"`
		Query                  *QueryVariableInline            `json:"query"`
		// Transformation Identifies this type of transformation, `ImQuery` in this case.
		Transformation ImQueryTransformation `json:"transformation"`
	}

	// ImQueryAllowedTransformations ...
	ImQueryAllowedTransformations string

	// ImQueryTransformation ...
	ImQueryTransformation string

	// IntegerVariableInline represents a type which stores either a value or a variable name
	IntegerVariableInline struct {
		Name  *string
		Value *int
	}

	// MaxColors Set the maximum number of colors in the images palette. Reducing the number of colors in an image can help to reduce file size.
	MaxColors struct {
		// Colors The value representing the maximum number of colors to use with the source image.
		Colors *IntegerVariableInline `json:"colors"`
		// Transformation Identifies this type of transformation, `MaxColors` in this case.
		Transformation MaxColorsTransformation `json:"transformation"`
	}

	// MaxColorsTransformation ...
	MaxColorsTransformation string

	// Mirror Flips an image horizontally, vertically, or both.
	Mirror struct {
		// Horizontal Flips the image horizontally.
		Horizontal *BooleanVariableInline `json:"horizontal,omitempty"`
		// Transformation Identifies this type of transformation, `Mirror` in this case.
		Transformation MirrorTransformation `json:"transformation"`
		// Vertical Flips the image vertically.
		Vertical *BooleanVariableInline `json:"vertical,omitempty"`
	}

	// MirrorTransformation ...
	MirrorTransformation string

	// MonoHue Allows you to set all hues in an image to a single specified hue of your choosing. Mono Hue maintains the original color’s lightness and saturation but sets the hue to that of the specified value. This has the effect of making the image shades of the specified hue.
	MonoHue struct {
		// Hue Specify a hue by indicating the degree of rotation between 0 and 360 degrees around the color wheel. By default Mono Hue applies a red hue, 0.0 on the color wheel.
		Hue *NumberVariableInline `json:"hue,omitempty"`
		// Transformation Identifies this type of transformation, `MonoHue` in this case.
		Transformation MonoHueTransformation `json:"transformation"`
	}

	// MonoHueTransformation ...
	MonoHueTransformation string

	// NumberVariableInline represents a type which stores either a value or a variable name
	NumberVariableInline struct {
		Name  *string
		Value *float64
	}

	// Opacity Adjusts the level of transparency of an image. Use this transformation to make an image more or less transparent.
	Opacity struct {
		// Opacity Represents alpha values on a scale of `0` to `1`. Values below `1` increase transparency, and `0` is invisible. For images that have some transparency, values above `1` increase the opacity of the transparent portions.
		Opacity *NumberVariableInline `json:"opacity"`
		// Transformation Identifies this type of transformation, `Opacity` in this case.
		Transformation OpacityTransformation `json:"transformation"`
	}

	// OpacityTransformation ...
	OpacityTransformation string

	// OutputImage Dictates the output quality (either `quality` or `perceptualQuality`) and formats that are created for each resized image. If unspecified, image formats are created to support all browsers at the default quality level (`85`), which includes formats such as WEBP, JPEG2000 and JPEG-XR for specific browsers.
	OutputImage struct {
		// AdaptiveQuality Override the quality of image to serve when Image & Video Manager detects a slow connection. Specifying lower values lets users with slow connections browse your site with reduced load times without impacting the quality of images for users with faster connections.
		AdaptiveQuality *int `json:"adaptiveQuality,omitempty"`
		// AllowedFormats The graphics file formats allowed for browser specific results.
		AllowedFormats []OutputImageAllowedFormats `json:"allowedFormats,omitempty"`
		// ForcedFormats The forced extra formats for the `imFormat` query parameter, which requests a specific browser type. By default, Image and Video Manager detects the browser and returns the appropriate image.
		ForcedFormats []OutputImageForcedFormats `json:"forcedFormats,omitempty"`
		// PerceptualQuality Mutually exclusive with quality. The perceptual quality to use when comparing resulting images, which overrides the `quality` setting. Perceptual quality tunes each image format's quality parameter dynamically based on the human-perceived quality of the output image. This can result in better byte savings (as compared to using regular quality) as many images can be encoded at a much lower quality without compromising perception of the image. In addition, certain images may need to be encoded at a slightly higher quality in order to maintain human-perceived quality. Values are tiered `high`, `mediumHigh`, `medium`, `mediumLow`, or `low`.
		PerceptualQuality *OutputImagePerceptualQualityVariableInline `json:"perceptualQuality,omitempty"`
		// PerceptualQualityFloor Only applies with perceptualQuality set. Sets a minimum image quality to respect when using perceptual quality. Perceptual quality will not reduce the quality below this value even if it determines the compressed image to be acceptably visually similar.
		PerceptualQualityFloor *int `json:"perceptualQualityFloor,omitempty"`
		// Quality Mutually exclusive with perceptualQuality, used by default if neither is specified. The chosen quality of the output images. Using a quality value from 1-100 resembles JPEG quality across output formats.
		Quality *IntegerVariableInline `json:"quality,omitempty"`
	}

	// OutputImageAllowedFormats ...
	OutputImageAllowedFormats string

	// OutputImageForcedFormats ...
	OutputImageForcedFormats string

	// OutputImagePerceptualQuality ...
	OutputImagePerceptualQuality string

	// OutputImagePerceptualQualityVariableInline represents a type which stores either a value or a variable name
	OutputImagePerceptualQualityVariableInline struct {
		Name  *string
		Value *OutputImagePerceptualQuality
	}

	// PointShapeType Defines coordinates for a single point, to help define polygons and rectangles. Each point may be an object with `x`and `y` members, or a two-element array.
	PointShapeType struct {
		// X The horizontal position of the point, measured in pixels.
		X *NumberVariableInline `json:"x"`
		// Y The vertical position of the point, measured in pixels.
		Y *NumberVariableInline `json:"y"`
	}

	// PolicyOutputImage Specifies details for each policy, such as transformations to apply and variations in image size and formats.
	PolicyOutputImage struct {
		// Breakpoints The breakpoint widths (in pixels) to use to create derivative images/videos.
		Breakpoints *Breakpoints `json:"breakpoints,omitempty"`
		// DateCreated Date this policy version was created in ISO 8601 extended notation format.
		DateCreated string `json:"dateCreated"`
		// Hosts Hosts that are allowed for image/video URLs within transformations or variables.
		Hosts []string `json:"hosts,omitempty"`
		// ID Unique identifier for a policy, up to 64 alphanumeric characters including underscores or dashes.
		ID string `json:"id"`
		// Output Dictates the output quality (either `quality` or `perceptualQuality`) and formats that are created for each resized image. If unspecified, image formats are created to support all browsers at the default quality level (`85`), which includes formats such as WEBP, JPEG2000 and JPEG-XR for specific browsers.
		Output *OutputImage `json:"output,omitempty"`
		// PostBreakpointTransformations Post-processing Transformations are applied to the image after image and quality settings have been applied.
		PostBreakpointTransformations PostBreakpointTransformations `json:"postBreakpointTransformations,omitempty"`
		// PreviousVersion The previous version number of this policy version
		PreviousVersion int `json:"previousVersion"`
		// RolloutInfo Contains information about policy rollout start and completion times.
		RolloutInfo *RolloutInfo `json:"rolloutInfo"`
		// Transformations Set of image transformations to apply to the source image. If unspecified, no operations are performed.
		Transformations Transformations `json:"transformations,omitempty"`
		// User The user who created this policy version
		User string `json:"user"`
		// Variables Declares variables for use within the policy. Any variable declared here can be invoked throughout transformations as a [Variable](#variable) object, so that you don't have to specify values separately. You can also pass in these variable names and values dynamically as query parameters in the image's request URL.
		Variables []Variable `json:"variables,omitempty"`
		// Version The version number of this policy version
		Version int `json:"version"`
		// Video Identifies this as an image policy.
		Video *bool `json:"video,omitempty"`
	}

	// PolicyOutputImageVideo ...
	PolicyOutputImageVideo bool

	// PolygonShapeType Defines a polygon from a series of connected points.
	PolygonShapeType struct {
		// Points Series of [PointShapeType](#pointshapetype) objects. The last and first points connect to close the shape automatically.
		Points []PointShapeType `json:"points"`
	}

	// QueryVariableInline represents a type which stores either a value or a variable name
	QueryVariableInline struct {
		Name *string
	}

	// RectangleShapeType Defines a rectangle's `width` and `height` relative to an `anchor` point at the top left corner.
	RectangleShapeType struct {
		Anchor *PointShapeType `json:"anchor"`
		// Height Extends the rectangle down from the `anchor` point.
		Height *NumberVariableInline `json:"height"`
		// Width Extends the rectangle right from the `anchor` point.
		Width *NumberVariableInline `json:"width"`
	}

	// RegionOfInterestCrop Crops to a region around a specified area of interest relative to the specified `width` and `height` values.
	RegionOfInterestCrop struct {
		// Gravity The placement of the crop area relative to the specified area of interest.
		Gravity *GravityVariableInline `json:"gravity,omitempty"`
		// Height The height in pixels of the output image relative to the specified `style` value.
		Height           *IntegerVariableInline `json:"height"`
		RegionOfInterest ShapeType              `json:"regionOfInterest"`
		// Style Specifies how to crop or scale a crop area for the specified area of interest in the source image, `zoom` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` places raw crop around the point of interest, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and point of interest as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the point of interest, relative to the specified `width` and `height` values.
		Style *RegionOfInterestCropStyleVariableInline `json:"style,omitempty"`
		// Transformation Identifies this type of transformation, `RegionOfInterestCrop` in this case.
		Transformation RegionOfInterestCropTransformation `json:"transformation"`
		// Width The width in pixels of the output image relative to the specified `style` value.
		Width *IntegerVariableInline `json:"width"`
	}

	// RegionOfInterestCropStyle ...
	RegionOfInterestCropStyle string

	// RegionOfInterestCropStyleVariableInline represents a type which stores either a value or a variable name
	RegionOfInterestCropStyleVariableInline struct {
		Name  *string
		Value *RegionOfInterestCropStyle
	}

	// RegionOfInterestCropTransformation ...
	RegionOfInterestCropTransformation string

	// RelativeCrop Shrinks or expands an image relative to the image's specified dimensions. Image and Video Manager fills the expanded areas with transparency. Positive values shrink the side, while negative values expand it.
	RelativeCrop struct {
		// East The number of pixels to shrink or expand the right side of the image.
		East *IntegerVariableInline `json:"east,omitempty"`
		// North The number of pixels to shrink or expand the top side of the image.
		North *IntegerVariableInline `json:"north,omitempty"`
		// South The number of pixels to shrink or expand the bottom side of the image.
		South *IntegerVariableInline `json:"south,omitempty"`
		// Transformation Identifies this type of transformation, `RelativeCrop` in this case.
		Transformation RelativeCropTransformation `json:"transformation"`
		// West The number of pixels to shrink or expand the left side of the image.
		West *IntegerVariableInline `json:"west,omitempty"`
	}

	// RelativeCropTransformation ...
	RelativeCropTransformation string

	// RemoveColor Removes a specified color from an image and replaces it with transparent pixels. This transformation is ideal for removing solid background colors from product images photographed on clean, consistent backgrounds without any shadows.
	RemoveColor struct {
		// Color The hexadecimal CSS color value to remove.
		Color *StringVariableInline `json:"color"`
		// Feather The RemoveColor transformation may create a hard edge around an image. To minimize these hard edges and make the removal of the color more gradual in appearance, use the Feather option. This option allows you to extend the color removal beyond the specified Tolerance. The pixels in this extended tolerance become semi-transparent - creating a softer edge.  The first realtime request for an image using the feather option may result in a slow transformation time. Subsequent requests are not impacted as they are served directly out of cache.
		Feather *NumberVariableInline `json:"feather,omitempty"`
		// Tolerance The Tolerance option defines how close the color needs to be to the selected color before it's changed to fully transparent. Set the Tolerance to 0.0 to remove only the exact color specified.
		Tolerance *NumberVariableInline `json:"tolerance,omitempty"`
		// Transformation Identifies this type of transformation, `RemoveColor` in this case.
		Transformation RemoveColorTransformation `json:"transformation"`
	}

	// RemoveColorTransformation ...
	RemoveColorTransformation string

	// Resize Resizes an image to a particular, absolute dimension. If you don't enter a `width` or a `height`, the image is resized with the `fit` aspect preservation mode, which selects a value for the missing dimension that preserves the image's aspect.
	Resize struct {
		// Aspect Preserves the aspect ratio. Select `fit` to make the image fit entirely within the selected width and height. When using `fit`, the resulting image has the largest possible size for the specified dimensions. Select `fill` to size the image so it both completely fills the dimensions and has the smallest possible file size. Otherwise `ignore` changes the original aspect ratio to fit within an arbitrarily shaped rectangle.
		Aspect *ResizeAspectVariableInline `json:"aspect,omitempty"`
		// Height The height to resize the source image to. Must be set if height is not specified.
		Height *IntegerVariableInline `json:"height,omitempty"`
		// Transformation Identifies this type of transformation, `Resize` in this case.
		Transformation ResizeTransformation `json:"transformation"`
		// Type Sets constraints for the image resize. Select `normal` to resize in all cases, either increasing or decreasing the dimensions. Select `downsize` to ignore this transformation if the result would be larger than the original. Select `upsize` to ignore this transformation if the result would be smaller.
		Type *ResizeTypeVariableInline `json:"type,omitempty"`
		// Width The width to resize the source image to. Must be set if width is not specified.
		Width *IntegerVariableInline `json:"width,omitempty"`
	}

	// ResizeAspect ...
	ResizeAspect string

	// ResizeAspectVariableInline represents a type which stores either a value or a variable name
	ResizeAspectVariableInline struct {
		Name  *string
		Value *ResizeAspect
	}

	// ResizeTransformation ...
	ResizeTransformation string

	// ResizeType ...
	ResizeType string

	// ResizeTypeVariableInline represents a type which stores either a value or a variable name
	ResizeTypeVariableInline struct {
		Name  *string
		Value *ResizeType
	}

	// RolloutInfo Contains information about policy rollout start and completion times.
	RolloutInfo struct {
		// EndTime The estimated time that rollout for this policy will end. Value is a unix timestamp.
		EndTime int `json:"endTime"`
		// RolloutDuration The amount of time in seconds that the policy takes to rollout. During the rollout an increasing proportion of images/videos will begin to use the new policy instead of the cached images/videos from the previous version. Policies on the staging network deploy as quickly as possible without rollout. For staging policies this value will always be 1.
		RolloutDuration int `json:"rolloutDuration"`
		// StartTime The estimated time that rollout for this policy will begin. Value is a unix timestamp.
		StartTime int `json:"startTime"`
	}

	// Rotate Rotate the image around its center by indicating the degrees of rotation.
	Rotate struct {
		// Degrees The value to rotate the image by. Positive values rotate clockwise, while negative values rotate counter-clockwise.
		Degrees *NumberVariableInline `json:"degrees"`
		// Transformation Identifies this type of transformation, `Rotate` in this case.
		Transformation RotateTransformation `json:"transformation"`
	}

	// RotateTransformation ...
	RotateTransformation string

	// Scale Changes the image's size to different dimensions relative to its starting size.
	Scale struct {
		// Height Scaling factor for the input height to determine the output height of the image, where values between `0` and `1` decrease size. Image dimensions need to be non-zero positive numbers.
		Height *NumberVariableInline `json:"height"`
		// Transformation Identifies this type of transformation, `Scale` in this case.
		Transformation ScaleTransformation `json:"transformation"`
		// Width Scaling factor for the input width to determine the output width of the image, where `1` leaves the width unchanged. Values greater than `1` increase the image size. Image dimensions need to be non-zero positive numbers.
		Width *NumberVariableInline `json:"width"`
	}

	// ScaleTransformation ...
	ScaleTransformation string

	// Shear Slants an image into a parallelogram, as a percent of the starting dimension as represented in decimal format. You need to specify at least one axis property. Transparent pixels fill empty areas around the sheared image as needed, so it's often useful to use a `BackgroundColor` transformation for these areas.
	Shear struct {
		// Transformation Identifies this type of transformation, `Shear` in this case.
		Transformation ShearTransformation `json:"transformation"`
		// XShear The amount to shear along the x-axis, measured in multiples of the image's width. Must be set if yShear is not specified.
		XShear *NumberVariableInline `json:"xShear,omitempty"`
		// YShear The amount to shear along the y-axis, measured in multiples of the image's height. Must be set if xShear is not specified.
		YShear *NumberVariableInline `json:"yShear,omitempty"`
	}

	// ShearTransformation ...
	ShearTransformation string

	// StringVariableInline represents a type which stores either a value or a variable name
	StringVariableInline struct {
		Name  *string
		Value *string
	}

	// TextImageType A snippet of text. Defines font family and size, fill color, and outline stroke width and color.
	TextImageType struct {
		// Fill The main fill color of the text.
		Fill *StringVariableInline `json:"fill,omitempty"`
		// Size The size in pixels to render the text.
		Size *NumberVariableInline `json:"size,omitempty"`
		// Stroke The color for the outline of the text.
		Stroke *StringVariableInline `json:"stroke,omitempty"`
		// StrokeSize The thickness in points for the outline of the text.
		StrokeSize *NumberVariableInline `json:"strokeSize,omitempty"`
		// Text The line of text to render.
		Text           *StringVariableInline `json:"text"`
		Transformation TransformationType    `json:"transformation,omitempty"`
		// Type Identifies the type of image, `Text` in this case.
		Type TextImageTypeType `json:"type"`
		// Typeface The font family to apply to the text image. This may be a URL to a TrueType or WOFF (v1) typeface, or a string that refers to one of the standard built-in browser fonts.
		Typeface *StringVariableInline `json:"typeface,omitempty"`
	}

	// TextImageTypePost A snippet of text. Defines font family and size, fill color, and outline stroke width and color.
	TextImageTypePost struct {
		// Fill The main fill color of the text.
		Fill *StringVariableInline `json:"fill,omitempty"`
		// Size The size in pixels to render the text.
		Size *NumberVariableInline `json:"size,omitempty"`
		// Stroke The color for the outline of the text.
		Stroke *StringVariableInline `json:"stroke,omitempty"`
		// StrokeSize The thickness in points for the outline of the text.
		StrokeSize *NumberVariableInline `json:"strokeSize,omitempty"`
		// Text The line of text to render.
		Text           *StringVariableInline  `json:"text"`
		Transformation TransformationTypePost `json:"transformation,omitempty"`
		// Type Identifies the type of image, `Text` in this case.
		Type TextImageTypePostType `json:"type"`
		// Typeface The font family to apply to the text image. This may be a URL to a TrueType or WOFF (v1) typeface, or a string that refers to one of the standard built-in browser fonts.
		Typeface *StringVariableInline `json:"typeface,omitempty"`
	}

	// TextImageTypePostType ...
	TextImageTypePostType string

	// TextImageTypeType ...
	TextImageTypeType string

	// Trim Automatically crops uniform backgrounds from the edges of an image.
	Trim struct {
		// Fuzz The fuzz tolerance of the trim, a value between `0` and `1` that determines the acceptable amount of background variation before trimming stops.
		Fuzz *NumberVariableInline `json:"fuzz,omitempty"`
		// Padding The amount of padding in pixels to add to the trimmed image.
		Padding *IntegerVariableInline `json:"padding,omitempty"`
		// Transformation Identifies this type of transformation, `Trim` in this case.
		Transformation TrimTransformation `json:"transformation"`
	}

	// TrimTransformation ...
	TrimTransformation string

	// URLImageType An image loaded from a URL.
	URLImageType struct {
		Transformation TransformationType `json:"transformation,omitempty"`
		// Type Identifies the type of image, `URL` in this case.
		Type URLImageTypeType `json:"type,omitempty"`
		// URL The URL of the image.
		URL *StringVariableInline `json:"url"`
	}

	// URLImageTypePost An image loaded from a URL.
	URLImageTypePost struct {
		Transformation TransformationTypePost `json:"transformation,omitempty"`
		// Type Identifies the type of image, `URL` in this case.
		Type URLImageTypePostType `json:"type,omitempty"`
		// URL The URL of the image.
		URL *StringVariableInline `json:"url"`
	}

	// URLImageTypePostType ...
	URLImageTypePostType string

	// URLImageTypeType ...
	URLImageTypeType string

	// UnionShapeType Identifies a combined shape based on a set of other shapes. You can use a full JSON object to represent a union or an array of shapes that describe it.
	UnionShapeType struct {
		Shapes []ShapeType `json:"shapes"`
	}

	// UnsharpMask Emphasizes edges and details in source images without distorting the colors. Although this effect is often referred to as _sharpening_ an image, it actually creates a blurred, inverted copy of the image known as an unsharp mask. Image and Video Manager combines the unsharp mask with the source image to create an image perceived as clearer.
	UnsharpMask struct {
		// Gain Set how much emphasis the filter applies to details. Higher values increase apparent sharpness of details.
		Gain *NumberVariableInline `json:"gain,omitempty"`
		// Sigma The standard deviation of the Gaussian distribution used in the in unsharp mask, measured in pixels, `1.0` by default. High values emphasize large details and low values emphasize small details.
		Sigma *NumberVariableInline `json:"sigma,omitempty"`
		// Threshold Set the minimum change required to include a detail in the filter. Higher values discard more changes.
		Threshold *NumberVariableInline `json:"threshold,omitempty"`
		// Transformation Identifies this type of transformation, `UnsharpMask` in this case.
		Transformation UnsharpMaskTransformation `json:"transformation"`
	}

	// UnsharpMaskTransformation ...
	UnsharpMaskTransformation string

	// Variable ...
	Variable struct {
		// DefaultValue The default value of the variable if no query parameter is provided. It needs to be one of the `enumOptions` if any are provided.
		DefaultValue string         `json:"defaultValue"`
		EnumOptions  []*EnumOptions `json:"enumOptions,omitempty"`
		// Name The name of the variable, also available as the query parameter name to set the variable's value dynamically. Use up to 50 alphanumeric characters.
		Name string `json:"name"`
		// Postfix A postfix added to the value provided for the variable, or to the default value.
		Postfix *string `json:"postfix,omitempty"`
		// Prefix A prefix added to the value provided for the variable, or to the default value.
		Prefix *string `json:"prefix,omitempty"`
		// Type The type of value for the variable.
		Type VariableType `json:"type"`
	}

	// VariableInline References the name of a variable defined [by the policy](#63c7bea4). Use this object to substitute preset values within transformations, or to pass in values dynamically using image URL query parameters.
	VariableInline struct {
		// Var Corresponds to the `name` of the variable declared by the policy, to insert the corresponding value.
		Var string `json:"var"`
	}

	// VariableType ...
	VariableType string

	// OutputVideo Dictates the output quality that are created for each resized video.
	OutputVideo struct {
		// PerceptualQuality The quality of derivative videos. High preserves video quality with reduced byte savings while low reduces video quality to increase byte savings.
		PerceptualQuality *OutputVideoPerceptualQualityVariableInline `json:"perceptualQuality,omitempty"`
		// PlaceholderVideoURL Allows you to add a specific placeholder video that appears when a user first requests a video, but before Image & Video Manager processes the video. If not specified the original video plays during the processing time.
		PlaceholderVideoURL *StringVariableInline `json:"placeholderVideoUrl,omitempty"`
		// VideoAdaptiveQuality Override the quality of video to serve when Image & Video Manager detects a slow connection. Specifying lower values lets users with slow connections browse your site with reduced load times without impacting the quality of videos for users with faster connections.
		VideoAdaptiveQuality *OutputVideoVideoAdaptiveQualityVariableInline `json:"videoAdaptiveQuality,omitempty"`
	}

	// OutputVideoPerceptualQuality ...
	OutputVideoPerceptualQuality string

	// OutputVideoPerceptualQualityVariableInline represents a type which stores either a value or a variable name
	OutputVideoPerceptualQualityVariableInline struct {
		Name  *string
		Value *OutputVideoPerceptualQuality
	}

	// OutputVideoVideoAdaptiveQuality ...
	OutputVideoVideoAdaptiveQuality string

	// OutputVideoVideoAdaptiveQualityVariableInline represents a type which stores either a value or a variable name
	OutputVideoVideoAdaptiveQualityVariableInline struct {
		Name  *string
		Value *OutputVideoVideoAdaptiveQuality
	}

	// PolicyOutputVideo Specifies details for each policy such as video size.
	PolicyOutputVideo struct {
		// Breakpoints The breakpoint widths (in pixels) to use to create derivative images/videos.
		Breakpoints *Breakpoints `json:"breakpoints,omitempty"`
		// DateCreated Date this policy version was created in ISO 8601 extended notation format.
		DateCreated string `json:"dateCreated"`
		// Hosts Hosts that are allowed for image/video URLs within transformations or variables.
		Hosts []string `json:"hosts,omitempty"`
		// ID Unique identifier for a policy, up to 64 alphanumeric characters including underscores or dashes.
		ID string `json:"id"`
		// Output Dictates the output quality that are created for each resized video.
		Output *OutputVideo `json:"output,omitempty"`
		// PreviousVersion The previous version number of this policy version
		PreviousVersion int `json:"previousVersion"`
		// RolloutInfo Contains information about policy rollout start and completion times.
		RolloutInfo *RolloutInfo `json:"rolloutInfo"`
		// User The user who created this policy version
		User string `json:"user"`
		// Variables Declares variables for use within the policy. Any variable declared here can be invoked throughout transformations as a [Variable](#variable) object, so that you don't have to specify values separately. You can also pass in these variable names and values dynamically as query parameters in the image's request URL.
		Variables []Variable `json:"variables,omitempty"`
		// Version The version number of this policy version
		Version int `json:"version"`
		// Video Identifies this as a video policy.
		Video *bool `json:"video,omitempty"`
	}

	// PolicyOutputVideoVideo ...
	PolicyOutputVideoVideo bool
)

/*-----------------------------------------------*/
/////////////// Generated constants ///////////////
/*-----------------------------------------------*/
const (

	// AppendGravityPriorityHorizontal const
	AppendGravityPriorityHorizontal AppendGravityPriority = "horizontal"
	// AppendGravityPriorityVertical const
	AppendGravityPriorityVertical AppendGravityPriority = "vertical"

	// AppendTransformationAppend const
	AppendTransformationAppend AppendTransformation = "Append"

	// AspectCropTransformationAspectCrop const
	AspectCropTransformationAspectCrop AspectCropTransformation = "AspectCrop"

	// BackgroundColorTransformationBackgroundColor const
	BackgroundColorTransformationBackgroundColor BackgroundColorTransformation = "BackgroundColor"

	// BlurTransformationBlur const
	BlurTransformationBlur BlurTransformation = "Blur"

	// BoxImageTypePostTypeBox const
	BoxImageTypePostTypeBox BoxImageTypePostType = "Box"

	// BoxImageTypeTypeBox const
	BoxImageTypeTypeBox BoxImageTypeType = "Box"

	// ChromaKeyTransformationChromaKey const
	ChromaKeyTransformationChromaKey ChromaKeyTransformation = "ChromaKey"

	// CircleImageTypePostTypeCircle const
	CircleImageTypePostTypeCircle CircleImageTypePostType = "Circle"

	// CircleImageTypeTypeCircle const
	CircleImageTypeTypeCircle CircleImageTypeType = "Circle"

	// CompositePlacementOver const
	CompositePlacementOver CompositePlacement = "Over"
	// CompositePlacementUnder const
	CompositePlacementUnder CompositePlacement = "Under"
	// CompositePlacementMask const
	CompositePlacementMask CompositePlacement = "Mask"
	// CompositePlacementStencil const
	CompositePlacementStencil CompositePlacement = "Stencil"

	// CompositePostPlacementOver const
	CompositePostPlacementOver CompositePostPlacement = "Over"
	// CompositePostPlacementUnder const
	CompositePostPlacementUnder CompositePostPlacement = "Under"
	// CompositePostPlacementMask const
	CompositePostPlacementMask CompositePostPlacement = "Mask"
	// CompositePostPlacementStencil const
	CompositePostPlacementStencil CompositePostPlacement = "Stencil"

	// CompositePostScaleDimensionWidth const
	CompositePostScaleDimensionWidth CompositePostScaleDimension = "width"
	// CompositePostScaleDimensionHeight const
	CompositePostScaleDimensionHeight CompositePostScaleDimension = "height"

	// CompositePostTransformationComposite const
	CompositePostTransformationComposite CompositePostTransformation = "Composite"

	// CompositeScaleDimensionWidth const
	CompositeScaleDimensionWidth CompositeScaleDimension = "width"
	// CompositeScaleDimensionHeight const
	CompositeScaleDimensionHeight CompositeScaleDimension = "height"

	// CompositeTransformationComposite const
	CompositeTransformationComposite CompositeTransformation = "Composite"

	// CompoundPostTransformationCompound const
	CompoundPostTransformationCompound CompoundPostTransformation = "Compound"

	// CompoundTransformationCompound const
	CompoundTransformationCompound CompoundTransformation = "Compound"

	// ContrastTransformationContrast const
	ContrastTransformationContrast ContrastTransformation = "Contrast"

	// CropTransformationCrop const
	CropTransformationCrop CropTransformation = "Crop"

	// FaceCropAlgorithmCascade const
	FaceCropAlgorithmCascade FaceCropAlgorithm = "cascade"
	// FaceCropAlgorithmDnn const
	FaceCropAlgorithmDnn FaceCropAlgorithm = "dnn"

	// FaceCropFocusAllFaces const
	FaceCropFocusAllFaces FaceCropFocus = "allFaces"
	// FaceCropFocusBiggestFace const
	FaceCropFocusBiggestFace FaceCropFocus = "biggestFace"

	// FaceCropStyleCrop const
	FaceCropStyleCrop FaceCropStyle = "crop"
	// FaceCropStyleFill const
	FaceCropStyleFill FaceCropStyle = "fill"
	// FaceCropStyleZoom const
	FaceCropStyleZoom FaceCropStyle = "zoom"

	// FaceCropTransformationFaceCrop const
	FaceCropTransformationFaceCrop FaceCropTransformation = "FaceCrop"

	// FeatureCropStyleCrop const
	FeatureCropStyleCrop FeatureCropStyle = "crop"
	// FeatureCropStyleFill const
	FeatureCropStyleFill FeatureCropStyle = "fill"
	// FeatureCropStyleZoom const
	FeatureCropStyleZoom FeatureCropStyle = "zoom"

	// FeatureCropTransformationFeatureCrop const
	FeatureCropTransformationFeatureCrop FeatureCropTransformation = "FeatureCrop"

	// FitAndFillTransformationFitAndFill const
	FitAndFillTransformationFitAndFill FitAndFillTransformation = "FitAndFill"

	// GoopTransformationGoop const
	GoopTransformationGoop GoopTransformation = "Goop"

	// GravityNorth const
	GravityNorth Gravity = "North"
	// GravityNorthEast const
	GravityNorthEast Gravity = "NorthEast"
	// GravityNorthWest const
	GravityNorthWest Gravity = "NorthWest"
	// GravitySouth const
	GravitySouth Gravity = "South"
	// GravitySouthEast const
	GravitySouthEast Gravity = "SouthEast"
	// GravitySouthWest const
	GravitySouthWest Gravity = "SouthWest"
	// GravityCenter const
	GravityCenter Gravity = "Center"
	// GravityEast const
	GravityEast Gravity = "East"
	// GravityWest const
	GravityWest Gravity = "West"

	// GravityPostNorth const
	GravityPostNorth GravityPost = "North"
	// GravityPostNorthEast const
	GravityPostNorthEast GravityPost = "NorthEast"
	// GravityPostNorthWest const
	GravityPostNorthWest GravityPost = "NorthWest"
	// GravityPostSouth const
	GravityPostSouth GravityPost = "South"
	// GravityPostSouthEast const
	GravityPostSouthEast GravityPost = "SouthEast"
	// GravityPostSouthWest const
	GravityPostSouthWest GravityPost = "SouthWest"
	// GravityPostCenter const
	GravityPostCenter GravityPost = "Center"
	// GravityPostEast const
	GravityPostEast GravityPost = "East"
	// GravityPostWest const
	GravityPostWest GravityPost = "West"

	// GrayscaleTransformationGrayscale const
	GrayscaleTransformationGrayscale GrayscaleTransformation = "Grayscale"

	// GrayscaleTypeRec601 const
	GrayscaleTypeRec601 GrayscaleType = "Rec601"
	// GrayscaleTypeRec709 const
	GrayscaleTypeRec709 GrayscaleType = "Rec709"
	// GrayscaleTypeBrightness const
	GrayscaleTypeBrightness GrayscaleType = "Brightness"
	// GrayscaleTypeLightness const
	GrayscaleTypeLightness GrayscaleType = "Lightness"

	// HSLTransformationHSL const
	HSLTransformationHSL HSLTransformation = "HSL"

	// HSVTransformationHSV const
	HSVTransformationHSV HSVTransformation = "HSV"

	// IfDimensionDimensionWidth const
	IfDimensionDimensionWidth IfDimensionDimension = "width"
	// IfDimensionDimensionHeight const
	IfDimensionDimensionHeight IfDimensionDimension = "height"
	// IfDimensionDimensionBoth const
	IfDimensionDimensionBoth IfDimensionDimension = "both"

	// IfDimensionPostDimensionWidth const
	IfDimensionPostDimensionWidth IfDimensionPostDimension = "width"
	// IfDimensionPostDimensionHeight const
	IfDimensionPostDimensionHeight IfDimensionPostDimension = "height"
	// IfDimensionPostDimensionBoth const
	IfDimensionPostDimensionBoth IfDimensionPostDimension = "both"

	// IfDimensionPostTransformationIfDimension const
	IfDimensionPostTransformationIfDimension IfDimensionPostTransformation = "IfDimension"

	// IfDimensionTransformationIfDimension const
	IfDimensionTransformationIfDimension IfDimensionTransformation = "IfDimension"

	// IfOrientationPostTransformationIfOrientation const
	IfOrientationPostTransformationIfOrientation IfOrientationPostTransformation = "IfOrientation"

	// IfOrientationTransformationIfOrientation const
	IfOrientationTransformationIfOrientation IfOrientationTransformation = "IfOrientation"

	// ImQueryAllowedTransformationsAppend const
	ImQueryAllowedTransformationsAppend ImQueryAllowedTransformations = "Append"
	// ImQueryAllowedTransformationsAspectCrop const
	ImQueryAllowedTransformationsAspectCrop ImQueryAllowedTransformations = "AspectCrop"
	// ImQueryAllowedTransformationsBackgroundColor const
	ImQueryAllowedTransformationsBackgroundColor ImQueryAllowedTransformations = "BackgroundColor"
	// ImQueryAllowedTransformationsBlur const
	ImQueryAllowedTransformationsBlur ImQueryAllowedTransformations = "Blur"
	// ImQueryAllowedTransformationsComposite const
	ImQueryAllowedTransformationsComposite ImQueryAllowedTransformations = "Composite"
	// ImQueryAllowedTransformationsContrast const
	ImQueryAllowedTransformationsContrast ImQueryAllowedTransformations = "Contrast"
	// ImQueryAllowedTransformationsCrop const
	ImQueryAllowedTransformationsCrop ImQueryAllowedTransformations = "Crop"
	// ImQueryAllowedTransformationsChromaKey const
	ImQueryAllowedTransformationsChromaKey ImQueryAllowedTransformations = "ChromaKey"
	// ImQueryAllowedTransformationsFaceCrop const
	ImQueryAllowedTransformationsFaceCrop ImQueryAllowedTransformations = "FaceCrop"
	// ImQueryAllowedTransformationsFeatureCrop const
	ImQueryAllowedTransformationsFeatureCrop ImQueryAllowedTransformations = "FeatureCrop"
	// ImQueryAllowedTransformationsFitAndFill const
	ImQueryAllowedTransformationsFitAndFill ImQueryAllowedTransformations = "FitAndFill"
	// ImQueryAllowedTransformationsGoop const
	ImQueryAllowedTransformationsGoop ImQueryAllowedTransformations = "Goop"
	// ImQueryAllowedTransformationsGrayscale const
	ImQueryAllowedTransformationsGrayscale ImQueryAllowedTransformations = "Grayscale"
	// ImQueryAllowedTransformationsHSL const
	ImQueryAllowedTransformationsHSL ImQueryAllowedTransformations = "HSL"
	// ImQueryAllowedTransformationsHSV const
	ImQueryAllowedTransformationsHSV ImQueryAllowedTransformations = "HSV"
	// ImQueryAllowedTransformationsMaxColors const
	ImQueryAllowedTransformationsMaxColors ImQueryAllowedTransformations = "MaxColors"
	// ImQueryAllowedTransformationsMirror const
	ImQueryAllowedTransformationsMirror ImQueryAllowedTransformations = "Mirror"
	// ImQueryAllowedTransformationsMonoHue const
	ImQueryAllowedTransformationsMonoHue ImQueryAllowedTransformations = "MonoHue"
	// ImQueryAllowedTransformationsOpacity const
	ImQueryAllowedTransformationsOpacity ImQueryAllowedTransformations = "Opacity"
	// ImQueryAllowedTransformationsRegionOfInterestCrop const
	ImQueryAllowedTransformationsRegionOfInterestCrop ImQueryAllowedTransformations = "RegionOfInterestCrop"
	// ImQueryAllowedTransformationsRelativeCrop const
	ImQueryAllowedTransformationsRelativeCrop ImQueryAllowedTransformations = "RelativeCrop"
	// ImQueryAllowedTransformationsRemoveColor const
	ImQueryAllowedTransformationsRemoveColor ImQueryAllowedTransformations = "RemoveColor"
	// ImQueryAllowedTransformationsResize const
	ImQueryAllowedTransformationsResize ImQueryAllowedTransformations = "Resize"
	// ImQueryAllowedTransformationsRotate const
	ImQueryAllowedTransformationsRotate ImQueryAllowedTransformations = "Rotate"
	// ImQueryAllowedTransformationsScale const
	ImQueryAllowedTransformationsScale ImQueryAllowedTransformations = "Scale"
	// ImQueryAllowedTransformationsShear const
	ImQueryAllowedTransformationsShear ImQueryAllowedTransformations = "Shear"
	// ImQueryAllowedTransformationsTrim const
	ImQueryAllowedTransformationsTrim ImQueryAllowedTransformations = "Trim"
	// ImQueryAllowedTransformationsUnsharpMask const
	ImQueryAllowedTransformationsUnsharpMask ImQueryAllowedTransformations = "UnsharpMask"
	// ImQueryAllowedTransformationsIfDimension const
	ImQueryAllowedTransformationsIfDimension ImQueryAllowedTransformations = "IfDimension"
	// ImQueryAllowedTransformationsIfOrientation const
	ImQueryAllowedTransformationsIfOrientation ImQueryAllowedTransformations = "IfOrientation"

	// ImQueryTransformationImQuery const
	ImQueryTransformationImQuery ImQueryTransformation = "ImQuery"

	// MaxColorsTransformationMaxColors const
	MaxColorsTransformationMaxColors MaxColorsTransformation = "MaxColors"

	// MirrorTransformationMirror const
	MirrorTransformationMirror MirrorTransformation = "Mirror"

	// MonoHueTransformationMonoHue const
	MonoHueTransformationMonoHue MonoHueTransformation = "MonoHue"

	// OpacityTransformationOpacity const
	OpacityTransformationOpacity OpacityTransformation = "Opacity"

	// OutputImageAllowedFormatsGif const
	OutputImageAllowedFormatsGif OutputImageAllowedFormats = "gif"
	// OutputImageAllowedFormatsJpeg const
	OutputImageAllowedFormatsJpeg OutputImageAllowedFormats = "jpeg"
	// OutputImageAllowedFormatsPng const
	OutputImageAllowedFormatsPng OutputImageAllowedFormats = "png"
	// OutputImageAllowedFormatsWebp const
	OutputImageAllowedFormatsWebp OutputImageAllowedFormats = "webp"
	// OutputImageAllowedFormatsJpegxr const
	OutputImageAllowedFormatsJpegxr OutputImageAllowedFormats = "jpegxr"
	// OutputImageAllowedFormatsJpeg2000 const
	OutputImageAllowedFormatsJpeg2000 OutputImageAllowedFormats = "jpeg2000"

	// OutputImageForcedFormatsGif const
	OutputImageForcedFormatsGif OutputImageForcedFormats = "gif"
	// OutputImageForcedFormatsJpeg const
	OutputImageForcedFormatsJpeg OutputImageForcedFormats = "jpeg"
	// OutputImageForcedFormatsPng const
	OutputImageForcedFormatsPng OutputImageForcedFormats = "png"
	// OutputImageForcedFormatsWebp const
	OutputImageForcedFormatsWebp OutputImageForcedFormats = "webp"
	// OutputImageForcedFormatsJpegxr const
	OutputImageForcedFormatsJpegxr OutputImageForcedFormats = "jpegxr"
	// OutputImageForcedFormatsJpeg2000 const
	OutputImageForcedFormatsJpeg2000 OutputImageForcedFormats = "jpeg2000"

	// OutputImagePerceptualQualityHigh const
	OutputImagePerceptualQualityHigh OutputImagePerceptualQuality = "high"
	// OutputImagePerceptualQualityMediumHigh const
	OutputImagePerceptualQualityMediumHigh OutputImagePerceptualQuality = "mediumHigh"
	// OutputImagePerceptualQualityMedium const
	OutputImagePerceptualQualityMedium OutputImagePerceptualQuality = "medium"
	// OutputImagePerceptualQualityMediumLow const
	OutputImagePerceptualQualityMediumLow OutputImagePerceptualQuality = "mediumLow"
	// OutputImagePerceptualQualityLow const
	OutputImagePerceptualQualityLow OutputImagePerceptualQuality = "low"

	// PolicyOutputImageVideoFalse const
	PolicyOutputImageVideoFalse PolicyOutputImageVideo = false

	// RegionOfInterestCropStyleCrop const
	RegionOfInterestCropStyleCrop RegionOfInterestCropStyle = "crop"
	// RegionOfInterestCropStyleFill const
	RegionOfInterestCropStyleFill RegionOfInterestCropStyle = "fill"
	// RegionOfInterestCropStyleZoom const
	RegionOfInterestCropStyleZoom RegionOfInterestCropStyle = "zoom"

	// RegionOfInterestCropTransformationRegionOfInterestCrop const
	RegionOfInterestCropTransformationRegionOfInterestCrop RegionOfInterestCropTransformation = "RegionOfInterestCrop"

	// RelativeCropTransformationRelativeCrop const
	RelativeCropTransformationRelativeCrop RelativeCropTransformation = "RelativeCrop"

	// RemoveColorTransformationRemoveColor const
	RemoveColorTransformationRemoveColor RemoveColorTransformation = "RemoveColor"

	// ResizeAspectFit const
	ResizeAspectFit ResizeAspect = "fit"
	// ResizeAspectFill const
	ResizeAspectFill ResizeAspect = "fill"
	// ResizeAspectIgnore const
	ResizeAspectIgnore ResizeAspect = "ignore"

	// ResizeTransformationResize const
	ResizeTransformationResize ResizeTransformation = "Resize"

	// ResizeTypeNormal const
	ResizeTypeNormal ResizeType = "normal"
	// ResizeTypeUpsize const
	ResizeTypeUpsize ResizeType = "upsize"
	// ResizeTypeDownsize const
	ResizeTypeDownsize ResizeType = "downsize"

	// RotateTransformationRotate const
	RotateTransformationRotate RotateTransformation = "Rotate"

	// ScaleTransformationScale const
	ScaleTransformationScale ScaleTransformation = "Scale"

	// ShearTransformationShear const
	ShearTransformationShear ShearTransformation = "Shear"

	// TextImageTypePostTypeText const
	TextImageTypePostTypeText TextImageTypePostType = "Text"

	// TextImageTypeTypeText const
	TextImageTypeTypeText TextImageTypeType = "Text"

	// TrimTransformationTrim const
	TrimTransformationTrim TrimTransformation = "Trim"

	// URLImageTypePostTypeURL const
	URLImageTypePostTypeURL URLImageTypePostType = "URL"

	// URLImageTypeTypeURL const
	URLImageTypeTypeURL URLImageTypeType = "URL"

	// UnsharpMaskTransformationUnsharpMask const
	UnsharpMaskTransformationUnsharpMask UnsharpMaskTransformation = "UnsharpMask"

	// VariableTypeBool const
	VariableTypeBool VariableType = "bool"
	// VariableTypeNumber const
	VariableTypeNumber VariableType = "number"
	// VariableTypeURL const
	VariableTypeURL VariableType = "url"
	// VariableTypeColor const
	VariableTypeColor VariableType = "color"
	// VariableTypeGravity const
	VariableTypeGravity VariableType = "gravity"
	// VariableTypePlacement const
	VariableTypePlacement VariableType = "placement"
	// VariableTypeScaleDimension const
	VariableTypeScaleDimension VariableType = "scaleDimension"
	// VariableTypeGrayscaleType const
	VariableTypeGrayscaleType VariableType = "grayscaleType"
	// VariableTypeAspect const
	VariableTypeAspect VariableType = "aspect"
	// VariableTypeResizeType const
	VariableTypeResizeType VariableType = "resizeType"
	// VariableTypeDimension const
	VariableTypeDimension VariableType = "dimension"
	// VariableTypePerceptualQuality const
	VariableTypePerceptualQuality VariableType = "perceptualQuality"
	// VariableTypeString const
	VariableTypeString VariableType = "string"
	// VariableTypeFocus const
	VariableTypeFocus VariableType = "focus"

	// OutputVideoPerceptualQualityHigh const
	OutputVideoPerceptualQualityHigh OutputVideoPerceptualQuality = "high"
	// OutputVideoPerceptualQualityMediumHigh const
	OutputVideoPerceptualQualityMediumHigh OutputVideoPerceptualQuality = "mediumHigh"
	// OutputVideoPerceptualQualityMedium const
	OutputVideoPerceptualQualityMedium OutputVideoPerceptualQuality = "medium"
	// OutputVideoPerceptualQualityMediumLow const
	OutputVideoPerceptualQualityMediumLow OutputVideoPerceptualQuality = "mediumLow"
	// OutputVideoPerceptualQualityLow const
	OutputVideoPerceptualQualityLow OutputVideoPerceptualQuality = "low"

	// OutputVideoVideoAdaptiveQualityHigh const
	OutputVideoVideoAdaptiveQualityHigh OutputVideoVideoAdaptiveQuality = "high"
	// OutputVideoVideoAdaptiveQualityMediumHigh const
	OutputVideoVideoAdaptiveQualityMediumHigh OutputVideoVideoAdaptiveQuality = "mediumHigh"
	// OutputVideoVideoAdaptiveQualityMedium const
	OutputVideoVideoAdaptiveQualityMedium OutputVideoVideoAdaptiveQuality = "medium"
	// OutputVideoVideoAdaptiveQualityMediumLow const
	OutputVideoVideoAdaptiveQualityMediumLow OutputVideoVideoAdaptiveQuality = "mediumLow"
	// OutputVideoVideoAdaptiveQualityLow const
	OutputVideoVideoAdaptiveQualityLow OutputVideoVideoAdaptiveQuality = "low"

	// PolicyOutputVideoVideoTrue const
	PolicyOutputVideoVideoTrue PolicyOutputVideoVideo = true
)

/*-----------------------------------------------*/
//////////// Interface implementations ////////////
/*-----------------------------------------------*/

func (Append) transformationType() string {
	return "Append"
}

func (AspectCrop) transformationType() string {
	return "AspectCrop"
}

func (BackgroundColor) transformationType() string {
	return "BackgroundColor"
}

func (Blur) transformationType() string {
	return "Blur"
}

func (BoxImageType) imageType() string {
	return "BoxImageType"
}

func (BoxImageTypePost) imageTypePost() string {
	return "BoxImageTypePost"
}

func (ChromaKey) transformationType() string {
	return "ChromaKey"
}

func (CircleImageType) imageType() string {
	return "CircleImageType"
}

func (CircleImageTypePost) imageTypePost() string {
	return "CircleImageTypePost"
}

func (CircleShapeType) shapeType() string {
	return "CircleShapeType"
}

func (Composite) transformationType() string {
	return "Composite"
}

func (Compound) transformationType() string {
	return "Compound"
}

func (Contrast) transformationType() string {
	return "Contrast"
}

func (Crop) transformationType() string {
	return "Crop"
}

func (FaceCrop) transformationType() string {
	return "FaceCrop"
}

func (FeatureCrop) transformationType() string {
	return "FeatureCrop"
}

func (FitAndFill) transformationType() string {
	return "FitAndFill"
}

func (Goop) transformationType() string {
	return "Goop"
}

func (Grayscale) transformationType() string {
	return "Grayscale"
}

func (HSL) transformationType() string {
	return "HSL"
}

func (HSV) transformationType() string {
	return "HSV"
}

func (IfDimension) transformationType() string {
	return "IfDimension"
}

func (IfOrientation) transformationType() string {
	return "IfOrientation"
}

func (ImQuery) transformationType() string {
	return "ImQuery"
}

func (MaxColors) transformationType() string {
	return "MaxColors"
}

func (Mirror) transformationType() string {
	return "Mirror"
}

func (MonoHue) transformationType() string {
	return "MonoHue"
}

func (Opacity) transformationType() string {
	return "Opacity"
}

func (PointShapeType) shapeType() string {
	return "PointShapeType"
}

func (PolygonShapeType) shapeType() string {
	return "PolygonShapeType"
}

func (RectangleShapeType) shapeType() string {
	return "RectangleShapeType"
}

func (RegionOfInterestCrop) transformationType() string {
	return "RegionOfInterestCrop"
}

func (RelativeCrop) transformationType() string {
	return "RelativeCrop"
}

func (RemoveColor) transformationType() string {
	return "RemoveColor"
}

func (Resize) transformationType() string {
	return "Resize"
}

func (Rotate) transformationType() string {
	return "Rotate"
}

func (Scale) transformationType() string {
	return "Scale"
}

func (Shear) transformationType() string {
	return "Shear"
}

func (TextImageType) imageType() string {
	return "TextImageType"
}

func (TextImageTypePost) imageTypePost() string {
	return "TextImageTypePost"
}

func (Trim) transformationType() string {
	return "Trim"
}

func (URLImageType) imageType() string {
	return "URLImageType"
}

func (URLImageTypePost) imageTypePost() string {
	return "URLImageTypePost"
}

func (UnionShapeType) shapeType() string {
	return "UnionShapeType"
}

func (UnsharpMask) transformationType() string {
	return "UnsharpMask"
}

func (BackgroundColor) transformationTypePost() string {
	return "BackgroundColor"
}

func (Blur) transformationTypePost() string {
	return "Blur"
}

func (ChromaKey) transformationTypePost() string {
	return "ChromaKey"
}

func (CompositePost) transformationTypePost() string {
	return "CompositePost"
}

func (CompoundPost) transformationTypePost() string {
	return "CompoundPost"
}

func (Contrast) transformationTypePost() string {
	return "Contrast"
}

func (Goop) transformationTypePost() string {
	return "Goop"
}

func (Grayscale) transformationTypePost() string {
	return "Grayscale"
}

func (HSL) transformationTypePost() string {
	return "HSL"
}

func (HSV) transformationTypePost() string {
	return "HSV"
}

func (IfDimensionPost) transformationTypePost() string {
	return "IfDimensionPost"
}

func (IfOrientationPost) transformationTypePost() string {
	return "IfOrientationPost"
}

func (MaxColors) transformationTypePost() string {
	return "MaxColors"
}

func (Mirror) transformationTypePost() string {
	return "Mirror"
}

func (MonoHue) transformationTypePost() string {
	return "MonoHue"
}

func (Opacity) transformationTypePost() string {
	return "Opacity"
}

func (RemoveColor) transformationTypePost() string {
	return "RemoveColor"
}

func (UnsharpMask) transformationTypePost() string {
	return "UnsharpMask"
}

/*-----------------------------------------------*/
//////////////// Pointer functions ////////////////
/*-----------------------------------------------*/

// AppendGravityPriorityPtr returns pointer of AppendGravityPriority
func AppendGravityPriorityPtr(v AppendGravityPriority) *AppendGravityPriority {
	return &v
}

// AppendTransformationPtr returns pointer of AppendTransformation
func AppendTransformationPtr(v AppendTransformation) *AppendTransformation {
	return &v
}

// AspectCropTransformationPtr returns pointer of AspectCropTransformation
func AspectCropTransformationPtr(v AspectCropTransformation) *AspectCropTransformation {
	return &v
}

// BackgroundColorTransformationPtr returns pointer of BackgroundColorTransformation
func BackgroundColorTransformationPtr(v BackgroundColorTransformation) *BackgroundColorTransformation {
	return &v
}

// BlurTransformationPtr returns pointer of BlurTransformation
func BlurTransformationPtr(v BlurTransformation) *BlurTransformation {
	return &v
}

// BoxImageTypePostTypePtr returns pointer of BoxImageTypePostType
func BoxImageTypePostTypePtr(v BoxImageTypePostType) *BoxImageTypePostType {
	return &v
}

// BoxImageTypeTypePtr returns pointer of BoxImageTypeType
func BoxImageTypeTypePtr(v BoxImageTypeType) *BoxImageTypeType {
	return &v
}

// ChromaKeyTransformationPtr returns pointer of ChromaKeyTransformation
func ChromaKeyTransformationPtr(v ChromaKeyTransformation) *ChromaKeyTransformation {
	return &v
}

// CircleImageTypePostTypePtr returns pointer of CircleImageTypePostType
func CircleImageTypePostTypePtr(v CircleImageTypePostType) *CircleImageTypePostType {
	return &v
}

// CircleImageTypeTypePtr returns pointer of CircleImageTypeType
func CircleImageTypeTypePtr(v CircleImageTypeType) *CircleImageTypeType {
	return &v
}

// CompositePlacementPtr returns pointer of CompositePlacement
func CompositePlacementPtr(v CompositePlacement) *CompositePlacement {
	return &v
}

// CompositePostPlacementPtr returns pointer of CompositePostPlacement
func CompositePostPlacementPtr(v CompositePostPlacement) *CompositePostPlacement {
	return &v
}

// CompositePostScaleDimensionPtr returns pointer of CompositePostScaleDimension
func CompositePostScaleDimensionPtr(v CompositePostScaleDimension) *CompositePostScaleDimension {
	return &v
}

// CompositePostTransformationPtr returns pointer of CompositePostTransformation
func CompositePostTransformationPtr(v CompositePostTransformation) *CompositePostTransformation {
	return &v
}

// CompositeScaleDimensionPtr returns pointer of CompositeScaleDimension
func CompositeScaleDimensionPtr(v CompositeScaleDimension) *CompositeScaleDimension {
	return &v
}

// CompositeTransformationPtr returns pointer of CompositeTransformation
func CompositeTransformationPtr(v CompositeTransformation) *CompositeTransformation {
	return &v
}

// CompoundPostTransformationPtr returns pointer of CompoundPostTransformation
func CompoundPostTransformationPtr(v CompoundPostTransformation) *CompoundPostTransformation {
	return &v
}

// CompoundTransformationPtr returns pointer of CompoundTransformation
func CompoundTransformationPtr(v CompoundTransformation) *CompoundTransformation {
	return &v
}

// ContrastTransformationPtr returns pointer of ContrastTransformation
func ContrastTransformationPtr(v ContrastTransformation) *ContrastTransformation {
	return &v
}

// CropTransformationPtr returns pointer of CropTransformation
func CropTransformationPtr(v CropTransformation) *CropTransformation {
	return &v
}

// FaceCropAlgorithmPtr returns pointer of FaceCropAlgorithm
func FaceCropAlgorithmPtr(v FaceCropAlgorithm) *FaceCropAlgorithm {
	return &v
}

// FaceCropFocusPtr returns pointer of FaceCropFocus
func FaceCropFocusPtr(v FaceCropFocus) *FaceCropFocus {
	return &v
}

// FaceCropStylePtr returns pointer of FaceCropStyle
func FaceCropStylePtr(v FaceCropStyle) *FaceCropStyle {
	return &v
}

// FaceCropTransformationPtr returns pointer of FaceCropTransformation
func FaceCropTransformationPtr(v FaceCropTransformation) *FaceCropTransformation {
	return &v
}

// FeatureCropStylePtr returns pointer of FeatureCropStyle
func FeatureCropStylePtr(v FeatureCropStyle) *FeatureCropStyle {
	return &v
}

// FeatureCropTransformationPtr returns pointer of FeatureCropTransformation
func FeatureCropTransformationPtr(v FeatureCropTransformation) *FeatureCropTransformation {
	return &v
}

// FitAndFillTransformationPtr returns pointer of FitAndFillTransformation
func FitAndFillTransformationPtr(v FitAndFillTransformation) *FitAndFillTransformation {
	return &v
}

// GoopTransformationPtr returns pointer of GoopTransformation
func GoopTransformationPtr(v GoopTransformation) *GoopTransformation {
	return &v
}

// GravityPtr returns pointer of Gravity
func GravityPtr(v Gravity) *Gravity {
	return &v
}

// GravityPostPtr returns pointer of GravityPost
func GravityPostPtr(v GravityPost) *GravityPost {
	return &v
}

// GrayscaleTransformationPtr returns pointer of GrayscaleTransformation
func GrayscaleTransformationPtr(v GrayscaleTransformation) *GrayscaleTransformation {
	return &v
}

// GrayscaleTypePtr returns pointer of GrayscaleType
func GrayscaleTypePtr(v GrayscaleType) *GrayscaleType {
	return &v
}

// HSLTransformationPtr returns pointer of HSLTransformation
func HSLTransformationPtr(v HSLTransformation) *HSLTransformation {
	return &v
}

// HSVTransformationPtr returns pointer of HSVTransformation
func HSVTransformationPtr(v HSVTransformation) *HSVTransformation {
	return &v
}

// IfDimensionDimensionPtr returns pointer of IfDimensionDimension
func IfDimensionDimensionPtr(v IfDimensionDimension) *IfDimensionDimension {
	return &v
}

// IfDimensionPostDimensionPtr returns pointer of IfDimensionPostDimension
func IfDimensionPostDimensionPtr(v IfDimensionPostDimension) *IfDimensionPostDimension {
	return &v
}

// IfDimensionPostTransformationPtr returns pointer of IfDimensionPostTransformation
func IfDimensionPostTransformationPtr(v IfDimensionPostTransformation) *IfDimensionPostTransformation {
	return &v
}

// IfDimensionTransformationPtr returns pointer of IfDimensionTransformation
func IfDimensionTransformationPtr(v IfDimensionTransformation) *IfDimensionTransformation {
	return &v
}

// IfOrientationPostTransformationPtr returns pointer of IfOrientationPostTransformation
func IfOrientationPostTransformationPtr(v IfOrientationPostTransformation) *IfOrientationPostTransformation {
	return &v
}

// IfOrientationTransformationPtr returns pointer of IfOrientationTransformation
func IfOrientationTransformationPtr(v IfOrientationTransformation) *IfOrientationTransformation {
	return &v
}

// ImQueryAllowedTransformationsPtr returns pointer of ImQueryAllowedTransformations
func ImQueryAllowedTransformationsPtr(v ImQueryAllowedTransformations) *ImQueryAllowedTransformations {
	return &v
}

// ImQueryTransformationPtr returns pointer of ImQueryTransformation
func ImQueryTransformationPtr(v ImQueryTransformation) *ImQueryTransformation {
	return &v
}

// MaxColorsTransformationPtr returns pointer of MaxColorsTransformation
func MaxColorsTransformationPtr(v MaxColorsTransformation) *MaxColorsTransformation {
	return &v
}

// MirrorTransformationPtr returns pointer of MirrorTransformation
func MirrorTransformationPtr(v MirrorTransformation) *MirrorTransformation {
	return &v
}

// MonoHueTransformationPtr returns pointer of MonoHueTransformation
func MonoHueTransformationPtr(v MonoHueTransformation) *MonoHueTransformation {
	return &v
}

// OpacityTransformationPtr returns pointer of OpacityTransformation
func OpacityTransformationPtr(v OpacityTransformation) *OpacityTransformation {
	return &v
}

// OutputImageAllowedFormatsPtr returns pointer of OutputImageAllowedFormats
func OutputImageAllowedFormatsPtr(v OutputImageAllowedFormats) *OutputImageAllowedFormats {
	return &v
}

// OutputImageForcedFormatsPtr returns pointer of OutputImageForcedFormats
func OutputImageForcedFormatsPtr(v OutputImageForcedFormats) *OutputImageForcedFormats {
	return &v
}

// OutputImagePerceptualQualityPtr returns pointer of OutputImagePerceptualQuality
func OutputImagePerceptualQualityPtr(v OutputImagePerceptualQuality) *OutputImagePerceptualQuality {
	return &v
}

// RegionOfInterestCropStylePtr returns pointer of RegionOfInterestCropStyle
func RegionOfInterestCropStylePtr(v RegionOfInterestCropStyle) *RegionOfInterestCropStyle {
	return &v
}

// RegionOfInterestCropTransformationPtr returns pointer of RegionOfInterestCropTransformation
func RegionOfInterestCropTransformationPtr(v RegionOfInterestCropTransformation) *RegionOfInterestCropTransformation {
	return &v
}

// RelativeCropTransformationPtr returns pointer of RelativeCropTransformation
func RelativeCropTransformationPtr(v RelativeCropTransformation) *RelativeCropTransformation {
	return &v
}

// RemoveColorTransformationPtr returns pointer of RemoveColorTransformation
func RemoveColorTransformationPtr(v RemoveColorTransformation) *RemoveColorTransformation {
	return &v
}

// ResizeAspectPtr returns pointer of ResizeAspect
func ResizeAspectPtr(v ResizeAspect) *ResizeAspect {
	return &v
}

// ResizeTransformationPtr returns pointer of ResizeTransformation
func ResizeTransformationPtr(v ResizeTransformation) *ResizeTransformation {
	return &v
}

// ResizeTypePtr returns pointer of ResizeType
func ResizeTypePtr(v ResizeType) *ResizeType {
	return &v
}

// RotateTransformationPtr returns pointer of RotateTransformation
func RotateTransformationPtr(v RotateTransformation) *RotateTransformation {
	return &v
}

// ScaleTransformationPtr returns pointer of ScaleTransformation
func ScaleTransformationPtr(v ScaleTransformation) *ScaleTransformation {
	return &v
}

// ShearTransformationPtr returns pointer of ShearTransformation
func ShearTransformationPtr(v ShearTransformation) *ShearTransformation {
	return &v
}

// TextImageTypePostTypePtr returns pointer of TextImageTypePostType
func TextImageTypePostTypePtr(v TextImageTypePostType) *TextImageTypePostType {
	return &v
}

// TextImageTypeTypePtr returns pointer of TextImageTypeType
func TextImageTypeTypePtr(v TextImageTypeType) *TextImageTypeType {
	return &v
}

// TrimTransformationPtr returns pointer of TrimTransformation
func TrimTransformationPtr(v TrimTransformation) *TrimTransformation {
	return &v
}

// URLImageTypePostTypePtr returns pointer of URLImageTypePostType
func URLImageTypePostTypePtr(v URLImageTypePostType) *URLImageTypePostType {
	return &v
}

// URLImageTypeTypePtr returns pointer of URLImageTypeType
func URLImageTypeTypePtr(v URLImageTypeType) *URLImageTypeType {
	return &v
}

// UnsharpMaskTransformationPtr returns pointer of UnsharpMaskTransformation
func UnsharpMaskTransformationPtr(v UnsharpMaskTransformation) *UnsharpMaskTransformation {
	return &v
}

// VariableTypePtr returns pointer of VariableType
func VariableTypePtr(v VariableType) *VariableType {
	return &v
}

// OutputVideoPerceptualQualityPtr returns pointer of OutputVideoPerceptualQuality
func OutputVideoPerceptualQualityPtr(v OutputVideoPerceptualQuality) *OutputVideoPerceptualQuality {
	return &v
}

// OutputVideoVideoAdaptiveQualityPtr returns pointer of OutputVideoVideoAdaptiveQuality
func OutputVideoVideoAdaptiveQualityPtr(v OutputVideoVideoAdaptiveQuality) *OutputVideoVideoAdaptiveQuality {
	return &v
}

/*-----------------------------------------------*/
/////////////////// Validators ////////////////////
/*-----------------------------------------------*/

// Validate validates Append
func (a Append) Validate() error {
	return validation.Errors{
		"Gravity":         validation.Validate(a.Gravity),
		"GravityPriority": validation.Validate(a.GravityPriority),
		"Image": validation.Validate(a.Image,
			validation.Required,
		),
		"PreserveMinorDimension": validation.Validate(a.PreserveMinorDimension),
		"Transformation": validation.Validate(a.Transformation,
			validation.Required,
			validation.In(AppendTransformationAppend),
		),
	}.Filter()
}

// Validate validates AppendGravityPriorityVariableInline
func (a AppendGravityPriorityVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(a.Name),
		"Value": validation.Validate(a.Value,
			validation.In(AppendGravityPriorityHorizontal, AppendGravityPriorityVertical),
		),
	}.Filter()
}

// Validate validates AspectCrop
func (a AspectCrop) Validate() error {
	return validation.Errors{
		"AllowExpansion": validation.Validate(a.AllowExpansion),
		"Height":         validation.Validate(a.Height),
		"Transformation": validation.Validate(a.Transformation,
			validation.Required,
			validation.In(AspectCropTransformationAspectCrop),
		),
		"Width":     validation.Validate(a.Width),
		"XPosition": validation.Validate(a.XPosition),
		"YPosition": validation.Validate(a.YPosition),
	}.Filter()
}

// Validate validates BackgroundColor
func (b BackgroundColor) Validate() error {
	return validation.Errors{
		"Color": validation.Validate(b.Color,
			validation.Required,
		),
		"Transformation": validation.Validate(b.Transformation,
			validation.Required,
			validation.In(BackgroundColorTransformationBackgroundColor),
		),
	}.Filter()
}

// Validate validates Blur
func (b Blur) Validate() error {
	return validation.Errors{
		"Sigma": validation.Validate(b.Sigma),
		"Transformation": validation.Validate(b.Transformation,
			validation.Required,
			validation.In(BlurTransformationBlur),
		),
	}.Filter()
}

// Validate validates BooleanVariableInline
func (b BooleanVariableInline) Validate() error {
	return validation.Errors{
		"Name":  validation.Validate(b.Name),
		"Value": validation.Validate(b.Value),
	}.Filter()
}

// Validate validates BoxImageType
func (b BoxImageType) Validate() error {
	return validation.Errors{
		"Color":          validation.Validate(b.Color),
		"Height":         validation.Validate(b.Height),
		"Transformation": validation.Validate(b.Transformation),
		"Type": validation.Validate(b.Type,
			validation.Required,
			validation.In(BoxImageTypeTypeBox),
		),
		"Width": validation.Validate(b.Width),
	}.Filter()
}

// Validate validates BoxImageTypePost
func (b BoxImageTypePost) Validate() error {
	return validation.Errors{
		"Color":          validation.Validate(b.Color),
		"Height":         validation.Validate(b.Height),
		"Transformation": validation.Validate(b.Transformation),
		"Type": validation.Validate(b.Type,
			validation.Required,
			validation.In(BoxImageTypePostTypeBox),
		),
		"Width": validation.Validate(b.Width),
	}.Filter()
}

// Validate validates Breakpoints
func (b Breakpoints) Validate() error {
	return validation.Errors{
		"Widths": validation.Validate(b.Widths, validation.Each()),
	}.Filter()
}

// Validate validates ChromaKey
func (c ChromaKey) Validate() error {
	return validation.Errors{
		"Hue":                 validation.Validate(c.Hue),
		"HueFeather":          validation.Validate(c.HueFeather),
		"HueTolerance":        validation.Validate(c.HueTolerance),
		"LightnessFeather":    validation.Validate(c.LightnessFeather),
		"LightnessTolerance":  validation.Validate(c.LightnessTolerance),
		"SaturationFeather":   validation.Validate(c.SaturationFeather),
		"SaturationTolerance": validation.Validate(c.SaturationTolerance),
		"Transformation": validation.Validate(c.Transformation,
			validation.Required,
			validation.In(ChromaKeyTransformationChromaKey),
		),
	}.Filter()
}

// Validate validates CircleImageType
func (c CircleImageType) Validate() error {
	return validation.Errors{
		"Color":          validation.Validate(c.Color),
		"Diameter":       validation.Validate(c.Diameter),
		"Transformation": validation.Validate(c.Transformation),
		"Type": validation.Validate(c.Type,
			validation.Required,
			validation.In(CircleImageTypeTypeCircle),
		),
		"Width": validation.Validate(c.Width),
	}.Filter()
}

// Validate validates CircleImageTypePost
func (c CircleImageTypePost) Validate() error {
	return validation.Errors{
		"Color":          validation.Validate(c.Color),
		"Diameter":       validation.Validate(c.Diameter),
		"Transformation": validation.Validate(c.Transformation),
		"Type": validation.Validate(c.Type,
			validation.Required,
			validation.In(CircleImageTypePostTypeCircle),
		),
		"Width": validation.Validate(c.Width),
	}.Filter()
}

// Validate validates CircleShapeType
func (c CircleShapeType) Validate() error {
	return validation.Errors{
		"Center": validation.Validate(c.Center,
			validation.Required,
		),
		"Radius": validation.Validate(c.Radius,
			validation.Required,
		),
	}.Filter()
}

// Validate validates Composite
func (c Composite) Validate() error {
	return validation.Errors{
		"Gravity": validation.Validate(c.Gravity),
		"Image": validation.Validate(c.Image,
			validation.Required,
		),
		"Placement":      validation.Validate(c.Placement),
		"Scale":          validation.Validate(c.Scale),
		"ScaleDimension": validation.Validate(c.ScaleDimension),
		"Transformation": validation.Validate(c.Transformation,
			validation.Required,
			validation.In(CompositeTransformationComposite),
		),
		"XPosition": validation.Validate(c.XPosition),
		"YPosition": validation.Validate(c.YPosition),
	}.Filter()
}

// Validate validates CompositePlacementVariableInline
func (c CompositePlacementVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(c.Name),
		"Value": validation.Validate(c.Value,
			validation.In(CompositePlacementOver,
				CompositePlacementUnder,
				CompositePlacementMask,
				CompositePlacementStencil),
		),
	}.Filter()
}

// Validate validates CompositePost
func (c CompositePost) Validate() error {
	return validation.Errors{
		"Gravity": validation.Validate(c.Gravity),
		"Image": validation.Validate(c.Image,
			validation.Required,
		),
		"Placement":      validation.Validate(c.Placement),
		"Scale":          validation.Validate(c.Scale),
		"ScaleDimension": validation.Validate(c.ScaleDimension),
		"Transformation": validation.Validate(c.Transformation,
			validation.Required,
			validation.In(CompositePostTransformationComposite),
		),
		"XPosition": validation.Validate(c.XPosition),
		"YPosition": validation.Validate(c.YPosition),
	}.Filter()
}

// Validate validates CompositePostPlacementVariableInline
func (c CompositePostPlacementVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(c.Name),
		"Value": validation.Validate(c.Value,
			validation.In(CompositePostPlacementOver,
				CompositePostPlacementUnder,
				CompositePostPlacementMask,
				CompositePostPlacementStencil),
		),
	}.Filter()
}

// Validate validates CompositePostScaleDimensionVariableInline
func (c CompositePostScaleDimensionVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(c.Name),
		"Value": validation.Validate(c.Value,
			validation.In(CompositePostScaleDimensionWidth, CompositePostScaleDimensionHeight),
		),
	}.Filter()
}

// Validate validates CompositeScaleDimensionVariableInline
func (c CompositeScaleDimensionVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(c.Name),
		"Value": validation.Validate(c.Value,
			validation.In(CompositeScaleDimensionWidth, CompositeScaleDimensionHeight),
		),
	}.Filter()
}

// Validate validates Compound
func (c Compound) Validate() error {
	return validation.Errors{
		"Transformation": validation.Validate(c.Transformation,
			validation.Required,
			validation.In(CompoundTransformationCompound),
		),
		"Transformations": validation.Validate(c.Transformations),
	}.Filter()
}

// Validate validates CompoundPost
func (c CompoundPost) Validate() error {
	return validation.Errors{
		"Transformation": validation.Validate(c.Transformation,
			validation.Required,
			validation.In(CompoundPostTransformationCompound),
		),
		"Transformations": validation.Validate(c.Transformations),
	}.Filter()
}

// Validate validates Contrast
func (c Contrast) Validate() error {
	return validation.Errors{
		"Brightness": validation.Validate(c.Brightness),
		"Contrast":   validation.Validate(c.Contrast),
		"Transformation": validation.Validate(c.Transformation,
			validation.Required,
			validation.In(ContrastTransformationContrast),
		),
	}.Filter()
}

// Validate validates Crop
func (c Crop) Validate() error {
	return validation.Errors{
		"AllowExpansion": validation.Validate(c.AllowExpansion),
		"Gravity":        validation.Validate(c.Gravity),
		"Height": validation.Validate(c.Height,
			validation.Required,
		),
		"Transformation": validation.Validate(c.Transformation,
			validation.Required,
			validation.In(CropTransformationCrop),
		),
		"Width": validation.Validate(c.Width,
			validation.Required,
		),
		"XPosition": validation.Validate(c.XPosition),
		"YPosition": validation.Validate(c.YPosition),
	}.Filter()
}

// Validate validates EnumOptions
func (e EnumOptions) Validate() error {
	return validation.Errors{
		"ID": validation.Validate(e.ID,
			validation.Required,
		),
		"Value": validation.Validate(e.Value,
			validation.Required,
		),
	}.Filter()
}

// Validate validates FaceCrop
func (f FaceCrop) Validate() error {
	return validation.Errors{
		"Algorithm":   validation.Validate(f.Algorithm),
		"Confidence":  validation.Validate(f.Confidence),
		"FailGravity": validation.Validate(f.FailGravity),
		"Focus":       validation.Validate(f.Focus),
		"Gravity":     validation.Validate(f.Gravity),
		"Height": validation.Validate(f.Height,
			validation.Required,
		),
		"Padding": validation.Validate(f.Padding),
		"Style":   validation.Validate(f.Style),
		"Transformation": validation.Validate(f.Transformation,
			validation.Required,
			validation.In(FaceCropTransformationFaceCrop),
		),
		"Width": validation.Validate(f.Width,
			validation.Required,
		),
	}.Filter()
}

// Validate validates FaceCropAlgorithmVariableInline
func (f FaceCropAlgorithmVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(f.Name),
		"Value": validation.Validate(f.Value,
			validation.In(FaceCropAlgorithmCascade, FaceCropAlgorithmDnn),
		),
	}.Filter()
}

// Validate validates FaceCropFocusVariableInline
func (f FaceCropFocusVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(f.Name),
		"Value": validation.Validate(f.Value,
			validation.In(FaceCropFocusAllFaces, FaceCropFocusBiggestFace),
		),
	}.Filter()
}

// Validate validates FaceCropStyleVariableInline
func (f FaceCropStyleVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(f.Name),
		"Value": validation.Validate(f.Value,
			validation.In(FaceCropStyleCrop, FaceCropStyleFill, FaceCropStyleZoom),
		),
	}.Filter()
}

// Validate validates FeatureCrop
func (f FeatureCrop) Validate() error {
	return validation.Errors{
		"FailGravity":   validation.Validate(f.FailGravity),
		"FeatureRadius": validation.Validate(f.FeatureRadius),
		"Gravity":       validation.Validate(f.Gravity),
		"Height": validation.Validate(f.Height,
			validation.Required,
		),
		"MaxFeatures":       validation.Validate(f.MaxFeatures),
		"MinFeatureQuality": validation.Validate(f.MinFeatureQuality),
		"Padding":           validation.Validate(f.Padding),
		"Style":             validation.Validate(f.Style),
		"Transformation": validation.Validate(f.Transformation,
			validation.Required,
			validation.In(FeatureCropTransformationFeatureCrop),
		),
		"Width": validation.Validate(f.Width,
			validation.Required,
		),
	}.Filter()
}

// Validate validates FeatureCropStyleVariableInline
func (f FeatureCropStyleVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(f.Name),
		"Value": validation.Validate(f.Value,
			validation.In(FeatureCropStyleCrop, FeatureCropStyleFill, FeatureCropStyleZoom),
		),
	}.Filter()
}

// Validate validates FitAndFill
func (f FitAndFill) Validate() error {
	return validation.Errors{
		"FillTransformation": validation.Validate(f.FillTransformation),
		"Height": validation.Validate(f.Height,
			validation.Required,
		),
		"Transformation": validation.Validate(f.Transformation,
			validation.Required,
			validation.In(FitAndFillTransformationFitAndFill),
		),
		"Width": validation.Validate(f.Width,
			validation.Required,
		),
	}.Filter()
}

// Validate validates Goop
func (g Goop) Validate() error {
	return validation.Errors{
		"Chaos":   validation.Validate(g.Chaos),
		"Density": validation.Validate(g.Density),
		"Power":   validation.Validate(g.Power),
		"Seed":    validation.Validate(g.Seed),
		"Transformation": validation.Validate(g.Transformation,
			validation.Required,
			validation.In(GoopTransformationGoop),
		),
	}.Filter()
}

// Validate validates GravityPostVariableInline
func (g GravityPostVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(g.Name),
		"Value": validation.Validate(g.Value,
			validation.In(GravityPostNorth,
				GravityPostNorthEast,
				GravityPostNorthWest,
				GravityPostSouth,
				GravityPostSouthEast,
				GravityPostSouthWest,
				GravityPostCenter,
				GravityPostEast,
				GravityPostWest),
		),
	}.Filter()
}

// Validate validates GravityVariableInline
func (g GravityVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(g.Name),
		"Value": validation.Validate(g.Value,
			validation.In(GravityNorth,
				GravityNorthEast,
				GravityNorthWest,
				GravitySouth,
				GravitySouthEast,
				GravitySouthWest,
				GravityCenter,
				GravityEast,
				GravityWest),
		),
	}.Filter()
}

// Validate validates Grayscale
func (g Grayscale) Validate() error {
	return validation.Errors{
		"Transformation": validation.Validate(g.Transformation,
			validation.Required,
			validation.In(GrayscaleTransformationGrayscale),
		),
		"Type": validation.Validate(g.Type),
	}.Filter()
}

// Validate validates GrayscaleTypeVariableInline
func (g GrayscaleTypeVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(g.Name),
		"Value": validation.Validate(g.Value,
			validation.In(GrayscaleTypeRec601,
				GrayscaleTypeRec709,
				GrayscaleTypeBrightness,
				GrayscaleTypeLightness),
		),
	}.Filter()
}

// Validate validates HSL
func (h HSL) Validate() error {
	return validation.Errors{
		"Hue":        validation.Validate(h.Hue),
		"Lightness":  validation.Validate(h.Lightness),
		"Saturation": validation.Validate(h.Saturation),
		"Transformation": validation.Validate(h.Transformation,
			validation.Required,
			validation.In(HSLTransformationHSL),
		),
	}.Filter()
}

// Validate validates HSV
func (h HSV) Validate() error {
	return validation.Errors{
		"Hue":        validation.Validate(h.Hue),
		"Saturation": validation.Validate(h.Saturation),
		"Transformation": validation.Validate(h.Transformation,
			validation.Required,
			validation.In(HSVTransformationHSV),
		),
		"Value": validation.Validate(h.Value),
	}.Filter()
}

// Validate validates IfDimension
func (i IfDimension) Validate() error {
	return validation.Errors{
		"Default":     validation.Validate(i.Default),
		"Dimension":   validation.Validate(i.Dimension),
		"Equal":       validation.Validate(i.Equal),
		"GreaterThan": validation.Validate(i.GreaterThan),
		"LessThan":    validation.Validate(i.LessThan),
		"Transformation": validation.Validate(i.Transformation,
			validation.Required,
			validation.In(IfDimensionTransformationIfDimension),
		),
		"Value": validation.Validate(i.Value,
			validation.Required,
		),
	}.Filter()
}

// Validate validates IfDimensionDimensionVariableInline
func (i IfDimensionDimensionVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(i.Name),
		"Value": validation.Validate(i.Value,
			validation.In(IfDimensionDimensionWidth, IfDimensionDimensionHeight, IfDimensionDimensionBoth),
		),
	}.Filter()
}

// Validate validates IfDimensionPost
func (i IfDimensionPost) Validate() error {
	return validation.Errors{
		"Default":     validation.Validate(i.Default),
		"Dimension":   validation.Validate(i.Dimension),
		"Equal":       validation.Validate(i.Equal),
		"GreaterThan": validation.Validate(i.GreaterThan),
		"LessThan":    validation.Validate(i.LessThan),
		"Transformation": validation.Validate(i.Transformation,
			validation.Required,
			validation.In(IfDimensionPostTransformationIfDimension),
		),
		"Value": validation.Validate(i.Value,
			validation.Required,
		),
	}.Filter()
}

// Validate validates IfDimensionPostDimensionVariableInline
func (i IfDimensionPostDimensionVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(i.Name),
		"Value": validation.Validate(i.Value,
			validation.In(IfDimensionPostDimensionWidth, IfDimensionPostDimensionHeight, IfDimensionPostDimensionBoth),
		),
	}.Filter()
}

// Validate validates IfOrientation
func (i IfOrientation) Validate() error {
	return validation.Errors{
		"Default":   validation.Validate(i.Default),
		"Landscape": validation.Validate(i.Landscape),
		"Portrait":  validation.Validate(i.Portrait),
		"Square":    validation.Validate(i.Square),
		"Transformation": validation.Validate(i.Transformation,
			validation.Required,
			validation.In(IfOrientationTransformationIfOrientation),
		),
	}.Filter()
}

// Validate validates IfOrientationPost
func (i IfOrientationPost) Validate() error {
	return validation.Errors{
		"Default":   validation.Validate(i.Default),
		"Landscape": validation.Validate(i.Landscape),
		"Portrait":  validation.Validate(i.Portrait),
		"Square":    validation.Validate(i.Square),
		"Transformation": validation.Validate(i.Transformation,
			validation.Required,
			validation.In(IfOrientationPostTransformationIfOrientation),
		),
	}.Filter()
}

// Validate validates ImQuery
func (i ImQuery) Validate() error {
	return validation.Errors{
		"AllowedTransformations": validation.Validate(i.AllowedTransformations,
			validation.Required, validation.Each(
				validation.In(ImQueryAllowedTransformationsAppend,
					ImQueryAllowedTransformationsAspectCrop,
					ImQueryAllowedTransformationsBackgroundColor,
					ImQueryAllowedTransformationsBlur,
					ImQueryAllowedTransformationsComposite,
					ImQueryAllowedTransformationsContrast,
					ImQueryAllowedTransformationsCrop,
					ImQueryAllowedTransformationsChromaKey,
					ImQueryAllowedTransformationsFaceCrop,
					ImQueryAllowedTransformationsFeatureCrop,
					ImQueryAllowedTransformationsFitAndFill,
					ImQueryAllowedTransformationsGoop,
					ImQueryAllowedTransformationsGrayscale,
					ImQueryAllowedTransformationsHSL,
					ImQueryAllowedTransformationsHSV,
					ImQueryAllowedTransformationsMaxColors,
					ImQueryAllowedTransformationsMirror,
					ImQueryAllowedTransformationsMonoHue,
					ImQueryAllowedTransformationsOpacity,
					ImQueryAllowedTransformationsRegionOfInterestCrop,
					ImQueryAllowedTransformationsRelativeCrop,
					ImQueryAllowedTransformationsRemoveColor,
					ImQueryAllowedTransformationsResize,
					ImQueryAllowedTransformationsRotate,
					ImQueryAllowedTransformationsScale,
					ImQueryAllowedTransformationsShear,
					ImQueryAllowedTransformationsTrim,
					ImQueryAllowedTransformationsUnsharpMask,
					ImQueryAllowedTransformationsIfDimension,
					ImQueryAllowedTransformationsIfOrientation)),
		),
		"Query": validation.Validate(i.Query,
			validation.Required,
		),
		"Transformation": validation.Validate(i.Transformation,
			validation.Required,
			validation.In(ImQueryTransformationImQuery),
		),
	}.Filter()
}

// Validate validates IntegerVariableInline
func (i IntegerVariableInline) Validate() error {
	return validation.Errors{
		"Name":  validation.Validate(i.Name),
		"Value": validation.Validate(i.Value),
	}.Filter()
}

// Validate validates MaxColors
func (m MaxColors) Validate() error {
	return validation.Errors{
		"Colors": validation.Validate(m.Colors,
			validation.Required,
		),
		"Transformation": validation.Validate(m.Transformation,
			validation.Required,
			validation.In(MaxColorsTransformationMaxColors),
		),
	}.Filter()
}

// Validate validates Mirror
func (m Mirror) Validate() error {
	return validation.Errors{
		"Horizontal": validation.Validate(m.Horizontal),
		"Transformation": validation.Validate(m.Transformation,
			validation.Required,
			validation.In(MirrorTransformationMirror),
		),
		"Vertical": validation.Validate(m.Vertical),
	}.Filter()
}

// Validate validates MonoHue
func (m MonoHue) Validate() error {
	return validation.Errors{
		"Hue": validation.Validate(m.Hue),
		"Transformation": validation.Validate(m.Transformation,
			validation.Required,
			validation.In(MonoHueTransformationMonoHue),
		),
	}.Filter()
}

// Validate validates NumberVariableInline
func (n NumberVariableInline) Validate() error {
	return validation.Errors{
		"Name":  validation.Validate(n.Name),
		"Value": validation.Validate(n.Value),
	}.Filter()
}

// Validate validates Opacity
func (o Opacity) Validate() error {
	return validation.Errors{
		"Opacity": validation.Validate(o.Opacity,
			validation.Required,
		),
		"Transformation": validation.Validate(o.Transformation,
			validation.Required,
			validation.In(OpacityTransformationOpacity),
		),
	}.Filter()
}

// Validate validates OutputImage
func (o OutputImage) Validate() error {
	return validation.Errors{
		"AdaptiveQuality": validation.Validate(o.AdaptiveQuality,
			validation.Min(1),
			validation.Max(100),
		),
		"AllowedFormats": validation.Validate(o.AllowedFormats, validation.Each(
			validation.In(OutputImageAllowedFormatsGif,
				OutputImageAllowedFormatsJpeg,
				OutputImageAllowedFormatsPng,
				OutputImageAllowedFormatsWebp,
				OutputImageAllowedFormatsJpegxr,
				OutputImageAllowedFormatsJpeg2000)),
		),
		"ForcedFormats": validation.Validate(o.ForcedFormats, validation.Each(
			validation.In(OutputImageForcedFormatsGif,
				OutputImageForcedFormatsJpeg,
				OutputImageForcedFormatsPng,
				OutputImageForcedFormatsWebp,
				OutputImageForcedFormatsJpegxr,
				OutputImageForcedFormatsJpeg2000)),
		),
		"PerceptualQuality": validation.Validate(o.PerceptualQuality),
		"PerceptualQualityFloor": validation.Validate(o.PerceptualQualityFloor,
			validation.Min(1),
			validation.Max(100),
		),
		"Quality": validation.Validate(o.Quality),
	}.Filter()
}

// Validate validates OutputImagePerceptualQualityVariableInline
func (o OutputImagePerceptualQualityVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(o.Name),
		"Value": validation.Validate(o.Value,
			validation.In(OutputImagePerceptualQualityHigh,
				OutputImagePerceptualQualityMediumHigh,
				OutputImagePerceptualQualityMedium,
				OutputImagePerceptualQualityMediumLow,
				OutputImagePerceptualQualityLow),
		),
	}.Filter()
}

// Validate validates PointShapeType
func (p PointShapeType) Validate() error {
	return validation.Errors{
		"X": validation.Validate(p.X,
			validation.Required,
		),
		"Y": validation.Validate(p.Y,
			validation.Required,
		),
	}.Filter()
}

// Validate validates PolicyOutputImage
func (p PolicyOutputImage) Validate() error {
	return validation.Errors{
		"Breakpoints": validation.Validate(p.Breakpoints),
		"DateCreated": validation.Validate(p.DateCreated,
			validation.Required,
		),
		"Hosts": validation.Validate(p.Hosts, validation.Each()),
		"ID": validation.Validate(p.ID,
			validation.Required,
		),
		"Output":                        validation.Validate(p.Output),
		"PostBreakpointTransformations": validation.Validate(p.PostBreakpointTransformations),
		"PreviousVersion": validation.Validate(p.PreviousVersion,
			validation.Required,
		),
		"RolloutInfo": validation.Validate(p.RolloutInfo,
			validation.Required,
		),
		"Transformations": validation.Validate(p.Transformations),
		"User": validation.Validate(p.User,
			validation.Required,
		),
		"Variables": validation.Validate(p.Variables, validation.Each()),
		"Version": validation.Validate(p.Version,
			validation.Required,
		),
		"Video": validation.Validate(p.Video,
			validation.In(PolicyOutputImageVideoFalse),
		),
	}.Filter()
}

// Validate validates PolygonShapeType
func (p PolygonShapeType) Validate() error {
	return validation.Errors{
		"Points": validation.Validate(p.Points,
			validation.Required, validation.Each(),
		),
	}.Filter()
}

// Validate validates QueryVariableInline
func (q QueryVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(q.Name,
			validation.Required,
		),
	}.Filter()
}

// Validate validates RectangleShapeType
func (r RectangleShapeType) Validate() error {
	return validation.Errors{
		"Anchor": validation.Validate(r.Anchor,
			validation.Required,
		),
		"Height": validation.Validate(r.Height,
			validation.Required,
		),
		"Width": validation.Validate(r.Width,
			validation.Required,
		),
	}.Filter()
}

// Validate validates RegionOfInterestCrop
func (r RegionOfInterestCrop) Validate() error {
	return validation.Errors{
		"Gravity": validation.Validate(r.Gravity),
		"Height": validation.Validate(r.Height,
			validation.Required,
		),
		"RegionOfInterest": validation.Validate(r.RegionOfInterest,
			validation.Required,
		),
		"Style": validation.Validate(r.Style),
		"Transformation": validation.Validate(r.Transformation,
			validation.Required,
			validation.In(RegionOfInterestCropTransformationRegionOfInterestCrop),
		),
		"Width": validation.Validate(r.Width,
			validation.Required,
		),
	}.Filter()
}

// Validate validates RegionOfInterestCropStyleVariableInline
func (r RegionOfInterestCropStyleVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(r.Name),
		"Value": validation.Validate(r.Value,
			validation.In(RegionOfInterestCropStyleCrop, RegionOfInterestCropStyleFill, RegionOfInterestCropStyleZoom),
		),
	}.Filter()
}

// Validate validates RelativeCrop
func (r RelativeCrop) Validate() error {
	return validation.Errors{
		"East":  validation.Validate(r.East),
		"North": validation.Validate(r.North),
		"South": validation.Validate(r.South),
		"Transformation": validation.Validate(r.Transformation,
			validation.Required,
			validation.In(RelativeCropTransformationRelativeCrop),
		),
		"West": validation.Validate(r.West),
	}.Filter()
}

// Validate validates RemoveColor
func (r RemoveColor) Validate() error {
	return validation.Errors{
		"Color": validation.Validate(r.Color,
			validation.Required,
		),
		"Feather":   validation.Validate(r.Feather),
		"Tolerance": validation.Validate(r.Tolerance),
		"Transformation": validation.Validate(r.Transformation,
			validation.Required,
			validation.In(RemoveColorTransformationRemoveColor),
		),
	}.Filter()
}

// Validate validates Resize
func (r Resize) Validate() error {
	return validation.Errors{
		"Aspect": validation.Validate(r.Aspect),
		"Height": validation.Validate(r.Height),
		"Transformation": validation.Validate(r.Transformation,
			validation.Required,
			validation.In(ResizeTransformationResize),
		),
		"Type":  validation.Validate(r.Type),
		"Width": validation.Validate(r.Width),
	}.Filter()
}

// Validate validates ResizeAspectVariableInline
func (r ResizeAspectVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(r.Name),
		"Value": validation.Validate(r.Value,
			validation.In(ResizeAspectFit, ResizeAspectFill, ResizeAspectIgnore),
		),
	}.Filter()
}

// Validate validates ResizeTypeVariableInline
func (r ResizeTypeVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(r.Name),
		"Value": validation.Validate(r.Value,
			validation.In(ResizeTypeNormal, ResizeTypeUpsize, ResizeTypeDownsize),
		),
	}.Filter()
}

// Validate validates RolloutInfo
func (r RolloutInfo) Validate() error {
	return validation.Errors{
		"EndTime": validation.Validate(r.EndTime,
			validation.Required,
		),
		"RolloutDuration": validation.Validate(r.RolloutDuration,
			validation.Required,
			validation.Min(3600),
			validation.Max(604800),
		),
		"StartTime": validation.Validate(r.StartTime,
			validation.Required,
		),
	}.Filter()
}

// Validate validates Rotate
func (r Rotate) Validate() error {
	return validation.Errors{
		"Degrees": validation.Validate(r.Degrees,
			validation.Required,
		),
		"Transformation": validation.Validate(r.Transformation,
			validation.Required,
			validation.In(RotateTransformationRotate),
		),
	}.Filter()
}

// Validate validates Scale
func (s Scale) Validate() error {
	return validation.Errors{
		"Height": validation.Validate(s.Height,
			validation.Required,
		),
		"Transformation": validation.Validate(s.Transformation,
			validation.Required,
			validation.In(ScaleTransformationScale),
		),
		"Width": validation.Validate(s.Width,
			validation.Required,
		),
	}.Filter()
}

// Validate validates Shear
func (s Shear) Validate() error {
	return validation.Errors{
		"Transformation": validation.Validate(s.Transformation,
			validation.Required,
			validation.In(ShearTransformationShear),
		),
		"XShear": validation.Validate(s.XShear),
		"YShear": validation.Validate(s.YShear),
	}.Filter()
}

// Validate validates StringVariableInline
func (s StringVariableInline) Validate() error {
	return validation.Errors{
		"Name":  validation.Validate(s.Name),
		"Value": validation.Validate(s.Value),
	}.Filter()
}

// Validate validates TextImageType
func (t TextImageType) Validate() error {
	return validation.Errors{
		"Fill":       validation.Validate(t.Fill),
		"Size":       validation.Validate(t.Size),
		"Stroke":     validation.Validate(t.Stroke),
		"StrokeSize": validation.Validate(t.StrokeSize),
		"Text": validation.Validate(t.Text,
			validation.Required,
		),
		"Transformation": validation.Validate(t.Transformation),
		"Type": validation.Validate(t.Type,
			validation.Required,
			validation.In(TextImageTypeTypeText),
		),
		"Typeface": validation.Validate(t.Typeface),
	}.Filter()
}

// Validate validates TextImageTypePost
func (t TextImageTypePost) Validate() error {
	return validation.Errors{
		"Fill":       validation.Validate(t.Fill),
		"Size":       validation.Validate(t.Size),
		"Stroke":     validation.Validate(t.Stroke),
		"StrokeSize": validation.Validate(t.StrokeSize),
		"Text": validation.Validate(t.Text,
			validation.Required,
		),
		"Transformation": validation.Validate(t.Transformation),
		"Type": validation.Validate(t.Type,
			validation.Required,
			validation.In(TextImageTypePostTypeText),
		),
		"Typeface": validation.Validate(t.Typeface),
	}.Filter()
}

// Validate validates Trim
func (t Trim) Validate() error {
	return validation.Errors{
		"Fuzz": validation.Validate(t.Fuzz,
			validation.Max(1),
		),
		"Padding": validation.Validate(t.Padding),
		"Transformation": validation.Validate(t.Transformation,
			validation.Required,
			validation.In(TrimTransformationTrim),
		),
	}.Filter()
}

// Validate validates URLImageType
func (u URLImageType) Validate() error {
	return validation.Errors{
		"Transformation": validation.Validate(u.Transformation),
		"Type": validation.Validate(u.Type,
			validation.In(URLImageTypeTypeURL),
		),
		"URL": validation.Validate(u.URL,
			validation.Required,
		),
	}.Filter()
}

// Validate validates URLImageTypePost
func (u URLImageTypePost) Validate() error {
	return validation.Errors{
		"Transformation": validation.Validate(u.Transformation),
		"Type": validation.Validate(u.Type,
			validation.In(URLImageTypePostTypeURL),
		),
		"URL": validation.Validate(u.URL,
			validation.Required,
		),
	}.Filter()
}

// Validate validates UnionShapeType
func (u UnionShapeType) Validate() error {
	return validation.Errors{
		"Shapes": validation.Validate(u.Shapes,
			validation.Required, validation.Each(),
		),
	}.Filter()
}

// Validate validates UnsharpMask
func (u UnsharpMask) Validate() error {
	return validation.Errors{
		"Gain":      validation.Validate(u.Gain),
		"Sigma":     validation.Validate(u.Sigma),
		"Threshold": validation.Validate(u.Threshold),
		"Transformation": validation.Validate(u.Transformation,
			validation.Required,
			validation.In(UnsharpMaskTransformationUnsharpMask),
		),
	}.Filter()
}

// Validate validates Variable
func (v Variable) Validate() error {
	return validation.Errors{
		"DefaultValue": validation.Validate(v.DefaultValue,
			validation.Required.When(v.Type != VariableTypeString),
		),
		"EnumOptions": validation.Validate(v.EnumOptions, validation.Each()),
		"Name": validation.Validate(v.Name,
			validation.Required,
		),
		"Postfix": validation.Validate(v.Postfix),
		"Prefix":  validation.Validate(v.Prefix),
		"Type": validation.Validate(v.Type,
			validation.Required,
			validation.In(VariableTypeBool,
				VariableTypeNumber,
				VariableTypeURL,
				VariableTypeColor,
				VariableTypeGravity,
				VariableTypePlacement,
				VariableTypeScaleDimension,
				VariableTypeGrayscaleType,
				VariableTypeAspect,
				VariableTypeResizeType,
				VariableTypeDimension,
				VariableTypePerceptualQuality,
				VariableTypeString,
				VariableTypeFocus),
		),
	}.Filter()
}

// Validate validates VariableInline
func (v VariableInline) Validate() error {
	return validation.Errors{
		"Var": validation.Validate(v.Var,
			validation.Required,
		),
	}.Filter()
}

// Validate validates OutputVideo
func (o OutputVideo) Validate() error {
	return validation.Errors{
		"PerceptualQuality":    validation.Validate(o.PerceptualQuality),
		"PlaceholderVideoURL":  validation.Validate(o.PlaceholderVideoURL),
		"VideoAdaptiveQuality": validation.Validate(o.VideoAdaptiveQuality),
	}.Filter()
}

// Validate validates OutputVideoPerceptualQualityVariableInline
func (o OutputVideoPerceptualQualityVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(o.Name),
		"Value": validation.Validate(o.Value,
			validation.In(OutputVideoPerceptualQualityHigh,
				OutputVideoPerceptualQualityMediumHigh,
				OutputVideoPerceptualQualityMedium,
				OutputVideoPerceptualQualityMediumLow,
				OutputVideoPerceptualQualityLow),
		),
	}.Filter()
}

// Validate validates OutputVideoVideoAdaptiveQualityVariableInline
func (o OutputVideoVideoAdaptiveQualityVariableInline) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(o.Name),
		"Value": validation.Validate(o.Value,
			validation.In(OutputVideoVideoAdaptiveQualityHigh,
				OutputVideoVideoAdaptiveQualityMediumHigh,
				OutputVideoVideoAdaptiveQualityMedium,
				OutputVideoVideoAdaptiveQualityMediumLow,
				OutputVideoVideoAdaptiveQualityLow),
		),
	}.Filter()
}

// Validate validates PolicyOutputVideo
func (p PolicyOutputVideo) Validate() error {
	return validation.Errors{
		"Breakpoints": validation.Validate(p.Breakpoints),
		"DateCreated": validation.Validate(p.DateCreated,
			validation.Required,
		),
		"Hosts": validation.Validate(p.Hosts, validation.Each()),
		"ID": validation.Validate(p.ID,
			validation.Required,
		),
		"Output": validation.Validate(p.Output),
		"PreviousVersion": validation.Validate(p.PreviousVersion,
			validation.Required,
		),
		"RolloutInfo": validation.Validate(p.RolloutInfo,
			validation.Required,
		),
		"User": validation.Validate(p.User,
			validation.Required,
		),
		"Variables": validation.Validate(p.Variables, validation.Each()),
		"Version": validation.Validate(p.Version,
			validation.Required,
		),
		"Video": validation.Validate(p.Video,
			validation.In(PolicyOutputVideoVideoTrue),
		),
	}.Filter()
}

/*-----------------------------------------------*/
//// Variable type marshalers and unmarshalers ////
/*-----------------------------------------------*/
var (

	// ErrUnmarshalVariableAppendGravityPriorityVariableInline represents an error while unmarshalling AppendGravityPriorityVariableInline
	ErrUnmarshalVariableAppendGravityPriorityVariableInline = errors.New("unmarshalling AppendGravityPriorityVariableInline")
	// ErrUnmarshalVariableBooleanVariableInline represents an error while unmarshalling BooleanVariableInline
	ErrUnmarshalVariableBooleanVariableInline = errors.New("unmarshalling BooleanVariableInline")
	// ErrUnmarshalVariableCompositePlacementVariableInline represents an error while unmarshalling CompositePlacementVariableInline
	ErrUnmarshalVariableCompositePlacementVariableInline = errors.New("unmarshalling CompositePlacementVariableInline")
	// ErrUnmarshalVariableCompositePostPlacementVariableInline represents an error while unmarshalling CompositePostPlacementVariableInline
	ErrUnmarshalVariableCompositePostPlacementVariableInline = errors.New("unmarshalling CompositePostPlacementVariableInline")
	// ErrUnmarshalVariableCompositePostScaleDimensionVariableInline represents an error while unmarshalling CompositePostScaleDimensionVariableInline
	ErrUnmarshalVariableCompositePostScaleDimensionVariableInline = errors.New("unmarshalling CompositePostScaleDimensionVariableInline")
	// ErrUnmarshalVariableCompositeScaleDimensionVariableInline represents an error while unmarshalling CompositeScaleDimensionVariableInline
	ErrUnmarshalVariableCompositeScaleDimensionVariableInline = errors.New("unmarshalling CompositeScaleDimensionVariableInline")
	// ErrUnmarshalVariableFaceCropAlgorithmVariableInline represents an error while unmarshalling FaceCropAlgorithmVariableInline
	ErrUnmarshalVariableFaceCropAlgorithmVariableInline = errors.New("unmarshalling FaceCropAlgorithmVariableInline")
	// ErrUnmarshalVariableFaceCropFocusVariableInline represents an error while unmarshalling FaceCropFocusVariableInline
	ErrUnmarshalVariableFaceCropFocusVariableInline = errors.New("unmarshalling FaceCropFocusVariableInline")
	// ErrUnmarshalVariableFaceCropStyleVariableInline represents an error while unmarshalling FaceCropStyleVariableInline
	ErrUnmarshalVariableFaceCropStyleVariableInline = errors.New("unmarshalling FaceCropStyleVariableInline")
	// ErrUnmarshalVariableFeatureCropStyleVariableInline represents an error while unmarshalling FeatureCropStyleVariableInline
	ErrUnmarshalVariableFeatureCropStyleVariableInline = errors.New("unmarshalling FeatureCropStyleVariableInline")
	// ErrUnmarshalVariableGravityPostVariableInline represents an error while unmarshalling GravityPostVariableInline
	ErrUnmarshalVariableGravityPostVariableInline = errors.New("unmarshalling GravityPostVariableInline")
	// ErrUnmarshalVariableGravityVariableInline represents an error while unmarshalling GravityVariableInline
	ErrUnmarshalVariableGravityVariableInline = errors.New("unmarshalling GravityVariableInline")
	// ErrUnmarshalVariableGrayscaleTypeVariableInline represents an error while unmarshalling GrayscaleTypeVariableInline
	ErrUnmarshalVariableGrayscaleTypeVariableInline = errors.New("unmarshalling GrayscaleTypeVariableInline")
	// ErrUnmarshalVariableIfDimensionDimensionVariableInline represents an error while unmarshalling IfDimensionDimensionVariableInline
	ErrUnmarshalVariableIfDimensionDimensionVariableInline = errors.New("unmarshalling IfDimensionDimensionVariableInline")
	// ErrUnmarshalVariableIfDimensionPostDimensionVariableInline represents an error while unmarshalling IfDimensionPostDimensionVariableInline
	ErrUnmarshalVariableIfDimensionPostDimensionVariableInline = errors.New("unmarshalling IfDimensionPostDimensionVariableInline")
	// ErrUnmarshalVariableIntegerVariableInline represents an error while unmarshalling IntegerVariableInline
	ErrUnmarshalVariableIntegerVariableInline = errors.New("unmarshalling IntegerVariableInline")
	// ErrUnmarshalVariableNumberVariableInline represents an error while unmarshalling NumberVariableInline
	ErrUnmarshalVariableNumberVariableInline = errors.New("unmarshalling NumberVariableInline")
	// ErrUnmarshalVariableOutputImagePerceptualQualityVariableInline represents an error while unmarshalling OutputImagePerceptualQualityVariableInline
	ErrUnmarshalVariableOutputImagePerceptualQualityVariableInline = errors.New("unmarshalling OutputImagePerceptualQualityVariableInline")
	// ErrUnmarshalVariableQueryVariableInline represents an error while unmarshalling QueryVariableInline
	ErrUnmarshalVariableQueryVariableInline = errors.New("unmarshalling QueryVariableInline")
	// ErrUnmarshalVariableRegionOfInterestCropStyleVariableInline represents an error while unmarshalling RegionOfInterestCropStyleVariableInline
	ErrUnmarshalVariableRegionOfInterestCropStyleVariableInline = errors.New("unmarshalling RegionOfInterestCropStyleVariableInline")
	// ErrUnmarshalVariableResizeAspectVariableInline represents an error while unmarshalling ResizeAspectVariableInline
	ErrUnmarshalVariableResizeAspectVariableInline = errors.New("unmarshalling ResizeAspectVariableInline")
	// ErrUnmarshalVariableResizeTypeVariableInline represents an error while unmarshalling ResizeTypeVariableInline
	ErrUnmarshalVariableResizeTypeVariableInline = errors.New("unmarshalling ResizeTypeVariableInline")
	// ErrUnmarshalVariableStringVariableInline represents an error while unmarshalling StringVariableInline
	ErrUnmarshalVariableStringVariableInline = errors.New("unmarshalling StringVariableInline")
	// ErrUnmarshalVariableOutputVideoPerceptualQualityVariableInline represents an error while unmarshalling OutputVideoPerceptualQualityVariableInline
	ErrUnmarshalVariableOutputVideoPerceptualQualityVariableInline = errors.New("unmarshalling OutputVideoPerceptualQualityVariableInline")
	// ErrUnmarshalVariableOutputVideoVideoAdaptiveQualityVariableInline represents an error while unmarshalling OutputVideoVideoAdaptiveQualityVariableInline
	ErrUnmarshalVariableOutputVideoVideoAdaptiveQualityVariableInline = errors.New("unmarshalling OutputVideoVideoAdaptiveQualityVariableInline")
)

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (a *AppendGravityPriorityVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		a.Name = &variable.Var
		a.Value = nil
		return nil
	}
	var value AppendGravityPriority
	if err = json.Unmarshal(in, &value); err == nil {
		a.Name = nil
		a.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableAppendGravityPriorityVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (a *AppendGravityPriorityVariableInline) MarshalJSON() ([]byte, error) {
	if a.Value != nil {
		return json.Marshal(*a.Value)
	}
	if a.Name != nil {
		return json.Marshal(VariableInline{Var: *a.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (b *BooleanVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		b.Name = &variable.Var
		b.Value = nil
		return nil
	}
	var value bool
	if err = json.Unmarshal(in, &value); err == nil {
		b.Name = nil
		b.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableBooleanVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (b *BooleanVariableInline) MarshalJSON() ([]byte, error) {
	if b.Value != nil {
		return json.Marshal(*b.Value)
	}
	if b.Name != nil {
		return json.Marshal(VariableInline{Var: *b.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (c *CompositePlacementVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		c.Name = &variable.Var
		c.Value = nil
		return nil
	}
	var value CompositePlacement
	if err = json.Unmarshal(in, &value); err == nil {
		c.Name = nil
		c.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableCompositePlacementVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (c *CompositePlacementVariableInline) MarshalJSON() ([]byte, error) {
	if c.Value != nil {
		return json.Marshal(*c.Value)
	}
	if c.Name != nil {
		return json.Marshal(VariableInline{Var: *c.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (c *CompositePostPlacementVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		c.Name = &variable.Var
		c.Value = nil
		return nil
	}
	var value CompositePostPlacement
	if err = json.Unmarshal(in, &value); err == nil {
		c.Name = nil
		c.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableCompositePostPlacementVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (c *CompositePostPlacementVariableInline) MarshalJSON() ([]byte, error) {
	if c.Value != nil {
		return json.Marshal(*c.Value)
	}
	if c.Name != nil {
		return json.Marshal(VariableInline{Var: *c.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (c *CompositePostScaleDimensionVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		c.Name = &variable.Var
		c.Value = nil
		return nil
	}
	var value CompositePostScaleDimension
	if err = json.Unmarshal(in, &value); err == nil {
		c.Name = nil
		c.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableCompositePostScaleDimensionVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (c *CompositePostScaleDimensionVariableInline) MarshalJSON() ([]byte, error) {
	if c.Value != nil {
		return json.Marshal(*c.Value)
	}
	if c.Name != nil {
		return json.Marshal(VariableInline{Var: *c.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (c *CompositeScaleDimensionVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		c.Name = &variable.Var
		c.Value = nil
		return nil
	}
	var value CompositeScaleDimension
	if err = json.Unmarshal(in, &value); err == nil {
		c.Name = nil
		c.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableCompositeScaleDimensionVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (c *CompositeScaleDimensionVariableInline) MarshalJSON() ([]byte, error) {
	if c.Value != nil {
		return json.Marshal(*c.Value)
	}
	if c.Name != nil {
		return json.Marshal(VariableInline{Var: *c.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (f *FaceCropAlgorithmVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		f.Name = &variable.Var
		f.Value = nil
		return nil
	}
	var value FaceCropAlgorithm
	if err = json.Unmarshal(in, &value); err == nil {
		f.Name = nil
		f.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableFaceCropAlgorithmVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (f *FaceCropAlgorithmVariableInline) MarshalJSON() ([]byte, error) {
	if f.Value != nil {
		return json.Marshal(*f.Value)
	}
	if f.Name != nil {
		return json.Marshal(VariableInline{Var: *f.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (f *FaceCropFocusVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		f.Name = &variable.Var
		f.Value = nil
		return nil
	}
	var value FaceCropFocus
	if err = json.Unmarshal(in, &value); err == nil {
		f.Name = nil
		f.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableFaceCropFocusVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (f *FaceCropFocusVariableInline) MarshalJSON() ([]byte, error) {
	if f.Value != nil {
		return json.Marshal(*f.Value)
	}
	if f.Name != nil {
		return json.Marshal(VariableInline{Var: *f.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (f *FaceCropStyleVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		f.Name = &variable.Var
		f.Value = nil
		return nil
	}
	var value FaceCropStyle
	if err = json.Unmarshal(in, &value); err == nil {
		f.Name = nil
		f.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableFaceCropStyleVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (f *FaceCropStyleVariableInline) MarshalJSON() ([]byte, error) {
	if f.Value != nil {
		return json.Marshal(*f.Value)
	}
	if f.Name != nil {
		return json.Marshal(VariableInline{Var: *f.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (f *FeatureCropStyleVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		f.Name = &variable.Var
		f.Value = nil
		return nil
	}
	var value FeatureCropStyle
	if err = json.Unmarshal(in, &value); err == nil {
		f.Name = nil
		f.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableFeatureCropStyleVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (f *FeatureCropStyleVariableInline) MarshalJSON() ([]byte, error) {
	if f.Value != nil {
		return json.Marshal(*f.Value)
	}
	if f.Name != nil {
		return json.Marshal(VariableInline{Var: *f.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (g *GravityPostVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		g.Name = &variable.Var
		g.Value = nil
		return nil
	}
	var value GravityPost
	if err = json.Unmarshal(in, &value); err == nil {
		g.Name = nil
		g.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableGravityPostVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (g *GravityPostVariableInline) MarshalJSON() ([]byte, error) {
	if g.Value != nil {
		return json.Marshal(*g.Value)
	}
	if g.Name != nil {
		return json.Marshal(VariableInline{Var: *g.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (g *GravityVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		g.Name = &variable.Var
		g.Value = nil
		return nil
	}
	var value Gravity
	if err = json.Unmarshal(in, &value); err == nil {
		g.Name = nil
		g.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableGravityVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (g *GravityVariableInline) MarshalJSON() ([]byte, error) {
	if g.Value != nil {
		return json.Marshal(*g.Value)
	}
	if g.Name != nil {
		return json.Marshal(VariableInline{Var: *g.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (g *GrayscaleTypeVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		g.Name = &variable.Var
		g.Value = nil
		return nil
	}
	var value GrayscaleType
	if err = json.Unmarshal(in, &value); err == nil {
		g.Name = nil
		g.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableGrayscaleTypeVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (g *GrayscaleTypeVariableInline) MarshalJSON() ([]byte, error) {
	if g.Value != nil {
		return json.Marshal(*g.Value)
	}
	if g.Name != nil {
		return json.Marshal(VariableInline{Var: *g.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (i *IfDimensionDimensionVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		i.Name = &variable.Var
		i.Value = nil
		return nil
	}
	var value IfDimensionDimension
	if err = json.Unmarshal(in, &value); err == nil {
		i.Name = nil
		i.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableIfDimensionDimensionVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (i *IfDimensionDimensionVariableInline) MarshalJSON() ([]byte, error) {
	if i.Value != nil {
		return json.Marshal(*i.Value)
	}
	if i.Name != nil {
		return json.Marshal(VariableInline{Var: *i.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (i *IfDimensionPostDimensionVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		i.Name = &variable.Var
		i.Value = nil
		return nil
	}
	var value IfDimensionPostDimension
	if err = json.Unmarshal(in, &value); err == nil {
		i.Name = nil
		i.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableIfDimensionPostDimensionVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (i *IfDimensionPostDimensionVariableInline) MarshalJSON() ([]byte, error) {
	if i.Value != nil {
		return json.Marshal(*i.Value)
	}
	if i.Name != nil {
		return json.Marshal(VariableInline{Var: *i.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (i *IntegerVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		i.Name = &variable.Var
		i.Value = nil
		return nil
	}
	var value int
	if err = json.Unmarshal(in, &value); err == nil {
		i.Name = nil
		i.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableIntegerVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (i *IntegerVariableInline) MarshalJSON() ([]byte, error) {
	if i.Value != nil {
		return json.Marshal(*i.Value)
	}
	if i.Name != nil {
		return json.Marshal(VariableInline{Var: *i.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (n *NumberVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		n.Name = &variable.Var
		n.Value = nil
		return nil
	}
	var value float64
	if err = json.Unmarshal(in, &value); err == nil {
		n.Name = nil
		n.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableNumberVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (n *NumberVariableInline) MarshalJSON() ([]byte, error) {
	if n.Value != nil {
		return json.Marshal(*n.Value)
	}
	if n.Name != nil {
		return json.Marshal(VariableInline{Var: *n.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (o *OutputImagePerceptualQualityVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		o.Name = &variable.Var
		o.Value = nil
		return nil
	}
	var value OutputImagePerceptualQuality
	if err = json.Unmarshal(in, &value); err == nil {
		o.Name = nil
		o.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableOutputImagePerceptualQualityVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (o *OutputImagePerceptualQualityVariableInline) MarshalJSON() ([]byte, error) {
	if o.Value != nil {
		return json.Marshal(*o.Value)
	}
	if o.Name != nil {
		return json.Marshal(VariableInline{Var: *o.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (q *QueryVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		q.Name = &variable.Var
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableQueryVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (q *QueryVariableInline) MarshalJSON() ([]byte, error) {
	if q.Name != nil {
		return json.Marshal(VariableInline{Var: *q.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (r *RegionOfInterestCropStyleVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		r.Name = &variable.Var
		r.Value = nil
		return nil
	}
	var value RegionOfInterestCropStyle
	if err = json.Unmarshal(in, &value); err == nil {
		r.Name = nil
		r.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableRegionOfInterestCropStyleVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (r *RegionOfInterestCropStyleVariableInline) MarshalJSON() ([]byte, error) {
	if r.Value != nil {
		return json.Marshal(*r.Value)
	}
	if r.Name != nil {
		return json.Marshal(VariableInline{Var: *r.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (r *ResizeAspectVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		r.Name = &variable.Var
		r.Value = nil
		return nil
	}
	var value ResizeAspect
	if err = json.Unmarshal(in, &value); err == nil {
		r.Name = nil
		r.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableResizeAspectVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (r *ResizeAspectVariableInline) MarshalJSON() ([]byte, error) {
	if r.Value != nil {
		return json.Marshal(*r.Value)
	}
	if r.Name != nil {
		return json.Marshal(VariableInline{Var: *r.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (r *ResizeTypeVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		r.Name = &variable.Var
		r.Value = nil
		return nil
	}
	var value ResizeType
	if err = json.Unmarshal(in, &value); err == nil {
		r.Name = nil
		r.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableResizeTypeVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (r *ResizeTypeVariableInline) MarshalJSON() ([]byte, error) {
	if r.Value != nil {
		return json.Marshal(*r.Value)
	}
	if r.Name != nil {
		return json.Marshal(VariableInline{Var: *r.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (s *StringVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		s.Name = &variable.Var
		s.Value = nil
		return nil
	}
	var value string
	if err = json.Unmarshal(in, &value); err == nil {
		s.Name = nil
		s.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableStringVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (s *StringVariableInline) MarshalJSON() ([]byte, error) {
	if s.Value != nil {
		return json.Marshal(*s.Value)
	}
	if s.Name != nil {
		return json.Marshal(VariableInline{Var: *s.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (o *OutputVideoPerceptualQualityVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		o.Name = &variable.Var
		o.Value = nil
		return nil
	}
	var value OutputVideoPerceptualQuality
	if err = json.Unmarshal(in, &value); err == nil {
		o.Name = nil
		o.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableOutputVideoPerceptualQualityVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (o *OutputVideoPerceptualQualityVariableInline) MarshalJSON() ([]byte, error) {
	if o.Value != nil {
		return json.Marshal(*o.Value)
	}
	if o.Name != nil {
		return json.Marshal(VariableInline{Var: *o.Name})
	}
	return nil, nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a variable which can be either a value or a variable object
func (o *OutputVideoVideoAdaptiveQualityVariableInline) UnmarshalJSON(in []byte) error {
	var err error
	var variable InlineVariable
	if err = json.Unmarshal(in, &variable); err == nil {
		o.Name = &variable.Var
		o.Value = nil
		return nil
	}
	var value OutputVideoVideoAdaptiveQuality
	if err = json.Unmarshal(in, &value); err == nil {
		o.Name = nil
		o.Value = &value
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnmarshalVariableOutputVideoVideoAdaptiveQualityVariableInline, err)
}

// MarshalJSON is a custom marshaler used to encode a variable which can be either a value or a variable object
func (o *OutputVideoVideoAdaptiveQualityVariableInline) MarshalJSON() ([]byte, error) {
	if o.Value != nil {
		return json.Marshal(*o.Value)
	}
	if o.Name != nil {
		return json.Marshal(VariableInline{Var: *o.Name})
	}
	return nil, nil
}

/*-----------------------------------------------*/
///////////// Image type unmarshalers /////////////
/*-----------------------------------------------*/

// ImageTypeValueHandlers is a map of available image types
var ImageTypeValueHandlers = map[string]func() ImageType{
	"box":    func() ImageType { return &BoxImageType{} },
	"text":   func() ImageType { return &TextImageType{} },
	"url":    func() ImageType { return &URLImageType{} },
	"circle": func() ImageType { return &CircleImageType{} },
}

// ImageTypePostValueHandlers is a map of available image post types
var ImageTypePostValueHandlers = map[string]func() ImageTypePost{
	"box":    func() ImageTypePost { return &BoxImageTypePost{} },
	"text":   func() ImageTypePost { return &TextImageTypePost{} },
	"url":    func() ImageTypePost { return &URLImageTypePost{} },
	"circle": func() ImageTypePost { return &CircleImageTypePost{} },
}

var (

	// ErrUnmarshalImageTypeAppend represents an error while unmarshalling Append
	ErrUnmarshalImageTypeAppend = errors.New("unmarshalling Append")
	// ErrUnmarshalImageTypeComposite represents an error while unmarshalling Composite
	ErrUnmarshalImageTypeComposite = errors.New("unmarshalling Composite")
)

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to ImageType interface
func (a *Append) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type AppendT Append
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalImageTypeAppend, err)
	}
	image, ok := data["image"]
	if !ok {
		var target AppendT
		err = json.Unmarshal(in, &target)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshalImageTypeAppend, err)
		}
		*a = Append(target)
		return nil
	}
	imageMap := image.(map[string]interface{})
	imageType, ok := imageMap["type"]
	if !ok {
		_, ok := imageMap["url"]
		if !ok {
			return fmt.Errorf("%w: missing image type", ErrUnmarshalImageTypeAppend)
		}
		imageType = "URL"
	}
	typeName, ok := imageType.(string)
	if !ok {
		return fmt.Errorf("%w: 'type' field on image should be a string", ErrUnmarshalImageTypeAppend)
	}
	var target AppendT
	targetImage, ok := ImageTypeValueHandlers[strings.ToLower(typeName)]
	if !ok {
		return fmt.Errorf("%w: invalid image type: %s", ErrUnmarshalImageTypeAppend, imageType)
	}
	target.Image = targetImage()
	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalImageTypeAppend, err)
	}
	*a = Append(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to ImageType interface
func (c *Composite) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type CompositeT Composite
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalImageTypeComposite, err)
	}
	image, ok := data["image"]
	if !ok {
		var target CompositeT
		err = json.Unmarshal(in, &target)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshalImageTypeComposite, err)
		}
		*c = Composite(target)
		return nil
	}
	imageMap := image.(map[string]interface{})
	imageType, ok := imageMap["type"]
	if !ok {
		_, ok := imageMap["url"]
		if !ok {
			return fmt.Errorf("%w: missing image type", ErrUnmarshalImageTypeComposite)
		}
		imageType = "URL"
	}
	typeName, ok := imageType.(string)
	if !ok {
		return fmt.Errorf("%w: 'type' field on image should be a string", ErrUnmarshalImageTypeComposite)
	}
	var target CompositeT
	targetImage, ok := ImageTypeValueHandlers[strings.ToLower(typeName)]
	if !ok {
		return fmt.Errorf("%w: invalid image type: %s", ErrUnmarshalImageTypeComposite, imageType)
	}
	target.Image = targetImage()
	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalImageTypeComposite, err)
	}
	*c = Composite(target)
	return nil
}

var (

	// ErrUnmarshalImageTypeCompositePost represents an error while unmarshalling CompositePost
	ErrUnmarshalImageTypeCompositePost = errors.New("unmarshalling CompositePost")
)

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to ImageType interface
func (c *CompositePost) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type CompositePostT CompositePost
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalImageTypeCompositePost, err)
	}
	image, ok := data["image"]
	if !ok {
		var target CompositePostT
		err = json.Unmarshal(in, &target)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshalImageTypeCompositePost, err)
		}
		*c = CompositePost(target)
		return nil
	}
	imageMap := image.(map[string]interface{})
	imageType, ok := imageMap["type"]
	if !ok {
		_, ok := imageMap["url"]
		if !ok {
			return fmt.Errorf("%w: missing image type", ErrUnmarshalImageTypeCompositePost)
		}
		imageType = "URL"
	}
	typeName, ok := imageType.(string)
	if !ok {
		return fmt.Errorf("%w: 'type' field on image should be a string", ErrUnmarshalImageTypeCompositePost)
	}
	var target CompositePostT
	targetImage, ok := ImageTypePostValueHandlers[strings.ToLower(typeName)]
	if !ok {
		return fmt.Errorf("%w: invalid image type: %s", ErrUnmarshalImageTypeCompositePost, imageType)
	}
	target.Image = targetImage()
	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalImageTypeCompositePost, err)
	}
	*c = CompositePost(target)
	return nil
}

/*-----------------------------------------------*/
///////////// Shape type unmarshalers /////////////
/*-----------------------------------------------*/

// ShapeTypes is a map of available shape types
var ShapeTypes = map[string]func() ShapeType{
	"circle":    func() ShapeType { return &CircleShapeType{} },
	"point":     func() ShapeType { return &PointShapeType{} },
	"polygon":   func() ShapeType { return &PolygonShapeType{} },
	"rectangle": func() ShapeType { return &RectangleShapeType{} },
	"union":     func() ShapeType { return &UnionShapeType{} },
}

// ShapeTypeValueHandlers returns a ShapeType based on fields specific for a concrete ShapeType
var ShapeTypeValueHandlers = func(m map[string]interface{}) ShapeType {
	if _, ok := m["radius"]; ok {
		return ShapeTypes["circle"]()
	}
	if _, ok := m["x"]; ok {
		return ShapeTypes["point"]()
	}
	if _, ok := m["points"]; ok {
		return ShapeTypes["polygon"]()
	}
	if _, ok := m["anchor"]; ok {
		return ShapeTypes["rectangle"]()
	}
	if _, ok := m["shapes"]; ok {
		return ShapeTypes["union"]()
	}
	return nil
}

var (

	// ErrUnmarshalShapeTypeRegionOfInterestCrop represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalShapeTypeRegionOfInterestCrop = errors.New("unmarshalling RegionOfInterestCrop")
)

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to ShapeType interface
func (r *RegionOfInterestCrop) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type RegionOfInterestCropT RegionOfInterestCrop
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalShapeTypeRegionOfInterestCrop, err)
	}
	shape, ok := data["regionOfInterest"]
	if !ok {
		var target RegionOfInterestCropT
		err = json.Unmarshal(in, &target)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshalShapeTypeRegionOfInterestCrop, err)
		}
		*r = RegionOfInterestCrop(target)
		return nil
	}
	shapeMap := shape.(map[string]interface{})
	var target RegionOfInterestCropT
	targetShape := ShapeTypeValueHandlers(shapeMap)
	if targetShape == nil {
		return fmt.Errorf("%w: invalid shape type", ErrUnmarshalShapeTypeRegionOfInterestCrop)
	}
	target.RegionOfInterest = targetShape
	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalShapeTypeRegionOfInterestCrop, err)
	}
	*r = RegionOfInterestCrop(target)
	return nil
}

/*-----------------------------------------------*/
////////// Transformation unmarshallers ///////////
/*-----------------------------------------------*/

var (

	// ErrUnmarshalTransformationBoxImageType represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalTransformationBoxImageType = errors.New("unmarshalling BoxImageType")
	// ErrUnmarshalTransformationCircleImageType represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalTransformationCircleImageType = errors.New("unmarshalling CircleImageType")
	// ErrUnmarshalTransformationFitAndFill represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalTransformationFitAndFill = errors.New("unmarshalling FitAndFill")
	// ErrUnmarshalTransformationIfDimension represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalTransformationIfDimension = errors.New("unmarshalling IfDimension")
	// ErrUnmarshalTransformationIfOrientation represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalTransformationIfOrientation = errors.New("unmarshalling IfOrientation")
	// ErrUnmarshalTransformationTextImageType represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalTransformationTextImageType = errors.New("unmarshalling TextImageType")
	// ErrUnmarshalTransformationURLImageType represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalTransformationURLImageType = errors.New("unmarshalling URLImageType")
)

var (

	// ErrUnmarshalPostBreakpointTransformationBoxImageTypePost represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalPostBreakpointTransformationBoxImageTypePost = errors.New("unmarshalling BoxImageTypePost")
	// ErrUnmarshalPostBreakpointTransformationCircleImageTypePost represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalPostBreakpointTransformationCircleImageTypePost = errors.New("unmarshalling CircleImageTypePost")
	// ErrUnmarshalPostBreakpointTransformationIfDimensionPost represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalPostBreakpointTransformationIfDimensionPost = errors.New("unmarshalling IfDimensionPost")
	// ErrUnmarshalPostBreakpointTransformationIfOrientationPost represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalPostBreakpointTransformationIfOrientationPost = errors.New("unmarshalling IfOrientationPost")
	// ErrUnmarshalPostBreakpointTransformationTextImageTypePost represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalPostBreakpointTransformationTextImageTypePost = errors.New("unmarshalling TextImageTypePost")
	// ErrUnmarshalPostBreakpointTransformationURLImageTypePost represents an error while unmarshalling {$compositeType}}
	ErrUnmarshalPostBreakpointTransformationURLImageTypePost = errors.New("unmarshalling URLImageTypePost")
)

// TransformationHandlers is a map of available transformations
var TransformationHandlers = map[string]func() TransformationType{
	"Append":               func() TransformationType { return &Append{} },
	"AspectCrop":           func() TransformationType { return &AspectCrop{} },
	"BackgroundColor":      func() TransformationType { return &BackgroundColor{} },
	"Blur":                 func() TransformationType { return &Blur{} },
	"ChromaKey":            func() TransformationType { return &ChromaKey{} },
	"Composite":            func() TransformationType { return &Composite{} },
	"Compound":             func() TransformationType { return &Compound{} },
	"Contrast":             func() TransformationType { return &Contrast{} },
	"Crop":                 func() TransformationType { return &Crop{} },
	"FaceCrop":             func() TransformationType { return &FaceCrop{} },
	"FeatureCrop":          func() TransformationType { return &FeatureCrop{} },
	"FitAndFill":           func() TransformationType { return &FitAndFill{} },
	"Goop":                 func() TransformationType { return &Goop{} },
	"Grayscale":            func() TransformationType { return &Grayscale{} },
	"HSL":                  func() TransformationType { return &HSL{} },
	"HSV":                  func() TransformationType { return &HSV{} },
	"IfDimension":          func() TransformationType { return &IfDimension{} },
	"IfOrientation":        func() TransformationType { return &IfOrientation{} },
	"ImQuery":              func() TransformationType { return &ImQuery{} },
	"MaxColors":            func() TransformationType { return &MaxColors{} },
	"Mirror":               func() TransformationType { return &Mirror{} },
	"MonoHue":              func() TransformationType { return &MonoHue{} },
	"Opacity":              func() TransformationType { return &Opacity{} },
	"RegionOfInterestCrop": func() TransformationType { return &RegionOfInterestCrop{} },
	"RelativeCrop":         func() TransformationType { return &RelativeCrop{} },
	"RemoveColor":          func() TransformationType { return &RemoveColor{} },
	"Resize":               func() TransformationType { return &Resize{} },
	"Rotate":               func() TransformationType { return &Rotate{} },
	"Scale":                func() TransformationType { return &Scale{} },
	"Shear":                func() TransformationType { return &Shear{} },
	"Trim":                 func() TransformationType { return &Trim{} },
	"UnsharpMask":          func() TransformationType { return &UnsharpMask{} },
}

// PostBreakpointTransformationHandlers is a map of available PostBreakpointTransformations
var PostBreakpointTransformationHandlers = map[string]func() TransformationTypePost{
	"BackgroundColor": func() TransformationTypePost { return &BackgroundColor{} },
	"Blur":            func() TransformationTypePost { return &Blur{} },
	"ChromaKey":       func() TransformationTypePost { return &ChromaKey{} },
	"Compound":        func() TransformationTypePost { return &CompoundPost{} },
	"Composite":       func() TransformationTypePost { return &CompositePost{} },
	"Contrast":        func() TransformationTypePost { return &Contrast{} },
	"Goop":            func() TransformationTypePost { return &Goop{} },
	"Grayscale":       func() TransformationTypePost { return &Grayscale{} },
	"HSL":             func() TransformationTypePost { return &HSL{} },
	"HSV":             func() TransformationTypePost { return &HSV{} },
	"IfDimension":     func() TransformationTypePost { return &IfDimensionPost{} },
	"IfOrientation":   func() TransformationTypePost { return &IfOrientationPost{} },
	"MaxColors":       func() TransformationTypePost { return &MaxColors{} },
	"Mirror":          func() TransformationTypePost { return &Mirror{} },
	"MonoHue":         func() TransformationTypePost { return &MonoHue{} },
	"Opacity":         func() TransformationTypePost { return &Opacity{} },
	"RemoveColor":     func() TransformationTypePost { return &RemoveColor{} },
	"UnsharpMask":     func() TransformationTypePost { return &UnsharpMask{} },
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to Transformation interface
func (b *BoxImageType) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type BoxImageTypeT BoxImageType
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationBoxImageType, err)
	}
	var target BoxImageTypeT

	transformationParam, ok := data["transformation"]
	if ok {
		transformationMap, ok := transformationParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on BoxImageType should be a map", ErrUnmarshalTransformationBoxImageType)
		}
		typeName := transformationMap["transformation"].(string)
		transformationTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationBoxImageType, typeName)
		}
		target.Transformation = transformationTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationBoxImageType, err)
	}
	*b = BoxImageType(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to Transformation interface
func (c *CircleImageType) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type CircleImageTypeT CircleImageType
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationCircleImageType, err)
	}
	var target CircleImageTypeT

	transformationParam, ok := data["transformation"]
	if ok {
		transformationMap, ok := transformationParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on CircleImageType should be a map", ErrUnmarshalTransformationCircleImageType)
		}
		typeName := transformationMap["transformation"].(string)
		transformationTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationCircleImageType, typeName)
		}
		target.Transformation = transformationTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationCircleImageType, err)
	}
	*c = CircleImageType(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to Transformation interface
func (f *FitAndFill) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type FitAndFillT FitAndFill
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationFitAndFill, err)
	}
	var target FitAndFillT

	fillTransformationParam, ok := data["fillTransformation"]
	if ok {
		fillTransformationMap, ok := fillTransformationParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'fillTransformation' field on FitAndFill should be a map", ErrUnmarshalTransformationFitAndFill)
		}
		typeName := fillTransformationMap["transformation"].(string)
		fillTransformationTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationFitAndFill, typeName)
		}
		target.FillTransformation = fillTransformationTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationFitAndFill, err)
	}
	*f = FitAndFill(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to Transformation interface
func (i *IfDimension) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type IfDimensionT IfDimension
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationIfDimension, err)
	}
	var target IfDimensionT

	defaultParam, ok := data["default"]
	if ok {
		defaultMap, ok := defaultParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'default' field on IfDimension should be a map", ErrUnmarshalTransformationIfDimension)
		}
		typeName := defaultMap["transformation"].(string)
		defaultTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationIfDimension, typeName)
		}
		target.Default = defaultTarget()
	}

	equalParam, ok := data["equal"]
	if ok {
		equalMap, ok := equalParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'equal' field on IfDimension should be a map", ErrUnmarshalTransformationIfDimension)
		}
		typeName := equalMap["transformation"].(string)
		equalTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationIfDimension, typeName)
		}
		target.Equal = equalTarget()
	}

	greaterThanParam, ok := data["greaterThan"]
	if ok {
		greaterThanMap, ok := greaterThanParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'greaterThan' field on IfDimension should be a map", ErrUnmarshalTransformationIfDimension)
		}
		typeName := greaterThanMap["transformation"].(string)
		greaterThanTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationIfDimension, typeName)
		}
		target.GreaterThan = greaterThanTarget()
	}

	lessThanParam, ok := data["lessThan"]
	if ok {
		lessThanMap, ok := lessThanParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'lessThan' field on IfDimension should be a map", ErrUnmarshalTransformationIfDimension)
		}
		typeName := lessThanMap["transformation"].(string)
		lessThanTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationIfDimension, typeName)
		}
		target.LessThan = lessThanTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationIfDimension, err)
	}
	*i = IfDimension(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to Transformation interface
func (i *IfOrientation) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type IfOrientationT IfOrientation
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationIfOrientation, err)
	}
	var target IfOrientationT

	defaultParam, ok := data["default"]
	if ok {
		defaultMap, ok := defaultParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'default' field on IfOrientation should be a map", ErrUnmarshalTransformationIfOrientation)
		}
		typeName := defaultMap["transformation"].(string)
		defaultTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationIfOrientation, typeName)
		}
		target.Default = defaultTarget()
	}

	landscapeParam, ok := data["landscape"]
	if ok {
		landscapeMap, ok := landscapeParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'landscape' field on IfOrientation should be a map", ErrUnmarshalTransformationIfOrientation)
		}
		typeName := landscapeMap["transformation"].(string)
		landscapeTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationIfOrientation, typeName)
		}
		target.Landscape = landscapeTarget()
	}

	portraitParam, ok := data["portrait"]
	if ok {
		portraitMap, ok := portraitParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'portrait' field on IfOrientation should be a map", ErrUnmarshalTransformationIfOrientation)
		}
		typeName := portraitMap["transformation"].(string)
		portraitTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationIfOrientation, typeName)
		}
		target.Portrait = portraitTarget()
	}

	squareParam, ok := data["square"]
	if ok {
		squareMap, ok := squareParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'square' field on IfOrientation should be a map", ErrUnmarshalTransformationIfOrientation)
		}
		typeName := squareMap["transformation"].(string)
		squareTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationIfOrientation, typeName)
		}
		target.Square = squareTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationIfOrientation, err)
	}
	*i = IfOrientation(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to Transformation interface
func (t *TextImageType) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type TextImageTypeT TextImageType
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationTextImageType, err)
	}
	var target TextImageTypeT

	transformationParam, ok := data["transformation"]
	if ok {
		transformationMap, ok := transformationParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on TextImageType should be a map", ErrUnmarshalTransformationTextImageType)
		}
		typeName := transformationMap["transformation"].(string)
		transformationTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationTextImageType, typeName)
		}
		target.Transformation = transformationTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationTextImageType, err)
	}
	*t = TextImageType(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to Transformation interface
func (u *URLImageType) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type URLImageTypeT URLImageType
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationURLImageType, err)
	}
	var target URLImageTypeT

	transformationParam, ok := data["transformation"]
	if ok {
		transformationMap, ok := transformationParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on URLImageType should be a map", ErrUnmarshalTransformationURLImageType)
		}
		typeName := transformationMap["transformation"].(string)
		transformationTarget, ok := TransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalTransformationURLImageType, typeName)
		}
		target.Transformation = transformationTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationURLImageType, err)
	}
	*u = URLImageType(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to PostBreakpointTransformation interface
func (b *BoxImageTypePost) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type BoxImageTypePostT BoxImageTypePost
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationBoxImageTypePost, err)
	}
	var target BoxImageTypePostT

	transformationParam, ok := data["transformation"]
	if ok {
		transformationMap, ok := transformationParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on BoxImageTypePost should be a map", ErrUnmarshalPostBreakpointTransformationBoxImageTypePost)
		}
		typeName := transformationMap["transformation"].(string)
		transformationTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationBoxImageTypePost, typeName)
		}
		target.Transformation = transformationTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationBoxImageTypePost, err)
	}
	*b = BoxImageTypePost(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to PostBreakpointTransformation interface
func (c *CircleImageTypePost) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type CircleImageTypePostT CircleImageTypePost
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationCircleImageTypePost, err)
	}
	var target CircleImageTypePostT

	transformationParam, ok := data["transformation"]
	if ok {
		transformationMap, ok := transformationParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on CircleImageTypePost should be a map", ErrUnmarshalPostBreakpointTransformationCircleImageTypePost)
		}
		typeName := transformationMap["transformation"].(string)
		transformationTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationCircleImageTypePost, typeName)
		}
		target.Transformation = transformationTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationCircleImageTypePost, err)
	}
	*c = CircleImageTypePost(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to PostBreakpointTransformation interface
func (i *IfDimensionPost) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type IfDimensionPostT IfDimensionPost
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationIfDimensionPost, err)
	}
	var target IfDimensionPostT

	defaultParam, ok := data["default"]
	if ok {
		defaultMap, ok := defaultParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'default' field on IfDimensionPost should be a map", ErrUnmarshalPostBreakpointTransformationIfDimensionPost)
		}
		typeName := defaultMap["transformation"].(string)
		defaultTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationIfDimensionPost, typeName)
		}
		target.Default = defaultTarget()
	}

	equalParam, ok := data["equal"]
	if ok {
		equalMap, ok := equalParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'equal' field on IfDimensionPost should be a map", ErrUnmarshalPostBreakpointTransformationIfDimensionPost)
		}
		typeName := equalMap["transformation"].(string)
		equalTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationIfDimensionPost, typeName)
		}
		target.Equal = equalTarget()
	}

	greaterThanParam, ok := data["greaterThan"]
	if ok {
		greaterThanMap, ok := greaterThanParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'greaterThan' field on IfDimensionPost should be a map", ErrUnmarshalPostBreakpointTransformationIfDimensionPost)
		}
		typeName := greaterThanMap["transformation"].(string)
		greaterThanTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationIfDimensionPost, typeName)
		}
		target.GreaterThan = greaterThanTarget()
	}

	lessThanParam, ok := data["lessThan"]
	if ok {
		lessThanMap, ok := lessThanParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'lessThan' field on IfDimensionPost should be a map", ErrUnmarshalPostBreakpointTransformationIfDimensionPost)
		}
		typeName := lessThanMap["transformation"].(string)
		lessThanTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationIfDimensionPost, typeName)
		}
		target.LessThan = lessThanTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationIfDimensionPost, err)
	}
	*i = IfDimensionPost(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to PostBreakpointTransformation interface
func (i *IfOrientationPost) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type IfOrientationPostT IfOrientationPost
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationIfOrientationPost, err)
	}
	var target IfOrientationPostT

	defaultParam, ok := data["default"]
	if ok {
		defaultMap, ok := defaultParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'default' field on IfOrientationPost should be a map", ErrUnmarshalPostBreakpointTransformationIfOrientationPost)
		}
		typeName := defaultMap["transformation"].(string)
		defaultTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationIfOrientationPost, typeName)
		}
		target.Default = defaultTarget()
	}

	landscapeParam, ok := data["landscape"]
	if ok {
		landscapeMap, ok := landscapeParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'landscape' field on IfOrientationPost should be a map", ErrUnmarshalPostBreakpointTransformationIfOrientationPost)
		}
		typeName := landscapeMap["transformation"].(string)
		landscapeTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationIfOrientationPost, typeName)
		}
		target.Landscape = landscapeTarget()
	}

	portraitParam, ok := data["portrait"]
	if ok {
		portraitMap, ok := portraitParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'portrait' field on IfOrientationPost should be a map", ErrUnmarshalPostBreakpointTransformationIfOrientationPost)
		}
		typeName := portraitMap["transformation"].(string)
		portraitTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationIfOrientationPost, typeName)
		}
		target.Portrait = portraitTarget()
	}

	squareParam, ok := data["square"]
	if ok {
		squareMap, ok := squareParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'square' field on IfOrientationPost should be a map", ErrUnmarshalPostBreakpointTransformationIfOrientationPost)
		}
		typeName := squareMap["transformation"].(string)
		squareTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationIfOrientationPost, typeName)
		}
		target.Square = squareTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationIfOrientationPost, err)
	}
	*i = IfOrientationPost(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to PostBreakpointTransformation interface
func (t *TextImageTypePost) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type TextImageTypePostT TextImageTypePost
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationTextImageTypePost, err)
	}
	var target TextImageTypePostT

	transformationParam, ok := data["transformation"]
	if ok {
		transformationMap, ok := transformationParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on TextImageTypePost should be a map", ErrUnmarshalPostBreakpointTransformationTextImageTypePost)
		}
		typeName := transformationMap["transformation"].(string)
		transformationTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationTextImageTypePost, typeName)
		}
		target.Transformation = transformationTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationTextImageTypePost, err)
	}
	*t = TextImageTypePost(target)
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a type containing a reference to PostBreakpointTransformation interface
func (u *URLImageTypePost) UnmarshalJSON(in []byte) error {
	data := make(map[string]interface{})
	type URLImageTypePostT URLImageTypePost
	err := json.Unmarshal(in, &data)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationURLImageTypePost, err)
	}
	var target URLImageTypePostT

	transformationParam, ok := data["transformation"]
	if ok {
		transformationMap, ok := transformationParam.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on URLImageTypePost should be a map", ErrUnmarshalPostBreakpointTransformationURLImageTypePost)
		}
		typeName := transformationMap["transformation"].(string)
		transformationTarget, ok := PostBreakpointTransformationHandlers[typeName]
		if !ok {
			return fmt.Errorf("%w: invalid transformation type: %s", ErrUnmarshalPostBreakpointTransformationURLImageTypePost, typeName)
		}
		target.Transformation = transformationTarget()
	}

	err = json.Unmarshal(in, &target)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationURLImageTypePost, err)
	}
	*u = URLImageTypePost(target)
	return nil
}

// ErrUnmarshalTransformationList represents an error while unmarshalling transformation list
var ErrUnmarshalTransformationList = errors.New("unmarshalling transformation list")

// ErrUnmarshalPostBreakpointTransformationList represents an error while unmarshalling post breakpoint transformation list
var ErrUnmarshalPostBreakpointTransformationList = errors.New("unmarshalling post breakpoint transformation list")

// UnmarshalJSON is a custom unmarshaler used to decode a slice of Transformation interfaces
func (t *Transformations) UnmarshalJSON(in []byte) error {
	data := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(in, &data); err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalTransformationList, err)
	}
	for _, transformation := range data {
		transformationType, ok := transformation["transformation"]
		if !ok {
			return fmt.Errorf("%w: transformation should contain 'transformation' field", ErrUnmarshalTransformationList)
		}
		transformationTypeName, ok := transformationType.(string)
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on transformation entry should be a string", ErrUnmarshalTransformationList)
		}

		bytes, err := json.Marshal(transformation)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshalTransformationList, err)
		}

		indicatedTransformationType, ok := TransformationHandlers[transformationTypeName]
		if !ok {
			return fmt.Errorf("%w: unsupported transformation type: %s", ErrUnmarshalTransformationList, transformationTypeName)
		}
		ipt := indicatedTransformationType()
		err = json.Unmarshal(bytes, ipt)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshalTransformationList, err)
		}
		*t = append(*t, ipt)
	}
	return nil
}

// UnmarshalJSON is a custom unmarshaler used to decode a slice of PostBreakpointTransformation interfaces
func (t *PostBreakpointTransformations) UnmarshalJSON(in []byte) error {
	data := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(in, &data); err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationList, err)
	}
	for _, transformation := range data {
		transformationType, ok := transformation["transformation"]
		if !ok {
			return fmt.Errorf("%w: transformation should contain 'transformation' field", ErrUnmarshalPostBreakpointTransformationList)
		}
		transformationTypeName, ok := transformationType.(string)
		if !ok {
			return fmt.Errorf("%w: 'transformation' field on transformation entry should be a string", ErrUnmarshalPostBreakpointTransformationList)
		}

		bytes, err := json.Marshal(transformation)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationList, err)
		}

		indicatedTransformationType, ok := PostBreakpointTransformationHandlers[transformationTypeName]
		if !ok {
			return fmt.Errorf("%w: unsupported transformation type: %s", ErrUnmarshalPostBreakpointTransformationList, transformationTypeName)
		}
		ipt := indicatedTransformationType()
		err = json.Unmarshal(bytes, ipt)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrUnmarshalPostBreakpointTransformationList, err)
		}
		*t = append(*t, ipt)
	}
	return nil
}
