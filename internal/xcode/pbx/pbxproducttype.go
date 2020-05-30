package pbx

type PBXProductType = string

const (
	AppExtension        PBXProductType = "com.apple.product-type.app-extension"
	Application         PBXProductType = "com.apple.product-type.application"
	Bundle              PBXProductType = "com.apple.product-type.bundle"
	CommandLineTool     PBXProductType = "com.apple.product-type.tool"
	DynamicLibrary      PBXProductType = "com.apple.product-type.library.dynamic"
	Framework           PBXProductType = "com.apple.product-type.framework"
	MessagesApplication PBXProductType = "com.apple.product-type.application.messages"
	MessagesExtension   PBXProductType = "com.apple.product-ype.app-extension.messages"
	OcUnitTestBundle    PBXProductType = "com.apple.product-type.bundle.ocunit-test"
	StaticLibrary       PBXProductType = "com.apple.product-type.library.static"
	StickerPack         PBXProductType = "com.apple.product-type.app-extension.messages-sticker-pack"
	TvExtension         PBXProductType = "com.apple.product-type.tv-app-extension"
	UiTestBundle        PBXProductType = "com.apple.product-type.bundle.ui-testing"
	UnitTestBundle      PBXProductType = "com.apple.product-type.bundle.unit-test"
	Watch2App           PBXProductType = "com.apple.product-type.application.watchapp2"
	Watch2Extension     PBXProductType = "com.apple.product-type.watchkit2-extension"
	WatchApp            PBXProductType = "com.apple.product-type.application.watchapp"
	WatchExtension      PBXProductType = "com.apple.product-type.watchkit-extension"
	XcodeExtension      PBXProductType = "com.apple.product-type.xcode-extension"
	XpcService          PBXProductType = "com.apple.product-type.xpc-service"
)
