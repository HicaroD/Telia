package sdl

use "vendor::sdl::raw"
use "std::strings"

type Environment = raw::Environment
type Void = raw::Void
type iconv_data_t = raw::iconv_data_t
type AssertData = raw::AssertData
type AssertState = raw::AssertState
type AsyncIO = raw::AsyncIO
type AsyncIOQueue = raw::AsyncIOQueue
type AsyncIOTaskType = raw::AsyncIOTaskType
type AsyncIOResult = raw::AsyncIOResult
type AsyncIOOutcome = raw::AsyncIOOutcome
type AtomicInt = raw::AtomicInt
type AtomicU32 = raw::AtomicU32
type PropertyType = raw::PropertyType
type Thread = raw::Thread
type ThreadPriority = raw::ThreadPriority
type ThreadState = raw::ThreadState
type Mutex = raw::Mutex
type RWLock = raw::RWLock
type Semaphore = raw::Semaphore
type Condition = raw::Condition
type InitState = raw::InitState
type IOStream = raw::IOStream
type IOWhence = raw::IOWhence
type IOStatus = raw::IOStatus
type IOStreamInterface = raw::IOStreamInterface
type AudioFormat = raw::AudioFormat
type AudioSpec = raw::AudioSpec
type AudioStream = raw::AudioStream
type BlendFactor = raw::BlendFactor
type BlendOperation = raw::BlendOperation
type PixelFormat = raw::PixelFormat
type PixelFormatDetails = raw::PixelFormatDetails
type Color = raw::Color
type Palette = raw::Palette
type Rect = raw::Rect
type FRect = raw::FRect
type Point = raw::Point
type FPoint = raw::FPoint
type Surface = raw::Surface
type Colorspace = raw::Colorspace
type FlipMode = raw::FlipMode
type ScaleMode = raw::ScaleMode
type CameraSpec = raw::CameraSpec
type CameraPosition = raw::CameraPosition
type Camera = raw::Camera
type SystemTheme = raw::SystemTheme
type DisplayOrientation = raw::DisplayOrientation
type DisplayModeData = raw::DisplayModeData
type DisplayMode = raw::DisplayMode
type Window = raw::Window
type HitTestResult = raw::HitTestResult
type FlashOperation = raw::FlashOperation
type GLAttr = raw::GLAttr
type GLContextState = raw::GLContextState
type DialogFileFilter = raw::DialogFileFilter
type FileDialogType = raw::FileDialogType
type GUID = raw::GUID
type PowerState = raw::PowerState
type SensorType = raw::SensorType
type Sensor = raw::Sensor
type JoystickType = raw::JoystickType
type Joystick = raw::Joystick
type VirtualJoystickTouchpadDesc = raw::VirtualJoystickTouchpadDesc
type VirtualJoystickSensorDesc = raw::VirtualJoystickSensorDesc
type VirtualJoystickDesc = raw::VirtualJoystickDesc
type JoystickConnectionState = raw::JoystickConnectionState
type Gamepad = raw::Gamepad
type GamepadType = raw::GamepadType
type GamepadBindingType = raw::GamepadBindingType
type GamepadBindingInputAxis = raw::GamepadBindingInputAxis
type GamepadBindingInputHat = raw::GamepadBindingInputHat
type GamepadBindingInput = raw::GamepadBindingInput
type GamepadButton = raw::GamepadButton
type GamepadAxis = raw::GamepadAxis
type GamepadBindingOutputAxis = raw::GamepadBindingOutputAxis
type GamepadBindingOutput = raw::GamepadBindingOutput
type GamepadBinding = raw::GamepadBinding
type GamepadButtonLabel = raw::GamepadButtonLabel
type Scancode = raw::Scancode
type Cursor = raw::Cursor
type SystemCursor = raw::SystemCursor
type TouchDeviceType = raw::TouchDeviceType
type Finger = raw::Finger
type CommonEvent = raw::CommonEvent
type EventType = raw::EventType
type DisplayEvent = raw::DisplayEvent
type WindowEvent = raw::WindowEvent
type KeyboardDeviceEvent = raw::KeyboardDeviceEvent
type KeyboardEvent = raw::KeyboardEvent
type TextEditingEvent = raw::TextEditingEvent
type TextEditingCandidatesEvent = raw::TextEditingCandidatesEvent
type TextInputEvent = raw::TextInputEvent
type MouseDeviceEvent = raw::MouseDeviceEvent
type MouseMotionEvent = raw::MouseMotionEvent
type MouseButtonEvent = raw::MouseButtonEvent
type MouseWheelDirection = raw::MouseWheelDirection
type MouseWheelEvent = raw::MouseWheelEvent
type JoyDeviceEvent = raw::JoyDeviceEvent
type JoyAxisEvent = raw::JoyAxisEvent
type JoyBallEvent = raw::JoyBallEvent
type JoyHatEvent = raw::JoyHatEvent
type JoyButtonEvent = raw::JoyButtonEvent
type JoyBatteryEvent = raw::JoyBatteryEvent
type GamepadDeviceEvent = raw::GamepadDeviceEvent
type GamepadAxisEvent = raw::GamepadAxisEvent
type GamepadButtonEvent = raw::GamepadButtonEvent
type GamepadTouchpadEvent = raw::GamepadTouchpadEvent
type GamepadSensorEvent = raw::GamepadSensorEvent
type AudioDeviceEvent = raw::AudioDeviceEvent
type CameraDeviceEvent = raw::CameraDeviceEvent
type SensorEvent = raw::SensorEvent
type QuitEvent = raw::QuitEvent
type UserEvent = raw::UserEvent
type TouchFingerEvent = raw::TouchFingerEvent
type PenProximityEvent = raw::PenProximityEvent
type PenTouchEvent = raw::PenTouchEvent
type PenMotionEvent = raw::PenMotionEvent
type PenButtonEvent = raw::PenButtonEvent
type PenAxis = raw::PenAxis
type PenAxisEvent = raw::PenAxisEvent
type RenderEvent = raw::RenderEvent
type DropEvent = raw::DropEvent
type ClipboardEvent = raw::ClipboardEvent
type Event = raw::Event
type EventAction = raw::EventAction
type Folder = raw::Folder
type EnumerationResult = raw::EnumerationResult
type PathType = raw::PathType
type PathInfo = raw::PathInfo
type GPUDevice = raw::GPUDevice
type GPUComputePipelineCreateInfo = raw::GPUComputePipelineCreateInfo
type GPUComputePipeline = raw::GPUComputePipeline
type GPUShader = raw::GPUShader
type GPUVertexInputRate = raw::GPUVertexInputRate
type GPUVertexBufferDescription = raw::GPUVertexBufferDescription
type GPUVertexElementFormat = raw::GPUVertexElementFormat
type GPUVertexAttribute = raw::GPUVertexAttribute
type GPUVertexInputState = raw::GPUVertexInputState
type GPUPrimitiveType = raw::GPUPrimitiveType
type GPUFillMode = raw::GPUFillMode
type GPUCullMode = raw::GPUCullMode
type GPUFrontFace = raw::GPUFrontFace
type GPURasterizerState = raw::GPURasterizerState
type GPUSampleCount = raw::GPUSampleCount
type GPUMultisampleState = raw::GPUMultisampleState
type GPUCompareOp = raw::GPUCompareOp
type GPUStencilOp = raw::GPUStencilOp
type GPUStencilOpState = raw::GPUStencilOpState
type GPUDepthStencilState = raw::GPUDepthStencilState
type GPUTextureFormat = raw::GPUTextureFormat
type GPUBlendFactor = raw::GPUBlendFactor
type GPUBlendOp = raw::GPUBlendOp
type GPUColorTargetBlendState = raw::GPUColorTargetBlendState
type GPUColorTargetDescription = raw::GPUColorTargetDescription
type GPUGraphicsPipelineTargetInfo = raw::GPUGraphicsPipelineTargetInfo
type GPUGraphicsPipelineCreateInfo = raw::GPUGraphicsPipelineCreateInfo
type GPUGraphicsPipeline = raw::GPUGraphicsPipeline
type GPUFilter = raw::GPUFilter
type GPUSamplerMipmapMode = raw::GPUSamplerMipmapMode
type GPUSamplerAddressMode = raw::GPUSamplerAddressMode
type GPUSamplerCreateInfo = raw::GPUSamplerCreateInfo
type GPUSampler = raw::GPUSampler
type GPUShaderStage = raw::GPUShaderStage
type GPUShaderCreateInfo = raw::GPUShaderCreateInfo
type GPUTextureType = raw::GPUTextureType
type GPUTextureCreateInfo = raw::GPUTextureCreateInfo
type GPUTexture = raw::GPUTexture
type GPUBufferCreateInfo = raw::GPUBufferCreateInfo
type GPUBuffer = raw::GPUBuffer
type GPUTransferBufferUsage = raw::GPUTransferBufferUsage
type GPUTransferBufferCreateInfo = raw::GPUTransferBufferCreateInfo
type GPUTransferBuffer = raw::GPUTransferBuffer
type GPUCommandBuffer = raw::GPUCommandBuffer
type FColor = raw::FColor
type GPULoadOp = raw::GPULoadOp
type GPUStoreOp = raw::GPUStoreOp
type GPUColorTargetInfo = raw::GPUColorTargetInfo
type GPUDepthStencilTargetInfo = raw::GPUDepthStencilTargetInfo
type GPURenderPass = raw::GPURenderPass
type GPUViewport = raw::GPUViewport
type GPUBufferBinding = raw::GPUBufferBinding
type GPUIndexElementSize = raw::GPUIndexElementSize
type GPUTextureSamplerBinding = raw::GPUTextureSamplerBinding
type GPUStorageTextureReadWriteBinding = raw::GPUStorageTextureReadWriteBinding
type GPUStorageBufferReadWriteBinding = raw::GPUStorageBufferReadWriteBinding
type GPUComputePass = raw::GPUComputePass
type GPUCopyPass = raw::GPUCopyPass
type GPUTextureTransferInfo = raw::GPUTextureTransferInfo
type GPUTextureRegion = raw::GPUTextureRegion
type GPUTransferBufferLocation = raw::GPUTransferBufferLocation
type GPUBufferRegion = raw::GPUBufferRegion
type GPUTextureLocation = raw::GPUTextureLocation
type GPUBufferLocation = raw::GPUBufferLocation
type GPUBlitRegion = raw::GPUBlitRegion
type GPUBlitInfo = raw::GPUBlitInfo
type GPUSwapchainComposition = raw::GPUSwapchainComposition
type GPUPresentMode = raw::GPUPresentMode
type GPUFence = raw::GPUFence
type Haptic = raw::Haptic
type HapticDirection = raw::HapticDirection
type HapticConstant = raw::HapticConstant
type HapticPeriodic = raw::HapticPeriodic
type HapticCondition = raw::HapticCondition
type HapticRamp = raw::HapticRamp
type HapticLeftRight = raw::HapticLeftRight
type HapticCustom = raw::HapticCustom
type HapticEffect = raw::HapticEffect
type hid_bus_type = raw::hid_bus_type
type hid_device_info = raw::hid_device_info
type hid_device = raw::hid_device
type HintPriority = raw::HintPriority
type SharedObject = raw::SharedObject
type Locale = raw::Locale
type LogPriority = raw::LogPriority
type MessageBoxButtonData = raw::MessageBoxButtonData
type MessageBoxColor = raw::MessageBoxColor
type MessageBoxColorScheme = raw::MessageBoxColorScheme
type MessageBoxData = raw::MessageBoxData
type Process = raw::Process
type Renderer = raw::Renderer
type TextureAccess = raw::TextureAccess
type Texture = raw::Texture
type RendererLogicalPresentation = raw::RendererLogicalPresentation
type Vertex = raw::Vertex
type Storage = raw::Storage
type StorageInterface = raw::StorageInterface
type XEvent = raw::XEvent
type Sandbox = raw::Sandbox
type DateFormat = raw::DateFormat
type TimeFormat = raw::TimeFormat
type DateTime = raw::DateTime
type Tray = raw::Tray
type TrayMenu = raw::TrayMenu
type TrayEntry = raw::TrayEntry
type AppResult = raw::AppResult
type VkInstance_T = raw::VkInstance_T
type VkAllocationCallbacks = raw::VkAllocationCallbacks
type VkSurfaceKHR_T = raw::VkSurfaceKHR_T
type VkPhysicalDevice_T = raw::VkPhysicalDevice_T
type malloc_func = raw::malloc_func
type calloc_func = raw::calloc_func
type realloc_func = raw::realloc_func
type free_func = raw::free_func
type CompareCallback = raw::CompareCallback
type CompareCallback_r = raw::CompareCallback_r
type iconv_t = raw::iconv_t
type AssertionHandler = raw::AssertionHandler
type SpinLock = raw::SpinLock
type PropertiesID = raw::PropertiesID
type CleanupPropertyCallback = raw::CleanupPropertyCallback
type EnumeratePropertiesCallback = raw::EnumeratePropertiesCallback
type ThreadFunction = raw::ThreadFunction
type FunctionPointer = raw::FunctionPointer
type ThreadID = raw::ThreadID
type TLSID = raw::TLSID
type TLSDestructorCallback = raw::TLSDestructorCallback
type IOStreamInterface_sizeCallback = raw::IOStreamInterface_sizeCallback
type IOStreamInterface_seekCallback = raw::IOStreamInterface_seekCallback
type IOStreamInterface_readCallback = raw::IOStreamInterface_readCallback
type IOStreamInterface_writeCallback = raw::IOStreamInterface_writeCallback
type IOStreamInterface_flushCallback = raw::IOStreamInterface_flushCallback
type IOStreamInterface_closeCallback = raw::IOStreamInterface_closeCallback
type AudioDeviceID = raw::AudioDeviceID
type AudioStreamCallback = raw::AudioStreamCallback
type AudioPostmixCallback = raw::AudioPostmixCallback
type BlendMode = raw::BlendMode
type SurfaceFlags = raw::SurfaceFlags
type CameraID = raw::CameraID
type ClipboardDataCallback = raw::ClipboardDataCallback
type ClipboardCleanupCallback = raw::ClipboardCleanupCallback
type DisplayID = raw::DisplayID
type WindowFlags = raw::WindowFlags
type WindowID = raw::WindowID
type HitTest = raw::HitTest
type GLContext = raw::GLContext
type EGLDisplay = raw::EGLDisplay
type EGLConfig = raw::EGLConfig
type EGLSurface = raw::EGLSurface
type EGLAttrib = raw::EGLAttrib
type EGLAttribArrayCallback = raw::EGLAttribArrayCallback
type EGLint = raw::EGLint
type EGLIntArrayCallback = raw::EGLIntArrayCallback
type DialogFileCallback = raw::DialogFileCallback
type SensorID = raw::SensorID
type JoystickID = raw::JoystickID
type VirtualJoystickDesc_UpdateCallback = raw::VirtualJoystickDesc_UpdateCallback
type VirtualJoystickDesc_SetPlayerIndexCallback = raw::VirtualJoystickDesc_SetPlayerIndexCallback
type VirtualJoystickDesc_RumbleCallback = raw::VirtualJoystickDesc_RumbleCallback
type VirtualJoystickDesc_RumbleTriggersCallback = raw::VirtualJoystickDesc_RumbleTriggersCallback
type VirtualJoystickDesc_SetLEDCallback = raw::VirtualJoystickDesc_SetLEDCallback
type VirtualJoystickDesc_SendEffectCallback = raw::VirtualJoystickDesc_SendEffectCallback
type VirtualJoystickDesc_SetSensorsEnabledCallback = raw::VirtualJoystickDesc_SetSensorsEnabledCallback
type VirtualJoystickDesc_CleanupCallback = raw::VirtualJoystickDesc_CleanupCallback
type KeyboardID = raw::KeyboardID
type Keymod = raw::Keymod
type Keycode = raw::Keycode
type MouseID = raw::MouseID
type MouseButtonFlags = raw::MouseButtonFlags
type TouchID = raw::TouchID
type FingerID = raw::FingerID
type PenID = raw::PenID
type PenInputFlags = raw::PenInputFlags
type EventFilter = raw::EventFilter
type EnumerateDirectoryCallback = raw::EnumerateDirectoryCallback
type Time = raw::Time
type GlobFlags = raw::GlobFlags
type GPUShaderFormat = raw::GPUShaderFormat
type GPUColorComponentFlags = raw::GPUColorComponentFlags
type GPUTextureUsageFlags = raw::GPUTextureUsageFlags
type GPUBufferUsageFlags = raw::GPUBufferUsageFlags
type HapticID = raw::HapticID
type HintCallback = raw::HintCallback
type InitFlags = raw::InitFlags
type MainThreadCallback = raw::MainThreadCallback
type LogOutputFunction = raw::LogOutputFunction
type MessageBoxFlags = raw::MessageBoxFlags
type MessageBoxButtonFlags = raw::MessageBoxButtonFlags
type MetalView = raw::MetalView
type StorageInterface_closeCallback = raw::StorageInterface_closeCallback
type StorageInterface_readyCallback = raw::StorageInterface_readyCallback
type StorageInterface_enumerateCallback = raw::StorageInterface_enumerateCallback
type StorageInterface_infoCallback = raw::StorageInterface_infoCallback
type StorageInterface_read_fileCallback = raw::StorageInterface_read_fileCallback
type StorageInterface_write_fileCallback = raw::StorageInterface_write_fileCallback
type StorageInterface_mkdirCallback = raw::StorageInterface_mkdirCallback
type StorageInterface_removeCallback = raw::StorageInterface_removeCallback
type StorageInterface_renameCallback = raw::StorageInterface_renameCallback
type StorageInterface_copyCallback = raw::StorageInterface_copyCallback
type StorageInterface_space_remainingCallback = raw::StorageInterface_space_remainingCallback
type X11EventHook = raw::X11EventHook
type TimerID = raw::TimerID
type TimerCallback = raw::TimerCallback
type NSTimerCallback = raw::NSTimerCallback
type TrayEntryFlags = raw::TrayEntryFlags
type TrayCallback = raw::TrayCallback
type main_func = raw::main_func
type AppInit_func = raw::AppInit_func
type AppIterate_func = raw::AppIterate_func
type AppEvent_func = raw::AppEvent_func
type AppQuit_func = raw::AppQuit_func
type VkInstance = raw::VkInstance
type VkSurfaceKHR = raw::VkSurfaceKHR
type VkPhysicalDevice = raw::VkPhysicalDevice

@[header="SDL3/SDL.h"]
extern {
  @[link_name="SDL_ASSERTION_RETRY"]
  const ASSERTION_RETRY u32

  @[link_name="SDL_ASSERTION_BREAK"]
  const ASSERTION_BREAK u32

  @[link_name="SDL_ASSERTION_ABORT"]
  const ASSERTION_ABORT u32

  @[link_name="SDL_ASSERTION_IGNORE"]
  const ASSERTION_IGNORE u32

  @[link_name="SDL_ASSERTION_ALWAYS_IGNORE"]
  const ASSERTION_ALWAYS_IGNORE u32

  @[link_name="SDL_ASYNCIO_TASK_READ"]
  const ASYNCIO_TASK_READ u32

  @[link_name="SDL_ASYNCIO_TASK_WRITE"]
  const ASYNCIO_TASK_WRITE u32

  @[link_name="SDL_ASYNCIO_TASK_CLOSE"]
  const ASYNCIO_TASK_CLOSE u32

  @[link_name="SDL_ASYNCIO_COMPLETE"]
  const ASYNCIO_COMPLETE u32

  @[link_name="SDL_ASYNCIO_FAILURE"]
  const ASYNCIO_FAILURE u32

  @[link_name="SDL_ASYNCIO_CANCELED"]
  const ASYNCIO_CANCELED u32

  @[link_name="SDL_PROPERTY_TYPE_INVALID"]
  const PROPERTY_TYPE_INVALID u32

  @[link_name="SDL_PROPERTY_TYPE_POINTER"]
  const PROPERTY_TYPE_POINTER u32

  @[link_name="SDL_PROPERTY_TYPE_STRING"]
  const PROPERTY_TYPE_STRING u32

  @[link_name="SDL_PROPERTY_TYPE_NUMBER"]
  const PROPERTY_TYPE_NUMBER u32

  @[link_name="SDL_PROPERTY_TYPE_FLOAT"]
  const PROPERTY_TYPE_FLOAT u32

  @[link_name="SDL_PROPERTY_TYPE_BOOLEAN"]
  const PROPERTY_TYPE_BOOLEAN u32

  @[link_name="SDL_THREAD_PRIORITY_LOW"]
  const THREAD_PRIORITY_LOW u32

  @[link_name="SDL_THREAD_PRIORITY_NORMAL"]
  const THREAD_PRIORITY_NORMAL u32

  @[link_name="SDL_THREAD_PRIORITY_HIGH"]
  const THREAD_PRIORITY_HIGH u32

  @[link_name="SDL_THREAD_PRIORITY_TIME_CRITICAL"]
  const THREAD_PRIORITY_TIME_CRITICAL u32

  @[link_name="SDL_THREAD_UNKNOWN"]
  const THREAD_UNKNOWN u32

  @[link_name="SDL_THREAD_ALIVE"]
  const THREAD_ALIVE u32

  @[link_name="SDL_THREAD_DETACHED"]
  const THREAD_DETACHED u32

  @[link_name="SDL_THREAD_COMPLETE"]
  const THREAD_COMPLETE u32

  @[link_name="SDL_INIT_STATUS_UNINITIALIZED"]
  const INIT_STATUS_UNINITIALIZED u32

  @[link_name="SDL_INIT_STATUS_INITIALIZING"]
  const INIT_STATUS_INITIALIZING u32

  @[link_name="SDL_INIT_STATUS_INITIALIZED"]
  const INIT_STATUS_INITIALIZED u32

  @[link_name="SDL_INIT_STATUS_UNINITIALIZING"]
  const INIT_STATUS_UNINITIALIZING u32

  @[link_name="SDL_IO_STATUS_READY"]
  const IO_STATUS_READY u32

  @[link_name="SDL_IO_STATUS_ERROR"]
  const IO_STATUS_ERROR u32

  @[link_name="SDL_IO_STATUS_EOF"]
  const IO_STATUS_EOF u32

  @[link_name="SDL_IO_STATUS_NOT_READY"]
  const IO_STATUS_NOT_READY u32

  @[link_name="SDL_IO_STATUS_READONLY"]
  const IO_STATUS_READONLY u32

  @[link_name="SDL_IO_STATUS_WRITEONLY"]
  const IO_STATUS_WRITEONLY u32

  @[link_name="SDL_IO_SEEK_SET"]
  const IO_SEEK_SET u32

  @[link_name="SDL_IO_SEEK_CUR"]
  const IO_SEEK_CUR u32

  @[link_name="SDL_IO_SEEK_END"]
  const IO_SEEK_END u32

  @[link_name="SDL_AUDIO_UNKNOWN"]
  const AUDIO_UNKNOWN u32

  @[link_name="SDL_AUDIO_U8"]
  const AUDIO_U8 u32

  @[link_name="SDL_AUDIO_S8"]
  const AUDIO_S8 u32

  @[link_name="SDL_AUDIO_S16LE"]
  const AUDIO_S16LE u32

  @[link_name="SDL_AUDIO_S16BE"]
  const AUDIO_S16BE u32

  @[link_name="SDL_AUDIO_S32LE"]
  const AUDIO_S32LE u32

  @[link_name="SDL_AUDIO_S32BE"]
  const AUDIO_S32BE u32

  @[link_name="SDL_AUDIO_F32LE"]
  const AUDIO_F32LE u32

  @[link_name="SDL_AUDIO_F32BE"]
  const AUDIO_F32BE u32

  @[link_name="SDL_AUDIO_S16"]
  const AUDIO_S16 u32

  @[link_name="SDL_AUDIO_S32"]
  const AUDIO_S32 u32

  @[link_name="SDL_AUDIO_F32"]
  const AUDIO_F32 u32

  @[link_name="SDL_BLENDOPERATION_ADD"]
  const BLENDOPERATION_ADD u32

  @[link_name="SDL_BLENDOPERATION_SUBTRACT"]
  const BLENDOPERATION_SUBTRACT u32

  @[link_name="SDL_BLENDOPERATION_REV_SUBTRACT"]
  const BLENDOPERATION_REV_SUBTRACT u32

  @[link_name="SDL_BLENDOPERATION_MINIMUM"]
  const BLENDOPERATION_MINIMUM u32

  @[link_name="SDL_BLENDOPERATION_MAXIMUM"]
  const BLENDOPERATION_MAXIMUM u32

  @[link_name="SDL_BLENDFACTOR_ZERO"]
  const BLENDFACTOR_ZERO u32

  @[link_name="SDL_BLENDFACTOR_ONE"]
  const BLENDFACTOR_ONE u32

  @[link_name="SDL_BLENDFACTOR_SRC_COLOR"]
  const BLENDFACTOR_SRC_COLOR u32

  @[link_name="SDL_BLENDFACTOR_ONE_MINUS_SRC_COLOR"]
  const BLENDFACTOR_ONE_MINUS_SRC_COLOR u32

  @[link_name="SDL_BLENDFACTOR_SRC_ALPHA"]
  const BLENDFACTOR_SRC_ALPHA u32

  @[link_name="SDL_BLENDFACTOR_ONE_MINUS_SRC_ALPHA"]
  const BLENDFACTOR_ONE_MINUS_SRC_ALPHA u32

  @[link_name="SDL_BLENDFACTOR_DST_COLOR"]
  const BLENDFACTOR_DST_COLOR u32

  @[link_name="SDL_BLENDFACTOR_ONE_MINUS_DST_COLOR"]
  const BLENDFACTOR_ONE_MINUS_DST_COLOR u32

  @[link_name="SDL_BLENDFACTOR_DST_ALPHA"]
  const BLENDFACTOR_DST_ALPHA u32

  @[link_name="SDL_BLENDFACTOR_ONE_MINUS_DST_ALPHA"]
  const BLENDFACTOR_ONE_MINUS_DST_ALPHA u32

  @[link_name="SDL_PIXELTYPE_UNKNOWN"]
  const PIXELTYPE_UNKNOWN u32

  @[link_name="SDL_PIXELTYPE_INDEX1"]
  const PIXELTYPE_INDEX1 u32

  @[link_name="SDL_PIXELTYPE_INDEX4"]
  const PIXELTYPE_INDEX4 u32

  @[link_name="SDL_PIXELTYPE_INDEX8"]
  const PIXELTYPE_INDEX8 u32

  @[link_name="SDL_PIXELTYPE_PACKED8"]
  const PIXELTYPE_PACKED8 u32

  @[link_name="SDL_PIXELTYPE_PACKED16"]
  const PIXELTYPE_PACKED16 u32

  @[link_name="SDL_PIXELTYPE_PACKED32"]
  const PIXELTYPE_PACKED32 u32

  @[link_name="SDL_PIXELTYPE_ARRAYU8"]
  const PIXELTYPE_ARRAYU8 u32

  @[link_name="SDL_PIXELTYPE_ARRAYU16"]
  const PIXELTYPE_ARRAYU16 u32

  @[link_name="SDL_PIXELTYPE_ARRAYU32"]
  const PIXELTYPE_ARRAYU32 u32

  @[link_name="SDL_PIXELTYPE_ARRAYF16"]
  const PIXELTYPE_ARRAYF16 u32

  @[link_name="SDL_PIXELTYPE_ARRAYF32"]
  const PIXELTYPE_ARRAYF32 u32

  @[link_name="SDL_PIXELTYPE_INDEX2"]
  const PIXELTYPE_INDEX2 u32

  @[link_name="SDL_BITMAPORDER_NONE"]
  const BITMAPORDER_NONE u32

  @[link_name="SDL_BITMAPORDER_4321"]
  const BITMAPORDER_4321 u32

  @[link_name="SDL_BITMAPORDER_1234"]
  const BITMAPORDER_1234 u32

  @[link_name="SDL_PACKEDORDER_NONE"]
  const PACKEDORDER_NONE u32

  @[link_name="SDL_PACKEDORDER_XRGB"]
  const PACKEDORDER_XRGB u32

  @[link_name="SDL_PACKEDORDER_RGBX"]
  const PACKEDORDER_RGBX u32

  @[link_name="SDL_PACKEDORDER_ARGB"]
  const PACKEDORDER_ARGB u32

  @[link_name="SDL_PACKEDORDER_RGBA"]
  const PACKEDORDER_RGBA u32

  @[link_name="SDL_PACKEDORDER_XBGR"]
  const PACKEDORDER_XBGR u32

  @[link_name="SDL_PACKEDORDER_BGRX"]
  const PACKEDORDER_BGRX u32

  @[link_name="SDL_PACKEDORDER_ABGR"]
  const PACKEDORDER_ABGR u32

  @[link_name="SDL_PACKEDORDER_BGRA"]
  const PACKEDORDER_BGRA u32

  @[link_name="SDL_ARRAYORDER_NONE"]
  const ARRAYORDER_NONE u32

  @[link_name="SDL_ARRAYORDER_RGB"]
  const ARRAYORDER_RGB u32

  @[link_name="SDL_ARRAYORDER_RGBA"]
  const ARRAYORDER_RGBA u32

  @[link_name="SDL_ARRAYORDER_ARGB"]
  const ARRAYORDER_ARGB u32

  @[link_name="SDL_ARRAYORDER_BGR"]
  const ARRAYORDER_BGR u32

  @[link_name="SDL_ARRAYORDER_BGRA"]
  const ARRAYORDER_BGRA u32

  @[link_name="SDL_ARRAYORDER_ABGR"]
  const ARRAYORDER_ABGR u32

  @[link_name="SDL_PACKEDLAYOUT_NONE"]
  const PACKEDLAYOUT_NONE u32

  @[link_name="SDL_PACKEDLAYOUT_332"]
  const PACKEDLAYOUT_332 u32

  @[link_name="SDL_PACKEDLAYOUT_4444"]
  const PACKEDLAYOUT_4444 u32

  @[link_name="SDL_PACKEDLAYOUT_1555"]
  const PACKEDLAYOUT_1555 u32

  @[link_name="SDL_PACKEDLAYOUT_5551"]
  const PACKEDLAYOUT_5551 u32

  @[link_name="SDL_PACKEDLAYOUT_565"]
  const PACKEDLAYOUT_565 u32

  @[link_name="SDL_PACKEDLAYOUT_8888"]
  const PACKEDLAYOUT_8888 u32

  @[link_name="SDL_PACKEDLAYOUT_2101010"]
  const PACKEDLAYOUT_2101010 u32

  @[link_name="SDL_PACKEDLAYOUT_1010102"]
  const PACKEDLAYOUT_1010102 u32

  @[link_name="SDL_PIXELFORMAT_UNKNOWN"]
  const PIXELFORMAT_UNKNOWN u32

  @[link_name="SDL_PIXELFORMAT_INDEX1LSB"]
  const PIXELFORMAT_INDEX1LSB u32

  @[link_name="SDL_PIXELFORMAT_INDEX1MSB"]
  const PIXELFORMAT_INDEX1MSB u32

  @[link_name="SDL_PIXELFORMAT_INDEX2LSB"]
  const PIXELFORMAT_INDEX2LSB u32

  @[link_name="SDL_PIXELFORMAT_INDEX2MSB"]
  const PIXELFORMAT_INDEX2MSB u32

  @[link_name="SDL_PIXELFORMAT_INDEX4LSB"]
  const PIXELFORMAT_INDEX4LSB u32

  @[link_name="SDL_PIXELFORMAT_INDEX4MSB"]
  const PIXELFORMAT_INDEX4MSB u32

  @[link_name="SDL_PIXELFORMAT_INDEX8"]
  const PIXELFORMAT_INDEX8 u32

  @[link_name="SDL_PIXELFORMAT_RGB332"]
  const PIXELFORMAT_RGB332 u32

  @[link_name="SDL_PIXELFORMAT_XRGB4444"]
  const PIXELFORMAT_XRGB4444 u32

  @[link_name="SDL_PIXELFORMAT_XBGR4444"]
  const PIXELFORMAT_XBGR4444 u32

  @[link_name="SDL_PIXELFORMAT_XRGB1555"]
  const PIXELFORMAT_XRGB1555 u32

  @[link_name="SDL_PIXELFORMAT_XBGR1555"]
  const PIXELFORMAT_XBGR1555 u32

  @[link_name="SDL_PIXELFORMAT_ARGB4444"]
  const PIXELFORMAT_ARGB4444 u32

  @[link_name="SDL_PIXELFORMAT_RGBA4444"]
  const PIXELFORMAT_RGBA4444 u32

  @[link_name="SDL_PIXELFORMAT_ABGR4444"]
  const PIXELFORMAT_ABGR4444 u32

  @[link_name="SDL_PIXELFORMAT_BGRA4444"]
  const PIXELFORMAT_BGRA4444 u32

  @[link_name="SDL_PIXELFORMAT_ARGB1555"]
  const PIXELFORMAT_ARGB1555 u32

  @[link_name="SDL_PIXELFORMAT_RGBA5551"]
  const PIXELFORMAT_RGBA5551 u32

  @[link_name="SDL_PIXELFORMAT_ABGR1555"]
  const PIXELFORMAT_ABGR1555 u32

  @[link_name="SDL_PIXELFORMAT_BGRA5551"]
  const PIXELFORMAT_BGRA5551 u32

  @[link_name="SDL_PIXELFORMAT_RGB565"]
  const PIXELFORMAT_RGB565 u32

  @[link_name="SDL_PIXELFORMAT_BGR565"]
  const PIXELFORMAT_BGR565 u32

  @[link_name="SDL_PIXELFORMAT_RGB24"]
  const PIXELFORMAT_RGB24 u32

  @[link_name="SDL_PIXELFORMAT_BGR24"]
  const PIXELFORMAT_BGR24 u32

  @[link_name="SDL_PIXELFORMAT_XRGB8888"]
  const PIXELFORMAT_XRGB8888 u32

  @[link_name="SDL_PIXELFORMAT_RGBX8888"]
  const PIXELFORMAT_RGBX8888 u32

  @[link_name="SDL_PIXELFORMAT_XBGR8888"]
  const PIXELFORMAT_XBGR8888 u32

  @[link_name="SDL_PIXELFORMAT_BGRX8888"]
  const PIXELFORMAT_BGRX8888 u32

  @[link_name="SDL_PIXELFORMAT_ARGB8888"]
  const PIXELFORMAT_ARGB8888 u32

  @[link_name="SDL_PIXELFORMAT_RGBA8888"]
  const PIXELFORMAT_RGBA8888 u32

  @[link_name="SDL_PIXELFORMAT_ABGR8888"]
  const PIXELFORMAT_ABGR8888 u32

  @[link_name="SDL_PIXELFORMAT_BGRA8888"]
  const PIXELFORMAT_BGRA8888 u32

  @[link_name="SDL_PIXELFORMAT_XRGB2101010"]
  const PIXELFORMAT_XRGB2101010 u32

  @[link_name="SDL_PIXELFORMAT_XBGR2101010"]
  const PIXELFORMAT_XBGR2101010 u32

  @[link_name="SDL_PIXELFORMAT_ARGB2101010"]
  const PIXELFORMAT_ARGB2101010 u32

  @[link_name="SDL_PIXELFORMAT_ABGR2101010"]
  const PIXELFORMAT_ABGR2101010 u32

  @[link_name="SDL_PIXELFORMAT_RGB48"]
  const PIXELFORMAT_RGB48 u32

  @[link_name="SDL_PIXELFORMAT_BGR48"]
  const PIXELFORMAT_BGR48 u32

  @[link_name="SDL_PIXELFORMAT_RGBA64"]
  const PIXELFORMAT_RGBA64 u32

  @[link_name="SDL_PIXELFORMAT_ARGB64"]
  const PIXELFORMAT_ARGB64 u32

  @[link_name="SDL_PIXELFORMAT_BGRA64"]
  const PIXELFORMAT_BGRA64 u32

  @[link_name="SDL_PIXELFORMAT_ABGR64"]
  const PIXELFORMAT_ABGR64 u32

  @[link_name="SDL_PIXELFORMAT_RGB48_FLOAT"]
  const PIXELFORMAT_RGB48_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_BGR48_FLOAT"]
  const PIXELFORMAT_BGR48_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_RGBA64_FLOAT"]
  const PIXELFORMAT_RGBA64_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_ARGB64_FLOAT"]
  const PIXELFORMAT_ARGB64_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_BGRA64_FLOAT"]
  const PIXELFORMAT_BGRA64_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_ABGR64_FLOAT"]
  const PIXELFORMAT_ABGR64_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_RGB96_FLOAT"]
  const PIXELFORMAT_RGB96_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_BGR96_FLOAT"]
  const PIXELFORMAT_BGR96_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_RGBA128_FLOAT"]
  const PIXELFORMAT_RGBA128_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_ARGB128_FLOAT"]
  const PIXELFORMAT_ARGB128_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_BGRA128_FLOAT"]
  const PIXELFORMAT_BGRA128_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_ABGR128_FLOAT"]
  const PIXELFORMAT_ABGR128_FLOAT u32

  @[link_name="SDL_PIXELFORMAT_YV12"]
  const PIXELFORMAT_YV12 u32

  @[link_name="SDL_PIXELFORMAT_IYUV"]
  const PIXELFORMAT_IYUV u32

  @[link_name="SDL_PIXELFORMAT_YUY2"]
  const PIXELFORMAT_YUY2 u32

  @[link_name="SDL_PIXELFORMAT_UYVY"]
  const PIXELFORMAT_UYVY u32

  @[link_name="SDL_PIXELFORMAT_YVYU"]
  const PIXELFORMAT_YVYU u32

  @[link_name="SDL_PIXELFORMAT_NV12"]
  const PIXELFORMAT_NV12 u32

  @[link_name="SDL_PIXELFORMAT_NV21"]
  const PIXELFORMAT_NV21 u32

  @[link_name="SDL_PIXELFORMAT_P010"]
  const PIXELFORMAT_P010 u32

  @[link_name="SDL_PIXELFORMAT_EXTERNAL_OES"]
  const PIXELFORMAT_EXTERNAL_OES u32

  @[link_name="SDL_PIXELFORMAT_MJPG"]
  const PIXELFORMAT_MJPG u32

  @[link_name="SDL_PIXELFORMAT_RGBA32"]
  const PIXELFORMAT_RGBA32 u32

  @[link_name="SDL_PIXELFORMAT_ARGB32"]
  const PIXELFORMAT_ARGB32 u32

  @[link_name="SDL_PIXELFORMAT_BGRA32"]
  const PIXELFORMAT_BGRA32 u32

  @[link_name="SDL_PIXELFORMAT_ABGR32"]
  const PIXELFORMAT_ABGR32 u32

  @[link_name="SDL_PIXELFORMAT_RGBX32"]
  const PIXELFORMAT_RGBX32 u32

  @[link_name="SDL_PIXELFORMAT_XRGB32"]
  const PIXELFORMAT_XRGB32 u32

  @[link_name="SDL_PIXELFORMAT_BGRX32"]
  const PIXELFORMAT_BGRX32 u32

  @[link_name="SDL_PIXELFORMAT_XBGR32"]
  const PIXELFORMAT_XBGR32 u32

  @[link_name="SDL_COLOR_TYPE_UNKNOWN"]
  const COLOR_TYPE_UNKNOWN u32

  @[link_name="SDL_COLOR_TYPE_RGB"]
  const COLOR_TYPE_RGB u32

  @[link_name="SDL_COLOR_TYPE_YCBCR"]
  const COLOR_TYPE_YCBCR u32

  @[link_name="SDL_COLOR_RANGE_UNKNOWN"]
  const COLOR_RANGE_UNKNOWN u32

  @[link_name="SDL_COLOR_RANGE_LIMITED"]
  const COLOR_RANGE_LIMITED u32

  @[link_name="SDL_COLOR_RANGE_FULL"]
  const COLOR_RANGE_FULL u32

  @[link_name="SDL_COLOR_PRIMARIES_UNKNOWN"]
  const COLOR_PRIMARIES_UNKNOWN u32

  @[link_name="SDL_COLOR_PRIMARIES_BT709"]
  const COLOR_PRIMARIES_BT709 u32

  @[link_name="SDL_COLOR_PRIMARIES_UNSPECIFIED"]
  const COLOR_PRIMARIES_UNSPECIFIED u32

  @[link_name="SDL_COLOR_PRIMARIES_BT470M"]
  const COLOR_PRIMARIES_BT470M u32

  @[link_name="SDL_COLOR_PRIMARIES_BT470BG"]
  const COLOR_PRIMARIES_BT470BG u32

  @[link_name="SDL_COLOR_PRIMARIES_BT601"]
  const COLOR_PRIMARIES_BT601 u32

  @[link_name="SDL_COLOR_PRIMARIES_SMPTE240"]
  const COLOR_PRIMARIES_SMPTE240 u32

  @[link_name="SDL_COLOR_PRIMARIES_GENERIC_FILM"]
  const COLOR_PRIMARIES_GENERIC_FILM u32

  @[link_name="SDL_COLOR_PRIMARIES_BT2020"]
  const COLOR_PRIMARIES_BT2020 u32

  @[link_name="SDL_COLOR_PRIMARIES_XYZ"]
  const COLOR_PRIMARIES_XYZ u32

  @[link_name="SDL_COLOR_PRIMARIES_SMPTE431"]
  const COLOR_PRIMARIES_SMPTE431 u32

  @[link_name="SDL_COLOR_PRIMARIES_SMPTE432"]
  const COLOR_PRIMARIES_SMPTE432 u32

  @[link_name="SDL_COLOR_PRIMARIES_EBU3213"]
  const COLOR_PRIMARIES_EBU3213 u32

  @[link_name="SDL_COLOR_PRIMARIES_CUSTOM"]
  const COLOR_PRIMARIES_CUSTOM u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_UNKNOWN"]
  const TRANSFER_CHARACTERISTICS_UNKNOWN u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_BT709"]
  const TRANSFER_CHARACTERISTICS_BT709 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_UNSPECIFIED"]
  const TRANSFER_CHARACTERISTICS_UNSPECIFIED u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_GAMMA22"]
  const TRANSFER_CHARACTERISTICS_GAMMA22 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_GAMMA28"]
  const TRANSFER_CHARACTERISTICS_GAMMA28 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_BT601"]
  const TRANSFER_CHARACTERISTICS_BT601 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_SMPTE240"]
  const TRANSFER_CHARACTERISTICS_SMPTE240 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_LINEAR"]
  const TRANSFER_CHARACTERISTICS_LINEAR u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_LOG100"]
  const TRANSFER_CHARACTERISTICS_LOG100 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_LOG100_SQRT10"]
  const TRANSFER_CHARACTERISTICS_LOG100_SQRT10 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_IEC61966"]
  const TRANSFER_CHARACTERISTICS_IEC61966 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_BT1361"]
  const TRANSFER_CHARACTERISTICS_BT1361 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_SRGB"]
  const TRANSFER_CHARACTERISTICS_SRGB u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_BT2020_10BIT"]
  const TRANSFER_CHARACTERISTICS_BT2020_10BIT u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_BT2020_12BIT"]
  const TRANSFER_CHARACTERISTICS_BT2020_12BIT u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_PQ"]
  const TRANSFER_CHARACTERISTICS_PQ u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_SMPTE428"]
  const TRANSFER_CHARACTERISTICS_SMPTE428 u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_HLG"]
  const TRANSFER_CHARACTERISTICS_HLG u32

  @[link_name="SDL_TRANSFER_CHARACTERISTICS_CUSTOM"]
  const TRANSFER_CHARACTERISTICS_CUSTOM u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_IDENTITY"]
  const MATRIX_COEFFICIENTS_IDENTITY u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_BT709"]
  const MATRIX_COEFFICIENTS_BT709 u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_UNSPECIFIED"]
  const MATRIX_COEFFICIENTS_UNSPECIFIED u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_FCC"]
  const MATRIX_COEFFICIENTS_FCC u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_BT470BG"]
  const MATRIX_COEFFICIENTS_BT470BG u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_BT601"]
  const MATRIX_COEFFICIENTS_BT601 u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_SMPTE240"]
  const MATRIX_COEFFICIENTS_SMPTE240 u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_YCGCO"]
  const MATRIX_COEFFICIENTS_YCGCO u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_BT2020_NCL"]
  const MATRIX_COEFFICIENTS_BT2020_NCL u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_BT2020_CL"]
  const MATRIX_COEFFICIENTS_BT2020_CL u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_SMPTE2085"]
  const MATRIX_COEFFICIENTS_SMPTE2085 u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_CHROMA_DERIVED_NCL"]
  const MATRIX_COEFFICIENTS_CHROMA_DERIVED_NCL u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_CHROMA_DERIVED_CL"]
  const MATRIX_COEFFICIENTS_CHROMA_DERIVED_CL u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_ICTCP"]
  const MATRIX_COEFFICIENTS_ICTCP u32

  @[link_name="SDL_MATRIX_COEFFICIENTS_CUSTOM"]
  const MATRIX_COEFFICIENTS_CUSTOM u32

  @[link_name="SDL_CHROMA_LOCATION_NONE"]
  const CHROMA_LOCATION_NONE u32

  @[link_name="SDL_CHROMA_LOCATION_LEFT"]
  const CHROMA_LOCATION_LEFT u32

  @[link_name="SDL_CHROMA_LOCATION_CENTER"]
  const CHROMA_LOCATION_CENTER u32

  @[link_name="SDL_CHROMA_LOCATION_TOPLEFT"]
  const CHROMA_LOCATION_TOPLEFT u32

  @[link_name="SDL_COLORSPACE_UNKNOWN"]
  const COLORSPACE_UNKNOWN u32

  @[link_name="SDL_COLORSPACE_SRGB"]
  const COLORSPACE_SRGB u32

  @[link_name="SDL_COLORSPACE_SRGB_LINEAR"]
  const COLORSPACE_SRGB_LINEAR u32

  @[link_name="SDL_COLORSPACE_HDR10"]
  const COLORSPACE_HDR10 u32

  @[link_name="SDL_COLORSPACE_JPEG"]
  const COLORSPACE_JPEG u32

  @[link_name="SDL_COLORSPACE_BT601_LIMITED"]
  const COLORSPACE_BT601_LIMITED u32

  @[link_name="SDL_COLORSPACE_BT601_FULL"]
  const COLORSPACE_BT601_FULL u32

  @[link_name="SDL_COLORSPACE_BT709_LIMITED"]
  const COLORSPACE_BT709_LIMITED u32

  @[link_name="SDL_COLORSPACE_BT709_FULL"]
  const COLORSPACE_BT709_FULL u32

  @[link_name="SDL_COLORSPACE_BT2020_LIMITED"]
  const COLORSPACE_BT2020_LIMITED u32

  @[link_name="SDL_COLORSPACE_BT2020_FULL"]
  const COLORSPACE_BT2020_FULL u32

  @[link_name="SDL_COLORSPACE_RGB_DEFAULT"]
  const COLORSPACE_RGB_DEFAULT u32

  @[link_name="SDL_COLORSPACE_YUV_DEFAULT"]
  const COLORSPACE_YUV_DEFAULT u32

  @[link_name="SDL_SCALEMODE_INVALID"]
  const SCALEMODE_INVALID u32

  @[link_name="SDL_SCALEMODE_NEAREST"]
  const SCALEMODE_NEAREST u32

  @[link_name="SDL_SCALEMODE_LINEAR"]
  const SCALEMODE_LINEAR u32

  @[link_name="SDL_FLIP_NONE"]
  const FLIP_NONE u32

  @[link_name="SDL_FLIP_HORIZONTAL"]
  const FLIP_HORIZONTAL u32

  @[link_name="SDL_FLIP_VERTICAL"]
  const FLIP_VERTICAL u32

  @[link_name="SDL_CAMERA_POSITION_UNKNOWN"]
  const CAMERA_POSITION_UNKNOWN u32

  @[link_name="SDL_CAMERA_POSITION_FRONT_FACING"]
  const CAMERA_POSITION_FRONT_FACING u32

  @[link_name="SDL_CAMERA_POSITION_BACK_FACING"]
  const CAMERA_POSITION_BACK_FACING u32

  @[link_name="SDL_SYSTEM_THEME_UNKNOWN"]
  const SYSTEM_THEME_UNKNOWN u32

  @[link_name="SDL_SYSTEM_THEME_LIGHT"]
  const SYSTEM_THEME_LIGHT u32

  @[link_name="SDL_SYSTEM_THEME_DARK"]
  const SYSTEM_THEME_DARK u32

  @[link_name="SDL_ORIENTATION_UNKNOWN"]
  const ORIENTATION_UNKNOWN u32

  @[link_name="SDL_ORIENTATION_LANDSCAPE"]
  const ORIENTATION_LANDSCAPE u32

  @[link_name="SDL_ORIENTATION_LANDSCAPE_FLIPPED"]
  const ORIENTATION_LANDSCAPE_FLIPPED u32

  @[link_name="SDL_ORIENTATION_PORTRAIT"]
  const ORIENTATION_PORTRAIT u32

  @[link_name="SDL_ORIENTATION_PORTRAIT_FLIPPED"]
  const ORIENTATION_PORTRAIT_FLIPPED u32

  @[link_name="SDL_FLASH_CANCEL"]
  const FLASH_CANCEL u32

  @[link_name="SDL_FLASH_BRIEFLY"]
  const FLASH_BRIEFLY u32

  @[link_name="SDL_FLASH_UNTIL_FOCUSED"]
  const FLASH_UNTIL_FOCUSED u32

  @[link_name="SDL_GL_RED_SIZE"]
  const GL_RED_SIZE u32

  @[link_name="SDL_GL_GREEN_SIZE"]
  const GL_GREEN_SIZE u32

  @[link_name="SDL_GL_BLUE_SIZE"]
  const GL_BLUE_SIZE u32

  @[link_name="SDL_GL_ALPHA_SIZE"]
  const GL_ALPHA_SIZE u32

  @[link_name="SDL_GL_BUFFER_SIZE"]
  const GL_BUFFER_SIZE u32

  @[link_name="SDL_GL_DOUBLEBUFFER"]
  const GL_DOUBLEBUFFER u32

  @[link_name="SDL_GL_DEPTH_SIZE"]
  const GL_DEPTH_SIZE u32

  @[link_name="SDL_GL_STENCIL_SIZE"]
  const GL_STENCIL_SIZE u32

  @[link_name="SDL_GL_ACCUM_RED_SIZE"]
  const GL_ACCUM_RED_SIZE u32

  @[link_name="SDL_GL_ACCUM_GREEN_SIZE"]
  const GL_ACCUM_GREEN_SIZE u32

  @[link_name="SDL_GL_ACCUM_BLUE_SIZE"]
  const GL_ACCUM_BLUE_SIZE u32

  @[link_name="SDL_GL_ACCUM_ALPHA_SIZE"]
  const GL_ACCUM_ALPHA_SIZE u32

  @[link_name="SDL_GL_STEREO"]
  const GL_STEREO u32

  @[link_name="SDL_GL_MULTISAMPLEBUFFERS"]
  const GL_MULTISAMPLEBUFFERS u32

  @[link_name="SDL_GL_MULTISAMPLESAMPLES"]
  const GL_MULTISAMPLESAMPLES u32

  @[link_name="SDL_GL_ACCELERATED_VISUAL"]
  const GL_ACCELERATED_VISUAL u32

  @[link_name="SDL_GL_RETAINED_BACKING"]
  const GL_RETAINED_BACKING u32

  @[link_name="SDL_GL_CONTEXT_MAJOR_VERSION"]
  const GL_CONTEXT_MAJOR_VERSION u32

  @[link_name="SDL_GL_CONTEXT_MINOR_VERSION"]
  const GL_CONTEXT_MINOR_VERSION u32

  @[link_name="SDL_GL_CONTEXT_FLAGS"]
  const GL_CONTEXT_FLAGS u32

  @[link_name="SDL_GL_CONTEXT_PROFILE_MASK"]
  const GL_CONTEXT_PROFILE_MASK u32

  @[link_name="SDL_GL_SHARE_WITH_CURRENT_CONTEXT"]
  const GL_SHARE_WITH_CURRENT_CONTEXT u32

  @[link_name="SDL_GL_FRAMEBUFFER_SRGB_CAPABLE"]
  const GL_FRAMEBUFFER_SRGB_CAPABLE u32

  @[link_name="SDL_GL_CONTEXT_RELEASE_BEHAVIOR"]
  const GL_CONTEXT_RELEASE_BEHAVIOR u32

  @[link_name="SDL_GL_CONTEXT_RESET_NOTIFICATION"]
  const GL_CONTEXT_RESET_NOTIFICATION u32

  @[link_name="SDL_GL_CONTEXT_NO_ERROR"]
  const GL_CONTEXT_NO_ERROR u32

  @[link_name="SDL_GL_FLOATBUFFERS"]
  const GL_FLOATBUFFERS u32

  @[link_name="SDL_GL_EGL_PLATFORM"]
  const GL_EGL_PLATFORM u32

  @[link_name="SDL_HITTEST_NORMAL"]
  const HITTEST_NORMAL u32

  @[link_name="SDL_HITTEST_DRAGGABLE"]
  const HITTEST_DRAGGABLE u32

  @[link_name="SDL_HITTEST_RESIZE_TOPLEFT"]
  const HITTEST_RESIZE_TOPLEFT u32

  @[link_name="SDL_HITTEST_RESIZE_TOP"]
  const HITTEST_RESIZE_TOP u32

  @[link_name="SDL_HITTEST_RESIZE_TOPRIGHT"]
  const HITTEST_RESIZE_TOPRIGHT u32

  @[link_name="SDL_HITTEST_RESIZE_RIGHT"]
  const HITTEST_RESIZE_RIGHT u32

  @[link_name="SDL_HITTEST_RESIZE_BOTTOMRIGHT"]
  const HITTEST_RESIZE_BOTTOMRIGHT u32

  @[link_name="SDL_HITTEST_RESIZE_BOTTOM"]
  const HITTEST_RESIZE_BOTTOM u32

  @[link_name="SDL_HITTEST_RESIZE_BOTTOMLEFT"]
  const HITTEST_RESIZE_BOTTOMLEFT u32

  @[link_name="SDL_HITTEST_RESIZE_LEFT"]
  const HITTEST_RESIZE_LEFT u32

  @[link_name="SDL_FILEDIALOG_OPENFILE"]
  const FILEDIALOG_OPENFILE u32

  @[link_name="SDL_FILEDIALOG_SAVEFILE"]
  const FILEDIALOG_SAVEFILE u32

  @[link_name="SDL_FILEDIALOG_OPENFOLDER"]
  const FILEDIALOG_OPENFOLDER u32

  @[link_name="SDL_POWERSTATE_ERROR"]
  const POWERSTATE_ERROR u32

  @[link_name="SDL_POWERSTATE_UNKNOWN"]
  const POWERSTATE_UNKNOWN u32

  @[link_name="SDL_POWERSTATE_ON_BATTERY"]
  const POWERSTATE_ON_BATTERY u32

  @[link_name="SDL_POWERSTATE_NO_BATTERY"]
  const POWERSTATE_NO_BATTERY u32

  @[link_name="SDL_POWERSTATE_CHARGING"]
  const POWERSTATE_CHARGING u32

  @[link_name="SDL_POWERSTATE_CHARGED"]
  const POWERSTATE_CHARGED u32

  @[link_name="SDL_SENSOR_INVALID"]
  const SENSOR_INVALID u32

  @[link_name="SDL_SENSOR_UNKNOWN"]
  const SENSOR_UNKNOWN u32

  @[link_name="SDL_SENSOR_ACCEL"]
  const SENSOR_ACCEL u32

  @[link_name="SDL_SENSOR_GYRO"]
  const SENSOR_GYRO u32

  @[link_name="SDL_SENSOR_ACCEL_L"]
  const SENSOR_ACCEL_L u32

  @[link_name="SDL_SENSOR_GYRO_L"]
  const SENSOR_GYRO_L u32

  @[link_name="SDL_SENSOR_ACCEL_R"]
  const SENSOR_ACCEL_R u32

  @[link_name="SDL_SENSOR_GYRO_R"]
  const SENSOR_GYRO_R u32

  @[link_name="SDL_SENSOR_COUNT"]
  const SENSOR_COUNT u32

  @[link_name="SDL_JOYSTICK_TYPE_UNKNOWN"]
  const JOYSTICK_TYPE_UNKNOWN u32

  @[link_name="SDL_JOYSTICK_TYPE_GAMEPAD"]
  const JOYSTICK_TYPE_GAMEPAD u32

  @[link_name="SDL_JOYSTICK_TYPE_WHEEL"]
  const JOYSTICK_TYPE_WHEEL u32

  @[link_name="SDL_JOYSTICK_TYPE_ARCADE_STICK"]
  const JOYSTICK_TYPE_ARCADE_STICK u32

  @[link_name="SDL_JOYSTICK_TYPE_FLIGHT_STICK"]
  const JOYSTICK_TYPE_FLIGHT_STICK u32

  @[link_name="SDL_JOYSTICK_TYPE_DANCE_PAD"]
  const JOYSTICK_TYPE_DANCE_PAD u32

  @[link_name="SDL_JOYSTICK_TYPE_GUITAR"]
  const JOYSTICK_TYPE_GUITAR u32

  @[link_name="SDL_JOYSTICK_TYPE_DRUM_KIT"]
  const JOYSTICK_TYPE_DRUM_KIT u32

  @[link_name="SDL_JOYSTICK_TYPE_ARCADE_PAD"]
  const JOYSTICK_TYPE_ARCADE_PAD u32

  @[link_name="SDL_JOYSTICK_TYPE_THROTTLE"]
  const JOYSTICK_TYPE_THROTTLE u32

  @[link_name="SDL_JOYSTICK_TYPE_COUNT"]
  const JOYSTICK_TYPE_COUNT u32

  @[link_name="SDL_JOYSTICK_CONNECTION_INVALID"]
  const JOYSTICK_CONNECTION_INVALID u32

  @[link_name="SDL_JOYSTICK_CONNECTION_UNKNOWN"]
  const JOYSTICK_CONNECTION_UNKNOWN u32

  @[link_name="SDL_JOYSTICK_CONNECTION_WIRED"]
  const JOYSTICK_CONNECTION_WIRED u32

  @[link_name="SDL_JOYSTICK_CONNECTION_WIRELESS"]
  const JOYSTICK_CONNECTION_WIRELESS u32

  @[link_name="SDL_GAMEPAD_TYPE_UNKNOWN"]
  const GAMEPAD_TYPE_UNKNOWN u32

  @[link_name="SDL_GAMEPAD_TYPE_STANDARD"]
  const GAMEPAD_TYPE_STANDARD u32

  @[link_name="SDL_GAMEPAD_TYPE_XBOX360"]
  const GAMEPAD_TYPE_XBOX360 u32

  @[link_name="SDL_GAMEPAD_TYPE_XBOXONE"]
  const GAMEPAD_TYPE_XBOXONE u32

  @[link_name="SDL_GAMEPAD_TYPE_PS3"]
  const GAMEPAD_TYPE_PS3 u32

  @[link_name="SDL_GAMEPAD_TYPE_PS4"]
  const GAMEPAD_TYPE_PS4 u32

  @[link_name="SDL_GAMEPAD_TYPE_PS5"]
  const GAMEPAD_TYPE_PS5 u32

  @[link_name="SDL_GAMEPAD_TYPE_NINTENDO_SWITCH_PRO"]
  const GAMEPAD_TYPE_NINTENDO_SWITCH_PRO u32

  @[link_name="SDL_GAMEPAD_TYPE_NINTENDO_SWITCH_JOYCON_LEFT"]
  const GAMEPAD_TYPE_NINTENDO_SWITCH_JOYCON_LEFT u32

  @[link_name="SDL_GAMEPAD_TYPE_NINTENDO_SWITCH_JOYCON_RIGHT"]
  const GAMEPAD_TYPE_NINTENDO_SWITCH_JOYCON_RIGHT u32

  @[link_name="SDL_GAMEPAD_TYPE_NINTENDO_SWITCH_JOYCON_PAIR"]
  const GAMEPAD_TYPE_NINTENDO_SWITCH_JOYCON_PAIR u32

  @[link_name="SDL_GAMEPAD_TYPE_COUNT"]
  const GAMEPAD_TYPE_COUNT u32

  @[link_name="SDL_GAMEPAD_BUTTON_INVALID"]
  const GAMEPAD_BUTTON_INVALID u32

  @[link_name="SDL_GAMEPAD_BUTTON_SOUTH"]
  const GAMEPAD_BUTTON_SOUTH u32

  @[link_name="SDL_GAMEPAD_BUTTON_EAST"]
  const GAMEPAD_BUTTON_EAST u32

  @[link_name="SDL_GAMEPAD_BUTTON_WEST"]
  const GAMEPAD_BUTTON_WEST u32

  @[link_name="SDL_GAMEPAD_BUTTON_NORTH"]
  const GAMEPAD_BUTTON_NORTH u32

  @[link_name="SDL_GAMEPAD_BUTTON_BACK"]
  const GAMEPAD_BUTTON_BACK u32

  @[link_name="SDL_GAMEPAD_BUTTON_GUIDE"]
  const GAMEPAD_BUTTON_GUIDE u32

  @[link_name="SDL_GAMEPAD_BUTTON_START"]
  const GAMEPAD_BUTTON_START u32

  @[link_name="SDL_GAMEPAD_BUTTON_LEFT_STICK"]
  const GAMEPAD_BUTTON_LEFT_STICK u32

  @[link_name="SDL_GAMEPAD_BUTTON_RIGHT_STICK"]
  const GAMEPAD_BUTTON_RIGHT_STICK u32

  @[link_name="SDL_GAMEPAD_BUTTON_LEFT_SHOULDER"]
  const GAMEPAD_BUTTON_LEFT_SHOULDER u32

  @[link_name="SDL_GAMEPAD_BUTTON_RIGHT_SHOULDER"]
  const GAMEPAD_BUTTON_RIGHT_SHOULDER u32

  @[link_name="SDL_GAMEPAD_BUTTON_DPAD_UP"]
  const GAMEPAD_BUTTON_DPAD_UP u32

  @[link_name="SDL_GAMEPAD_BUTTON_DPAD_DOWN"]
  const GAMEPAD_BUTTON_DPAD_DOWN u32

  @[link_name="SDL_GAMEPAD_BUTTON_DPAD_LEFT"]
  const GAMEPAD_BUTTON_DPAD_LEFT u32

  @[link_name="SDL_GAMEPAD_BUTTON_DPAD_RIGHT"]
  const GAMEPAD_BUTTON_DPAD_RIGHT u32

  @[link_name="SDL_GAMEPAD_BUTTON_MISC1"]
  const GAMEPAD_BUTTON_MISC1 u32

  @[link_name="SDL_GAMEPAD_BUTTON_RIGHT_PADDLE1"]
  const GAMEPAD_BUTTON_RIGHT_PADDLE1 u32

  @[link_name="SDL_GAMEPAD_BUTTON_LEFT_PADDLE1"]
  const GAMEPAD_BUTTON_LEFT_PADDLE1 u32

  @[link_name="SDL_GAMEPAD_BUTTON_RIGHT_PADDLE2"]
  const GAMEPAD_BUTTON_RIGHT_PADDLE2 u32

  @[link_name="SDL_GAMEPAD_BUTTON_LEFT_PADDLE2"]
  const GAMEPAD_BUTTON_LEFT_PADDLE2 u32

  @[link_name="SDL_GAMEPAD_BUTTON_TOUCHPAD"]
  const GAMEPAD_BUTTON_TOUCHPAD u32

  @[link_name="SDL_GAMEPAD_BUTTON_MISC2"]
  const GAMEPAD_BUTTON_MISC2 u32

  @[link_name="SDL_GAMEPAD_BUTTON_MISC3"]
  const GAMEPAD_BUTTON_MISC3 u32

  @[link_name="SDL_GAMEPAD_BUTTON_MISC4"]
  const GAMEPAD_BUTTON_MISC4 u32

  @[link_name="SDL_GAMEPAD_BUTTON_MISC5"]
  const GAMEPAD_BUTTON_MISC5 u32

  @[link_name="SDL_GAMEPAD_BUTTON_MISC6"]
  const GAMEPAD_BUTTON_MISC6 u32

  @[link_name="SDL_GAMEPAD_BUTTON_COUNT"]
  const GAMEPAD_BUTTON_COUNT u32

  @[link_name="SDL_GAMEPAD_BUTTON_LABEL_UNKNOWN"]
  const GAMEPAD_BUTTON_LABEL_UNKNOWN u32

  @[link_name="SDL_GAMEPAD_BUTTON_LABEL_A"]
  const GAMEPAD_BUTTON_LABEL_A u32

  @[link_name="SDL_GAMEPAD_BUTTON_LABEL_B"]
  const GAMEPAD_BUTTON_LABEL_B u32

  @[link_name="SDL_GAMEPAD_BUTTON_LABEL_X"]
  const GAMEPAD_BUTTON_LABEL_X u32

  @[link_name="SDL_GAMEPAD_BUTTON_LABEL_Y"]
  const GAMEPAD_BUTTON_LABEL_Y u32

  @[link_name="SDL_GAMEPAD_BUTTON_LABEL_CROSS"]
  const GAMEPAD_BUTTON_LABEL_CROSS u32

  @[link_name="SDL_GAMEPAD_BUTTON_LABEL_CIRCLE"]
  const GAMEPAD_BUTTON_LABEL_CIRCLE u32

  @[link_name="SDL_GAMEPAD_BUTTON_LABEL_SQUARE"]
  const GAMEPAD_BUTTON_LABEL_SQUARE u32

  @[link_name="SDL_GAMEPAD_BUTTON_LABEL_TRIANGLE"]
  const GAMEPAD_BUTTON_LABEL_TRIANGLE u32

  @[link_name="SDL_GAMEPAD_AXIS_INVALID"]
  const GAMEPAD_AXIS_INVALID u32

  @[link_name="SDL_GAMEPAD_AXIS_LEFTX"]
  const GAMEPAD_AXIS_LEFTX u32

  @[link_name="SDL_GAMEPAD_AXIS_LEFTY"]
  const GAMEPAD_AXIS_LEFTY u32

  @[link_name="SDL_GAMEPAD_AXIS_RIGHTX"]
  const GAMEPAD_AXIS_RIGHTX u32

  @[link_name="SDL_GAMEPAD_AXIS_RIGHTY"]
  const GAMEPAD_AXIS_RIGHTY u32

  @[link_name="SDL_GAMEPAD_AXIS_LEFT_TRIGGER"]
  const GAMEPAD_AXIS_LEFT_TRIGGER u32

  @[link_name="SDL_GAMEPAD_AXIS_RIGHT_TRIGGER"]
  const GAMEPAD_AXIS_RIGHT_TRIGGER u32

  @[link_name="SDL_GAMEPAD_AXIS_COUNT"]
  const GAMEPAD_AXIS_COUNT u32

  @[link_name="SDL_GAMEPAD_BINDTYPE_NONE"]
  const GAMEPAD_BINDTYPE_NONE u32

  @[link_name="SDL_GAMEPAD_BINDTYPE_BUTTON"]
  const GAMEPAD_BINDTYPE_BUTTON u32

  @[link_name="SDL_GAMEPAD_BINDTYPE_AXIS"]
  const GAMEPAD_BINDTYPE_AXIS u32

  @[link_name="SDL_GAMEPAD_BINDTYPE_HAT"]
  const GAMEPAD_BINDTYPE_HAT u32

  @[link_name="SDL_SCANCODE_UNKNOWN"]
  const SCANCODE_UNKNOWN u32

  @[link_name="SDL_SCANCODE_A"]
  const SCANCODE_A u32

  @[link_name="SDL_SCANCODE_B"]
  const SCANCODE_B u32

  @[link_name="SDL_SCANCODE_C"]
  const SCANCODE_C u32

  @[link_name="SDL_SCANCODE_D"]
  const SCANCODE_D u32

  @[link_name="SDL_SCANCODE_E"]
  const SCANCODE_E u32

  @[link_name="SDL_SCANCODE_F"]
  const SCANCODE_F u32

  @[link_name="SDL_SCANCODE_G"]
  const SCANCODE_G u32

  @[link_name="SDL_SCANCODE_H"]
  const SCANCODE_H u32

  @[link_name="SDL_SCANCODE_I"]
  const SCANCODE_I u32

  @[link_name="SDL_SCANCODE_J"]
  const SCANCODE_J u32

  @[link_name="SDL_SCANCODE_K"]
  const SCANCODE_K u32

  @[link_name="SDL_SCANCODE_L"]
  const SCANCODE_L u32

  @[link_name="SDL_SCANCODE_M"]
  const SCANCODE_M u32

  @[link_name="SDL_SCANCODE_N"]
  const SCANCODE_N u32

  @[link_name="SDL_SCANCODE_O"]
  const SCANCODE_O u32

  @[link_name="SDL_SCANCODE_P"]
  const SCANCODE_P u32

  @[link_name="SDL_SCANCODE_Q"]
  const SCANCODE_Q u32

  @[link_name="SDL_SCANCODE_R"]
  const SCANCODE_R u32

  @[link_name="SDL_SCANCODE_S"]
  const SCANCODE_S u32

  @[link_name="SDL_SCANCODE_T"]
  const SCANCODE_T u32

  @[link_name="SDL_SCANCODE_U"]
  const SCANCODE_U u32

  @[link_name="SDL_SCANCODE_V"]
  const SCANCODE_V u32

  @[link_name="SDL_SCANCODE_W"]
  const SCANCODE_W u32

  @[link_name="SDL_SCANCODE_X"]
  const SCANCODE_X u32

  @[link_name="SDL_SCANCODE_Y"]
  const SCANCODE_Y u32

  @[link_name="SDL_SCANCODE_Z"]
  const SCANCODE_Z u32

  @[link_name="SDL_SCANCODE_1"]
  const SCANCODE_1 u32

  @[link_name="SDL_SCANCODE_2"]
  const SCANCODE_2 u32

  @[link_name="SDL_SCANCODE_3"]
  const SCANCODE_3 u32

  @[link_name="SDL_SCANCODE_4"]
  const SCANCODE_4 u32

  @[link_name="SDL_SCANCODE_5"]
  const SCANCODE_5 u32

  @[link_name="SDL_SCANCODE_6"]
  const SCANCODE_6 u32

  @[link_name="SDL_SCANCODE_7"]
  const SCANCODE_7 u32

  @[link_name="SDL_SCANCODE_8"]
  const SCANCODE_8 u32

  @[link_name="SDL_SCANCODE_9"]
  const SCANCODE_9 u32

  @[link_name="SDL_SCANCODE_0"]
  const SCANCODE_0 u32

  @[link_name="SDL_SCANCODE_RETURN"]
  const SCANCODE_RETURN u32

  @[link_name="SDL_SCANCODE_ESCAPE"]
  const SCANCODE_ESCAPE u32

  @[link_name="SDL_SCANCODE_BACKSPACE"]
  const SCANCODE_BACKSPACE u32

  @[link_name="SDL_SCANCODE_TAB"]
  const SCANCODE_TAB u32

  @[link_name="SDL_SCANCODE_SPACE"]
  const SCANCODE_SPACE u32

  @[link_name="SDL_SCANCODE_MINUS"]
  const SCANCODE_MINUS u32

  @[link_name="SDL_SCANCODE_EQUALS"]
  const SCANCODE_EQUALS u32

  @[link_name="SDL_SCANCODE_LEFTBRACKET"]
  const SCANCODE_LEFTBRACKET u32

  @[link_name="SDL_SCANCODE_RIGHTBRACKET"]
  const SCANCODE_RIGHTBRACKET u32

  @[link_name="SDL_SCANCODE_BACKSLASH"]
  const SCANCODE_BACKSLASH u32

  @[link_name="SDL_SCANCODE_NONUSHASH"]
  const SCANCODE_NONUSHASH u32

  @[link_name="SDL_SCANCODE_SEMICOLON"]
  const SCANCODE_SEMICOLON u32

  @[link_name="SDL_SCANCODE_APOSTROPHE"]
  const SCANCODE_APOSTROPHE u32

  @[link_name="SDL_SCANCODE_GRAVE"]
  const SCANCODE_GRAVE u32

  @[link_name="SDL_SCANCODE_COMMA"]
  const SCANCODE_COMMA u32

  @[link_name="SDL_SCANCODE_PERIOD"]
  const SCANCODE_PERIOD u32

  @[link_name="SDL_SCANCODE_SLASH"]
  const SCANCODE_SLASH u32

  @[link_name="SDL_SCANCODE_CAPSLOCK"]
  const SCANCODE_CAPSLOCK u32

  @[link_name="SDL_SCANCODE_F1"]
  const SCANCODE_F1 u32

  @[link_name="SDL_SCANCODE_F2"]
  const SCANCODE_F2 u32

  @[link_name="SDL_SCANCODE_F3"]
  const SCANCODE_F3 u32

  @[link_name="SDL_SCANCODE_F4"]
  const SCANCODE_F4 u32

  @[link_name="SDL_SCANCODE_F5"]
  const SCANCODE_F5 u32

  @[link_name="SDL_SCANCODE_F6"]
  const SCANCODE_F6 u32

  @[link_name="SDL_SCANCODE_F7"]
  const SCANCODE_F7 u32

  @[link_name="SDL_SCANCODE_F8"]
  const SCANCODE_F8 u32

  @[link_name="SDL_SCANCODE_F9"]
  const SCANCODE_F9 u32

  @[link_name="SDL_SCANCODE_F10"]
  const SCANCODE_F10 u32

  @[link_name="SDL_SCANCODE_F11"]
  const SCANCODE_F11 u32

  @[link_name="SDL_SCANCODE_F12"]
  const SCANCODE_F12 u32

  @[link_name="SDL_SCANCODE_PRINTSCREEN"]
  const SCANCODE_PRINTSCREEN u32

  @[link_name="SDL_SCANCODE_SCROLLLOCK"]
  const SCANCODE_SCROLLLOCK u32

  @[link_name="SDL_SCANCODE_PAUSE"]
  const SCANCODE_PAUSE u32

  @[link_name="SDL_SCANCODE_INSERT"]
  const SCANCODE_INSERT u32

  @[link_name="SDL_SCANCODE_HOME"]
  const SCANCODE_HOME u32

  @[link_name="SDL_SCANCODE_PAGEUP"]
  const SCANCODE_PAGEUP u32

  @[link_name="SDL_SCANCODE_DELETE"]
  const SCANCODE_DELETE u32

  @[link_name="SDL_SCANCODE_END"]
  const SCANCODE_END u32

  @[link_name="SDL_SCANCODE_PAGEDOWN"]
  const SCANCODE_PAGEDOWN u32

  @[link_name="SDL_SCANCODE_RIGHT"]
  const SCANCODE_RIGHT u32

  @[link_name="SDL_SCANCODE_LEFT"]
  const SCANCODE_LEFT u32

  @[link_name="SDL_SCANCODE_DOWN"]
  const SCANCODE_DOWN u32

  @[link_name="SDL_SCANCODE_UP"]
  const SCANCODE_UP u32

  @[link_name="SDL_SCANCODE_NUMLOCKCLEAR"]
  const SCANCODE_NUMLOCKCLEAR u32

  @[link_name="SDL_SCANCODE_KP_DIVIDE"]
  const SCANCODE_KP_DIVIDE u32

  @[link_name="SDL_SCANCODE_KP_MULTIPLY"]
  const SCANCODE_KP_MULTIPLY u32

  @[link_name="SDL_SCANCODE_KP_MINUS"]
  const SCANCODE_KP_MINUS u32

  @[link_name="SDL_SCANCODE_KP_PLUS"]
  const SCANCODE_KP_PLUS u32

  @[link_name="SDL_SCANCODE_KP_ENTER"]
  const SCANCODE_KP_ENTER u32

  @[link_name="SDL_SCANCODE_KP_1"]
  const SCANCODE_KP_1 u32

  @[link_name="SDL_SCANCODE_KP_2"]
  const SCANCODE_KP_2 u32

  @[link_name="SDL_SCANCODE_KP_3"]
  const SCANCODE_KP_3 u32

  @[link_name="SDL_SCANCODE_KP_4"]
  const SCANCODE_KP_4 u32

  @[link_name="SDL_SCANCODE_KP_5"]
  const SCANCODE_KP_5 u32

  @[link_name="SDL_SCANCODE_KP_6"]
  const SCANCODE_KP_6 u32

  @[link_name="SDL_SCANCODE_KP_7"]
  const SCANCODE_KP_7 u32

  @[link_name="SDL_SCANCODE_KP_8"]
  const SCANCODE_KP_8 u32

  @[link_name="SDL_SCANCODE_KP_9"]
  const SCANCODE_KP_9 u32

  @[link_name="SDL_SCANCODE_KP_0"]
  const SCANCODE_KP_0 u32

  @[link_name="SDL_SCANCODE_KP_PERIOD"]
  const SCANCODE_KP_PERIOD u32

  @[link_name="SDL_SCANCODE_NONUSBACKSLASH"]
  const SCANCODE_NONUSBACKSLASH u32

  @[link_name="SDL_SCANCODE_APPLICATION"]
  const SCANCODE_APPLICATION u32

  @[link_name="SDL_SCANCODE_POWER"]
  const SCANCODE_POWER u32

  @[link_name="SDL_SCANCODE_KP_EQUALS"]
  const SCANCODE_KP_EQUALS u32

  @[link_name="SDL_SCANCODE_F13"]
  const SCANCODE_F13 u32

  @[link_name="SDL_SCANCODE_F14"]
  const SCANCODE_F14 u32

  @[link_name="SDL_SCANCODE_F15"]
  const SCANCODE_F15 u32

  @[link_name="SDL_SCANCODE_F16"]
  const SCANCODE_F16 u32

  @[link_name="SDL_SCANCODE_F17"]
  const SCANCODE_F17 u32

  @[link_name="SDL_SCANCODE_F18"]
  const SCANCODE_F18 u32

  @[link_name="SDL_SCANCODE_F19"]
  const SCANCODE_F19 u32

  @[link_name="SDL_SCANCODE_F20"]
  const SCANCODE_F20 u32

  @[link_name="SDL_SCANCODE_F21"]
  const SCANCODE_F21 u32

  @[link_name="SDL_SCANCODE_F22"]
  const SCANCODE_F22 u32

  @[link_name="SDL_SCANCODE_F23"]
  const SCANCODE_F23 u32

  @[link_name="SDL_SCANCODE_F24"]
  const SCANCODE_F24 u32

  @[link_name="SDL_SCANCODE_EXECUTE"]
  const SCANCODE_EXECUTE u32

  @[link_name="SDL_SCANCODE_HELP"]
  const SCANCODE_HELP u32

  @[link_name="SDL_SCANCODE_MENU"]
  const SCANCODE_MENU u32

  @[link_name="SDL_SCANCODE_SELECT"]
  const SCANCODE_SELECT u32

  @[link_name="SDL_SCANCODE_STOP"]
  const SCANCODE_STOP u32

  @[link_name="SDL_SCANCODE_AGAIN"]
  const SCANCODE_AGAIN u32

  @[link_name="SDL_SCANCODE_UNDO"]
  const SCANCODE_UNDO u32

  @[link_name="SDL_SCANCODE_CUT"]
  const SCANCODE_CUT u32

  @[link_name="SDL_SCANCODE_COPY"]
  const SCANCODE_COPY u32

  @[link_name="SDL_SCANCODE_PASTE"]
  const SCANCODE_PASTE u32

  @[link_name="SDL_SCANCODE_FIND"]
  const SCANCODE_FIND u32

  @[link_name="SDL_SCANCODE_MUTE"]
  const SCANCODE_MUTE u32

  @[link_name="SDL_SCANCODE_VOLUMEUP"]
  const SCANCODE_VOLUMEUP u32

  @[link_name="SDL_SCANCODE_VOLUMEDOWN"]
  const SCANCODE_VOLUMEDOWN u32

  @[link_name="SDL_SCANCODE_KP_COMMA"]
  const SCANCODE_KP_COMMA u32

  @[link_name="SDL_SCANCODE_KP_EQUALSAS400"]
  const SCANCODE_KP_EQUALSAS400 u32

  @[link_name="SDL_SCANCODE_INTERNATIONAL1"]
  const SCANCODE_INTERNATIONAL1 u32

  @[link_name="SDL_SCANCODE_INTERNATIONAL2"]
  const SCANCODE_INTERNATIONAL2 u32

  @[link_name="SDL_SCANCODE_INTERNATIONAL3"]
  const SCANCODE_INTERNATIONAL3 u32

  @[link_name="SDL_SCANCODE_INTERNATIONAL4"]
  const SCANCODE_INTERNATIONAL4 u32

  @[link_name="SDL_SCANCODE_INTERNATIONAL5"]
  const SCANCODE_INTERNATIONAL5 u32

  @[link_name="SDL_SCANCODE_INTERNATIONAL6"]
  const SCANCODE_INTERNATIONAL6 u32

  @[link_name="SDL_SCANCODE_INTERNATIONAL7"]
  const SCANCODE_INTERNATIONAL7 u32

  @[link_name="SDL_SCANCODE_INTERNATIONAL8"]
  const SCANCODE_INTERNATIONAL8 u32

  @[link_name="SDL_SCANCODE_INTERNATIONAL9"]
  const SCANCODE_INTERNATIONAL9 u32

  @[link_name="SDL_SCANCODE_LANG1"]
  const SCANCODE_LANG1 u32

  @[link_name="SDL_SCANCODE_LANG2"]
  const SCANCODE_LANG2 u32

  @[link_name="SDL_SCANCODE_LANG3"]
  const SCANCODE_LANG3 u32

  @[link_name="SDL_SCANCODE_LANG4"]
  const SCANCODE_LANG4 u32

  @[link_name="SDL_SCANCODE_LANG5"]
  const SCANCODE_LANG5 u32

  @[link_name="SDL_SCANCODE_LANG6"]
  const SCANCODE_LANG6 u32

  @[link_name="SDL_SCANCODE_LANG7"]
  const SCANCODE_LANG7 u32

  @[link_name="SDL_SCANCODE_LANG8"]
  const SCANCODE_LANG8 u32

  @[link_name="SDL_SCANCODE_LANG9"]
  const SCANCODE_LANG9 u32

  @[link_name="SDL_SCANCODE_ALTERASE"]
  const SCANCODE_ALTERASE u32

  @[link_name="SDL_SCANCODE_SYSREQ"]
  const SCANCODE_SYSREQ u32

  @[link_name="SDL_SCANCODE_CANCEL"]
  const SCANCODE_CANCEL u32

  @[link_name="SDL_SCANCODE_CLEAR"]
  const SCANCODE_CLEAR u32

  @[link_name="SDL_SCANCODE_PRIOR"]
  const SCANCODE_PRIOR u32

  @[link_name="SDL_SCANCODE_RETURN2"]
  const SCANCODE_RETURN2 u32

  @[link_name="SDL_SCANCODE_SEPARATOR"]
  const SCANCODE_SEPARATOR u32

  @[link_name="SDL_SCANCODE_OUT"]
  const SCANCODE_OUT u32

  @[link_name="SDL_SCANCODE_OPER"]
  const SCANCODE_OPER u32

  @[link_name="SDL_SCANCODE_CLEARAGAIN"]
  const SCANCODE_CLEARAGAIN u32

  @[link_name="SDL_SCANCODE_CRSEL"]
  const SCANCODE_CRSEL u32

  @[link_name="SDL_SCANCODE_EXSEL"]
  const SCANCODE_EXSEL u32

  @[link_name="SDL_SCANCODE_KP_00"]
  const SCANCODE_KP_00 u32

  @[link_name="SDL_SCANCODE_KP_000"]
  const SCANCODE_KP_000 u32

  @[link_name="SDL_SCANCODE_THOUSANDSSEPARATOR"]
  const SCANCODE_THOUSANDSSEPARATOR u32

  @[link_name="SDL_SCANCODE_DECIMALSEPARATOR"]
  const SCANCODE_DECIMALSEPARATOR u32

  @[link_name="SDL_SCANCODE_CURRENCYUNIT"]
  const SCANCODE_CURRENCYUNIT u32

  @[link_name="SDL_SCANCODE_CURRENCYSUBUNIT"]
  const SCANCODE_CURRENCYSUBUNIT u32

  @[link_name="SDL_SCANCODE_KP_LEFTPAREN"]
  const SCANCODE_KP_LEFTPAREN u32

  @[link_name="SDL_SCANCODE_KP_RIGHTPAREN"]
  const SCANCODE_KP_RIGHTPAREN u32

  @[link_name="SDL_SCANCODE_KP_LEFTBRACE"]
  const SCANCODE_KP_LEFTBRACE u32

  @[link_name="SDL_SCANCODE_KP_RIGHTBRACE"]
  const SCANCODE_KP_RIGHTBRACE u32

  @[link_name="SDL_SCANCODE_KP_TAB"]
  const SCANCODE_KP_TAB u32

  @[link_name="SDL_SCANCODE_KP_BACKSPACE"]
  const SCANCODE_KP_BACKSPACE u32

  @[link_name="SDL_SCANCODE_KP_A"]
  const SCANCODE_KP_A u32

  @[link_name="SDL_SCANCODE_KP_B"]
  const SCANCODE_KP_B u32

  @[link_name="SDL_SCANCODE_KP_C"]
  const SCANCODE_KP_C u32

  @[link_name="SDL_SCANCODE_KP_D"]
  const SCANCODE_KP_D u32

  @[link_name="SDL_SCANCODE_KP_E"]
  const SCANCODE_KP_E u32

  @[link_name="SDL_SCANCODE_KP_F"]
  const SCANCODE_KP_F u32

  @[link_name="SDL_SCANCODE_KP_XOR"]
  const SCANCODE_KP_XOR u32

  @[link_name="SDL_SCANCODE_KP_POWER"]
  const SCANCODE_KP_POWER u32

  @[link_name="SDL_SCANCODE_KP_PERCENT"]
  const SCANCODE_KP_PERCENT u32

  @[link_name="SDL_SCANCODE_KP_LESS"]
  const SCANCODE_KP_LESS u32

  @[link_name="SDL_SCANCODE_KP_GREATER"]
  const SCANCODE_KP_GREATER u32

  @[link_name="SDL_SCANCODE_KP_AMPERSAND"]
  const SCANCODE_KP_AMPERSAND u32

  @[link_name="SDL_SCANCODE_KP_DBLAMPERSAND"]
  const SCANCODE_KP_DBLAMPERSAND u32

  @[link_name="SDL_SCANCODE_KP_VERTICALBAR"]
  const SCANCODE_KP_VERTICALBAR u32

  @[link_name="SDL_SCANCODE_KP_DBLVERTICALBAR"]
  const SCANCODE_KP_DBLVERTICALBAR u32

  @[link_name="SDL_SCANCODE_KP_COLON"]
  const SCANCODE_KP_COLON u32

  @[link_name="SDL_SCANCODE_KP_HASH"]
  const SCANCODE_KP_HASH u32

  @[link_name="SDL_SCANCODE_KP_SPACE"]
  const SCANCODE_KP_SPACE u32

  @[link_name="SDL_SCANCODE_KP_AT"]
  const SCANCODE_KP_AT u32

  @[link_name="SDL_SCANCODE_KP_EXCLAM"]
  const SCANCODE_KP_EXCLAM u32

  @[link_name="SDL_SCANCODE_KP_MEMSTORE"]
  const SCANCODE_KP_MEMSTORE u32

  @[link_name="SDL_SCANCODE_KP_MEMRECALL"]
  const SCANCODE_KP_MEMRECALL u32

  @[link_name="SDL_SCANCODE_KP_MEMCLEAR"]
  const SCANCODE_KP_MEMCLEAR u32

  @[link_name="SDL_SCANCODE_KP_MEMADD"]
  const SCANCODE_KP_MEMADD u32

  @[link_name="SDL_SCANCODE_KP_MEMSUBTRACT"]
  const SCANCODE_KP_MEMSUBTRACT u32

  @[link_name="SDL_SCANCODE_KP_MEMMULTIPLY"]
  const SCANCODE_KP_MEMMULTIPLY u32

  @[link_name="SDL_SCANCODE_KP_MEMDIVIDE"]
  const SCANCODE_KP_MEMDIVIDE u32

  @[link_name="SDL_SCANCODE_KP_PLUSMINUS"]
  const SCANCODE_KP_PLUSMINUS u32

  @[link_name="SDL_SCANCODE_KP_CLEAR"]
  const SCANCODE_KP_CLEAR u32

  @[link_name="SDL_SCANCODE_KP_CLEARENTRY"]
  const SCANCODE_KP_CLEARENTRY u32

  @[link_name="SDL_SCANCODE_KP_BINARY"]
  const SCANCODE_KP_BINARY u32

  @[link_name="SDL_SCANCODE_KP_OCTAL"]
  const SCANCODE_KP_OCTAL u32

  @[link_name="SDL_SCANCODE_KP_DECIMAL"]
  const SCANCODE_KP_DECIMAL u32

  @[link_name="SDL_SCANCODE_KP_HEXADECIMAL"]
  const SCANCODE_KP_HEXADECIMAL u32

  @[link_name="SDL_SCANCODE_LCTRL"]
  const SCANCODE_LCTRL u32

  @[link_name="SDL_SCANCODE_LSHIFT"]
  const SCANCODE_LSHIFT u32

  @[link_name="SDL_SCANCODE_LALT"]
  const SCANCODE_LALT u32

  @[link_name="SDL_SCANCODE_LGUI"]
  const SCANCODE_LGUI u32

  @[link_name="SDL_SCANCODE_RCTRL"]
  const SCANCODE_RCTRL u32

  @[link_name="SDL_SCANCODE_RSHIFT"]
  const SCANCODE_RSHIFT u32

  @[link_name="SDL_SCANCODE_RALT"]
  const SCANCODE_RALT u32

  @[link_name="SDL_SCANCODE_RGUI"]
  const SCANCODE_RGUI u32

  @[link_name="SDL_SCANCODE_MODE"]
  const SCANCODE_MODE u32

  @[link_name="SDL_SCANCODE_SLEEP"]
  const SCANCODE_SLEEP u32

  @[link_name="SDL_SCANCODE_WAKE"]
  const SCANCODE_WAKE u32

  @[link_name="SDL_SCANCODE_CHANNEL_INCREMENT"]
  const SCANCODE_CHANNEL_INCREMENT u32

  @[link_name="SDL_SCANCODE_CHANNEL_DECREMENT"]
  const SCANCODE_CHANNEL_DECREMENT u32

  @[link_name="SDL_SCANCODE_MEDIA_PLAY"]
  const SCANCODE_MEDIA_PLAY u32

  @[link_name="SDL_SCANCODE_MEDIA_PAUSE"]
  const SCANCODE_MEDIA_PAUSE u32

  @[link_name="SDL_SCANCODE_MEDIA_RECORD"]
  const SCANCODE_MEDIA_RECORD u32

  @[link_name="SDL_SCANCODE_MEDIA_FAST_FORWARD"]
  const SCANCODE_MEDIA_FAST_FORWARD u32

  @[link_name="SDL_SCANCODE_MEDIA_REWIND"]
  const SCANCODE_MEDIA_REWIND u32

  @[link_name="SDL_SCANCODE_MEDIA_NEXT_TRACK"]
  const SCANCODE_MEDIA_NEXT_TRACK u32

  @[link_name="SDL_SCANCODE_MEDIA_PREVIOUS_TRACK"]
  const SCANCODE_MEDIA_PREVIOUS_TRACK u32

  @[link_name="SDL_SCANCODE_MEDIA_STOP"]
  const SCANCODE_MEDIA_STOP u32

  @[link_name="SDL_SCANCODE_MEDIA_EJECT"]
  const SCANCODE_MEDIA_EJECT u32

  @[link_name="SDL_SCANCODE_MEDIA_PLAY_PAUSE"]
  const SCANCODE_MEDIA_PLAY_PAUSE u32

  @[link_name="SDL_SCANCODE_MEDIA_SELECT"]
  const SCANCODE_MEDIA_SELECT u32

  @[link_name="SDL_SCANCODE_AC_NEW"]
  const SCANCODE_AC_NEW u32

  @[link_name="SDL_SCANCODE_AC_OPEN"]
  const SCANCODE_AC_OPEN u32

  @[link_name="SDL_SCANCODE_AC_CLOSE"]
  const SCANCODE_AC_CLOSE u32

  @[link_name="SDL_SCANCODE_AC_EXIT"]
  const SCANCODE_AC_EXIT u32

  @[link_name="SDL_SCANCODE_AC_SAVE"]
  const SCANCODE_AC_SAVE u32

  @[link_name="SDL_SCANCODE_AC_PRINT"]
  const SCANCODE_AC_PRINT u32

  @[link_name="SDL_SCANCODE_AC_PROPERTIES"]
  const SCANCODE_AC_PROPERTIES u32

  @[link_name="SDL_SCANCODE_AC_SEARCH"]
  const SCANCODE_AC_SEARCH u32

  @[link_name="SDL_SCANCODE_AC_HOME"]
  const SCANCODE_AC_HOME u32

  @[link_name="SDL_SCANCODE_AC_BACK"]
  const SCANCODE_AC_BACK u32

  @[link_name="SDL_SCANCODE_AC_FORWARD"]
  const SCANCODE_AC_FORWARD u32

  @[link_name="SDL_SCANCODE_AC_STOP"]
  const SCANCODE_AC_STOP u32

  @[link_name="SDL_SCANCODE_AC_REFRESH"]
  const SCANCODE_AC_REFRESH u32

  @[link_name="SDL_SCANCODE_AC_BOOKMARKS"]
  const SCANCODE_AC_BOOKMARKS u32

  @[link_name="SDL_SCANCODE_SOFTLEFT"]
  const SCANCODE_SOFTLEFT u32

  @[link_name="SDL_SCANCODE_SOFTRIGHT"]
  const SCANCODE_SOFTRIGHT u32

  @[link_name="SDL_SCANCODE_CALL"]
  const SCANCODE_CALL u32

  @[link_name="SDL_SCANCODE_ENDCALL"]
  const SCANCODE_ENDCALL u32

  @[link_name="SDL_SCANCODE_RESERVED"]
  const SCANCODE_RESERVED u32

  @[link_name="SDL_SCANCODE_COUNT"]
  const SCANCODE_COUNT u32

  @[link_name="SDL_TEXTINPUT_TYPE_TEXT"]
  const TEXTINPUT_TYPE_TEXT u32

  @[link_name="SDL_TEXTINPUT_TYPE_TEXT_NAME"]
  const TEXTINPUT_TYPE_TEXT_NAME u32

  @[link_name="SDL_TEXTINPUT_TYPE_TEXT_EMAIL"]
  const TEXTINPUT_TYPE_TEXT_EMAIL u32

  @[link_name="SDL_TEXTINPUT_TYPE_TEXT_USERNAME"]
  const TEXTINPUT_TYPE_TEXT_USERNAME u32

  @[link_name="SDL_TEXTINPUT_TYPE_TEXT_PASSWORD_HIDDEN"]
  const TEXTINPUT_TYPE_TEXT_PASSWORD_HIDDEN u32

  @[link_name="SDL_TEXTINPUT_TYPE_TEXT_PASSWORD_VISIBLE"]
  const TEXTINPUT_TYPE_TEXT_PASSWORD_VISIBLE u32

  @[link_name="SDL_TEXTINPUT_TYPE_NUMBER"]
  const TEXTINPUT_TYPE_NUMBER u32

  @[link_name="SDL_TEXTINPUT_TYPE_NUMBER_PASSWORD_HIDDEN"]
  const TEXTINPUT_TYPE_NUMBER_PASSWORD_HIDDEN u32

  @[link_name="SDL_TEXTINPUT_TYPE_NUMBER_PASSWORD_VISIBLE"]
  const TEXTINPUT_TYPE_NUMBER_PASSWORD_VISIBLE u32

  @[link_name="SDL_CAPITALIZE_NONE"]
  const CAPITALIZE_NONE u32

  @[link_name="SDL_CAPITALIZE_SENTENCES"]
  const CAPITALIZE_SENTENCES u32

  @[link_name="SDL_CAPITALIZE_WORDS"]
  const CAPITALIZE_WORDS u32

  @[link_name="SDL_CAPITALIZE_LETTERS"]
  const CAPITALIZE_LETTERS u32

  @[link_name="SDL_SYSTEM_CURSOR_DEFAULT"]
  const SYSTEM_CURSOR_DEFAULT u32

  @[link_name="SDL_SYSTEM_CURSOR_TEXT"]
  const SYSTEM_CURSOR_TEXT u32

  @[link_name="SDL_SYSTEM_CURSOR_WAIT"]
  const SYSTEM_CURSOR_WAIT u32

  @[link_name="SDL_SYSTEM_CURSOR_CROSSHAIR"]
  const SYSTEM_CURSOR_CROSSHAIR u32

  @[link_name="SDL_SYSTEM_CURSOR_PROGRESS"]
  const SYSTEM_CURSOR_PROGRESS u32

  @[link_name="SDL_SYSTEM_CURSOR_NWSE_RESIZE"]
  const SYSTEM_CURSOR_NWSE_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_NESW_RESIZE"]
  const SYSTEM_CURSOR_NESW_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_EW_RESIZE"]
  const SYSTEM_CURSOR_EW_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_NS_RESIZE"]
  const SYSTEM_CURSOR_NS_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_MOVE"]
  const SYSTEM_CURSOR_MOVE u32

  @[link_name="SDL_SYSTEM_CURSOR_NOT_ALLOWED"]
  const SYSTEM_CURSOR_NOT_ALLOWED u32

  @[link_name="SDL_SYSTEM_CURSOR_POINTER"]
  const SYSTEM_CURSOR_POINTER u32

  @[link_name="SDL_SYSTEM_CURSOR_NW_RESIZE"]
  const SYSTEM_CURSOR_NW_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_N_RESIZE"]
  const SYSTEM_CURSOR_N_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_NE_RESIZE"]
  const SYSTEM_CURSOR_NE_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_E_RESIZE"]
  const SYSTEM_CURSOR_E_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_SE_RESIZE"]
  const SYSTEM_CURSOR_SE_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_S_RESIZE"]
  const SYSTEM_CURSOR_S_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_SW_RESIZE"]
  const SYSTEM_CURSOR_SW_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_W_RESIZE"]
  const SYSTEM_CURSOR_W_RESIZE u32

  @[link_name="SDL_SYSTEM_CURSOR_COUNT"]
  const SYSTEM_CURSOR_COUNT u32

  @[link_name="SDL_MOUSEWHEEL_NORMAL"]
  const MOUSEWHEEL_NORMAL u32

  @[link_name="SDL_MOUSEWHEEL_FLIPPED"]
  const MOUSEWHEEL_FLIPPED u32

  @[link_name="SDL_TOUCH_DEVICE_INVALID"]
  const TOUCH_DEVICE_INVALID u32

  @[link_name="SDL_TOUCH_DEVICE_DIRECT"]
  const TOUCH_DEVICE_DIRECT u32

  @[link_name="SDL_TOUCH_DEVICE_INDIRECT_ABSOLUTE"]
  const TOUCH_DEVICE_INDIRECT_ABSOLUTE u32

  @[link_name="SDL_TOUCH_DEVICE_INDIRECT_RELATIVE"]
  const TOUCH_DEVICE_INDIRECT_RELATIVE u32

  @[link_name="SDL_PEN_AXIS_PRESSURE"]
  const PEN_AXIS_PRESSURE u32

  @[link_name="SDL_PEN_AXIS_XTILT"]
  const PEN_AXIS_XTILT u32

  @[link_name="SDL_PEN_AXIS_YTILT"]
  const PEN_AXIS_YTILT u32

  @[link_name="SDL_PEN_AXIS_DISTANCE"]
  const PEN_AXIS_DISTANCE u32

  @[link_name="SDL_PEN_AXIS_ROTATION"]
  const PEN_AXIS_ROTATION u32

  @[link_name="SDL_PEN_AXIS_SLIDER"]
  const PEN_AXIS_SLIDER u32

  @[link_name="SDL_PEN_AXIS_TANGENTIAL_PRESSURE"]
  const PEN_AXIS_TANGENTIAL_PRESSURE u32

  @[link_name="SDL_PEN_AXIS_COUNT"]
  const PEN_AXIS_COUNT u32

  @[link_name="SDL_EVENT_FIRST"]
  const EVENT_FIRST u32

  @[link_name="SDL_EVENT_QUIT"]
  const EVENT_QUIT u32

  @[link_name="SDL_EVENT_TERMINATING"]
  const EVENT_TERMINATING u32

  @[link_name="SDL_EVENT_LOW_MEMORY"]
  const EVENT_LOW_MEMORY u32

  @[link_name="SDL_EVENT_WILL_ENTER_BACKGROUND"]
  const EVENT_WILL_ENTER_BACKGROUND u32

  @[link_name="SDL_EVENT_DID_ENTER_BACKGROUND"]
  const EVENT_DID_ENTER_BACKGROUND u32

  @[link_name="SDL_EVENT_WILL_ENTER_FOREGROUND"]
  const EVENT_WILL_ENTER_FOREGROUND u32

  @[link_name="SDL_EVENT_DID_ENTER_FOREGROUND"]
  const EVENT_DID_ENTER_FOREGROUND u32

  @[link_name="SDL_EVENT_LOCALE_CHANGED"]
  const EVENT_LOCALE_CHANGED u32

  @[link_name="SDL_EVENT_SYSTEM_THEME_CHANGED"]
  const EVENT_SYSTEM_THEME_CHANGED u32

  @[link_name="SDL_EVENT_DISPLAY_ORIENTATION"]
  const EVENT_DISPLAY_ORIENTATION u32

  @[link_name="SDL_EVENT_DISPLAY_ADDED"]
  const EVENT_DISPLAY_ADDED u32

  @[link_name="SDL_EVENT_DISPLAY_REMOVED"]
  const EVENT_DISPLAY_REMOVED u32

  @[link_name="SDL_EVENT_DISPLAY_MOVED"]
  const EVENT_DISPLAY_MOVED u32

  @[link_name="SDL_EVENT_DISPLAY_DESKTOP_MODE_CHANGED"]
  const EVENT_DISPLAY_DESKTOP_MODE_CHANGED u32

  @[link_name="SDL_EVENT_DISPLAY_CURRENT_MODE_CHANGED"]
  const EVENT_DISPLAY_CURRENT_MODE_CHANGED u32

  @[link_name="SDL_EVENT_DISPLAY_CONTENT_SCALE_CHANGED"]
  const EVENT_DISPLAY_CONTENT_SCALE_CHANGED u32

  @[link_name="SDL_EVENT_DISPLAY_FIRST"]
  const EVENT_DISPLAY_FIRST u32

  @[link_name="SDL_EVENT_DISPLAY_LAST"]
  const EVENT_DISPLAY_LAST u32

  @[link_name="SDL_EVENT_WINDOW_SHOWN"]
  const EVENT_WINDOW_SHOWN u32

  @[link_name="SDL_EVENT_WINDOW_HIDDEN"]
  const EVENT_WINDOW_HIDDEN u32

  @[link_name="SDL_EVENT_WINDOW_EXPOSED"]
  const EVENT_WINDOW_EXPOSED u32

  @[link_name="SDL_EVENT_WINDOW_MOVED"]
  const EVENT_WINDOW_MOVED u32

  @[link_name="SDL_EVENT_WINDOW_RESIZED"]
  const EVENT_WINDOW_RESIZED u32

  @[link_name="SDL_EVENT_WINDOW_PIXEL_SIZE_CHANGED"]
  const EVENT_WINDOW_PIXEL_SIZE_CHANGED u32

  @[link_name="SDL_EVENT_WINDOW_METAL_VIEW_RESIZED"]
  const EVENT_WINDOW_METAL_VIEW_RESIZED u32

  @[link_name="SDL_EVENT_WINDOW_MINIMIZED"]
  const EVENT_WINDOW_MINIMIZED u32

  @[link_name="SDL_EVENT_WINDOW_MAXIMIZED"]
  const EVENT_WINDOW_MAXIMIZED u32

  @[link_name="SDL_EVENT_WINDOW_RESTORED"]
  const EVENT_WINDOW_RESTORED u32

  @[link_name="SDL_EVENT_WINDOW_MOUSE_ENTER"]
  const EVENT_WINDOW_MOUSE_ENTER u32

  @[link_name="SDL_EVENT_WINDOW_MOUSE_LEAVE"]
  const EVENT_WINDOW_MOUSE_LEAVE u32

  @[link_name="SDL_EVENT_WINDOW_FOCUS_GAINED"]
  const EVENT_WINDOW_FOCUS_GAINED u32

  @[link_name="SDL_EVENT_WINDOW_FOCUS_LOST"]
  const EVENT_WINDOW_FOCUS_LOST u32

  @[link_name="SDL_EVENT_WINDOW_CLOSE_REQUESTED"]
  const EVENT_WINDOW_CLOSE_REQUESTED u32

  @[link_name="SDL_EVENT_WINDOW_HIT_TEST"]
  const EVENT_WINDOW_HIT_TEST u32

  @[link_name="SDL_EVENT_WINDOW_ICCPROF_CHANGED"]
  const EVENT_WINDOW_ICCPROF_CHANGED u32

  @[link_name="SDL_EVENT_WINDOW_DISPLAY_CHANGED"]
  const EVENT_WINDOW_DISPLAY_CHANGED u32

  @[link_name="SDL_EVENT_WINDOW_DISPLAY_SCALE_CHANGED"]
  const EVENT_WINDOW_DISPLAY_SCALE_CHANGED u32

  @[link_name="SDL_EVENT_WINDOW_SAFE_AREA_CHANGED"]
  const EVENT_WINDOW_SAFE_AREA_CHANGED u32

  @[link_name="SDL_EVENT_WINDOW_OCCLUDED"]
  const EVENT_WINDOW_OCCLUDED u32

  @[link_name="SDL_EVENT_WINDOW_ENTER_FULLSCREEN"]
  const EVENT_WINDOW_ENTER_FULLSCREEN u32

  @[link_name="SDL_EVENT_WINDOW_LEAVE_FULLSCREEN"]
  const EVENT_WINDOW_LEAVE_FULLSCREEN u32

  @[link_name="SDL_EVENT_WINDOW_DESTROYED"]
  const EVENT_WINDOW_DESTROYED u32

  @[link_name="SDL_EVENT_WINDOW_HDR_STATE_CHANGED"]
  const EVENT_WINDOW_HDR_STATE_CHANGED u32

  @[link_name="SDL_EVENT_WINDOW_FIRST"]
  const EVENT_WINDOW_FIRST u32

  @[link_name="SDL_EVENT_WINDOW_LAST"]
  const EVENT_WINDOW_LAST u32

  @[link_name="SDL_EVENT_KEY_DOWN"]
  const EVENT_KEY_DOWN u32

  @[link_name="SDL_EVENT_KEY_UP"]
  const EVENT_KEY_UP u32

  @[link_name="SDL_EVENT_TEXT_EDITING"]
  const EVENT_TEXT_EDITING u32

  @[link_name="SDL_EVENT_TEXT_INPUT"]
  const EVENT_TEXT_INPUT u32

  @[link_name="SDL_EVENT_KEYMAP_CHANGED"]
  const EVENT_KEYMAP_CHANGED u32

  @[link_name="SDL_EVENT_KEYBOARD_ADDED"]
  const EVENT_KEYBOARD_ADDED u32

  @[link_name="SDL_EVENT_KEYBOARD_REMOVED"]
  const EVENT_KEYBOARD_REMOVED u32

  @[link_name="SDL_EVENT_TEXT_EDITING_CANDIDATES"]
  const EVENT_TEXT_EDITING_CANDIDATES u32

  @[link_name="SDL_EVENT_MOUSE_MOTION"]
  const EVENT_MOUSE_MOTION u32

  @[link_name="SDL_EVENT_MOUSE_BUTTON_DOWN"]
  const EVENT_MOUSE_BUTTON_DOWN u32

  @[link_name="SDL_EVENT_MOUSE_BUTTON_UP"]
  const EVENT_MOUSE_BUTTON_UP u32

  @[link_name="SDL_EVENT_MOUSE_WHEEL"]
  const EVENT_MOUSE_WHEEL u32

  @[link_name="SDL_EVENT_MOUSE_ADDED"]
  const EVENT_MOUSE_ADDED u32

  @[link_name="SDL_EVENT_MOUSE_REMOVED"]
  const EVENT_MOUSE_REMOVED u32

  @[link_name="SDL_EVENT_JOYSTICK_AXIS_MOTION"]
  const EVENT_JOYSTICK_AXIS_MOTION u32

  @[link_name="SDL_EVENT_JOYSTICK_BALL_MOTION"]
  const EVENT_JOYSTICK_BALL_MOTION u32

  @[link_name="SDL_EVENT_JOYSTICK_HAT_MOTION"]
  const EVENT_JOYSTICK_HAT_MOTION u32

  @[link_name="SDL_EVENT_JOYSTICK_BUTTON_DOWN"]
  const EVENT_JOYSTICK_BUTTON_DOWN u32

  @[link_name="SDL_EVENT_JOYSTICK_BUTTON_UP"]
  const EVENT_JOYSTICK_BUTTON_UP u32

  @[link_name="SDL_EVENT_JOYSTICK_ADDED"]
  const EVENT_JOYSTICK_ADDED u32

  @[link_name="SDL_EVENT_JOYSTICK_REMOVED"]
  const EVENT_JOYSTICK_REMOVED u32

  @[link_name="SDL_EVENT_JOYSTICK_BATTERY_UPDATED"]
  const EVENT_JOYSTICK_BATTERY_UPDATED u32

  @[link_name="SDL_EVENT_JOYSTICK_UPDATE_COMPLETE"]
  const EVENT_JOYSTICK_UPDATE_COMPLETE u32

  @[link_name="SDL_EVENT_GAMEPAD_AXIS_MOTION"]
  const EVENT_GAMEPAD_AXIS_MOTION u32

  @[link_name="SDL_EVENT_GAMEPAD_BUTTON_DOWN"]
  const EVENT_GAMEPAD_BUTTON_DOWN u32

  @[link_name="SDL_EVENT_GAMEPAD_BUTTON_UP"]
  const EVENT_GAMEPAD_BUTTON_UP u32

  @[link_name="SDL_EVENT_GAMEPAD_ADDED"]
  const EVENT_GAMEPAD_ADDED u32

  @[link_name="SDL_EVENT_GAMEPAD_REMOVED"]
  const EVENT_GAMEPAD_REMOVED u32

  @[link_name="SDL_EVENT_GAMEPAD_REMAPPED"]
  const EVENT_GAMEPAD_REMAPPED u32

  @[link_name="SDL_EVENT_GAMEPAD_TOUCHPAD_DOWN"]
  const EVENT_GAMEPAD_TOUCHPAD_DOWN u32

  @[link_name="SDL_EVENT_GAMEPAD_TOUCHPAD_MOTION"]
  const EVENT_GAMEPAD_TOUCHPAD_MOTION u32

  @[link_name="SDL_EVENT_GAMEPAD_TOUCHPAD_UP"]
  const EVENT_GAMEPAD_TOUCHPAD_UP u32

  @[link_name="SDL_EVENT_GAMEPAD_SENSOR_UPDATE"]
  const EVENT_GAMEPAD_SENSOR_UPDATE u32

  @[link_name="SDL_EVENT_GAMEPAD_UPDATE_COMPLETE"]
  const EVENT_GAMEPAD_UPDATE_COMPLETE u32

  @[link_name="SDL_EVENT_GAMEPAD_STEAM_HANDLE_UPDATED"]
  const EVENT_GAMEPAD_STEAM_HANDLE_UPDATED u32

  @[link_name="SDL_EVENT_FINGER_DOWN"]
  const EVENT_FINGER_DOWN u32

  @[link_name="SDL_EVENT_FINGER_UP"]
  const EVENT_FINGER_UP u32

  @[link_name="SDL_EVENT_FINGER_MOTION"]
  const EVENT_FINGER_MOTION u32

  @[link_name="SDL_EVENT_FINGER_CANCELED"]
  const EVENT_FINGER_CANCELED u32

  @[link_name="SDL_EVENT_CLIPBOARD_UPDATE"]
  const EVENT_CLIPBOARD_UPDATE u32

  @[link_name="SDL_EVENT_DROP_FILE"]
  const EVENT_DROP_FILE u32

  @[link_name="SDL_EVENT_DROP_TEXT"]
  const EVENT_DROP_TEXT u32

  @[link_name="SDL_EVENT_DROP_BEGIN"]
  const EVENT_DROP_BEGIN u32

  @[link_name="SDL_EVENT_DROP_COMPLETE"]
  const EVENT_DROP_COMPLETE u32

  @[link_name="SDL_EVENT_DROP_POSITION"]
  const EVENT_DROP_POSITION u32

  @[link_name="SDL_EVENT_AUDIO_DEVICE_ADDED"]
  const EVENT_AUDIO_DEVICE_ADDED u32

  @[link_name="SDL_EVENT_AUDIO_DEVICE_REMOVED"]
  const EVENT_AUDIO_DEVICE_REMOVED u32

  @[link_name="SDL_EVENT_AUDIO_DEVICE_FORMAT_CHANGED"]
  const EVENT_AUDIO_DEVICE_FORMAT_CHANGED u32

  @[link_name="SDL_EVENT_SENSOR_UPDATE"]
  const EVENT_SENSOR_UPDATE u32

  @[link_name="SDL_EVENT_PEN_PROXIMITY_IN"]
  const EVENT_PEN_PROXIMITY_IN u32

  @[link_name="SDL_EVENT_PEN_PROXIMITY_OUT"]
  const EVENT_PEN_PROXIMITY_OUT u32

  @[link_name="SDL_EVENT_PEN_DOWN"]
  const EVENT_PEN_DOWN u32

  @[link_name="SDL_EVENT_PEN_UP"]
  const EVENT_PEN_UP u32

  @[link_name="SDL_EVENT_PEN_BUTTON_DOWN"]
  const EVENT_PEN_BUTTON_DOWN u32

  @[link_name="SDL_EVENT_PEN_BUTTON_UP"]
  const EVENT_PEN_BUTTON_UP u32

  @[link_name="SDL_EVENT_PEN_MOTION"]
  const EVENT_PEN_MOTION u32

  @[link_name="SDL_EVENT_PEN_AXIS"]
  const EVENT_PEN_AXIS u32

  @[link_name="SDL_EVENT_CAMERA_DEVICE_ADDED"]
  const EVENT_CAMERA_DEVICE_ADDED u32

  @[link_name="SDL_EVENT_CAMERA_DEVICE_REMOVED"]
  const EVENT_CAMERA_DEVICE_REMOVED u32

  @[link_name="SDL_EVENT_CAMERA_DEVICE_APPROVED"]
  const EVENT_CAMERA_DEVICE_APPROVED u32

  @[link_name="SDL_EVENT_CAMERA_DEVICE_DENIED"]
  const EVENT_CAMERA_DEVICE_DENIED u32

  @[link_name="SDL_EVENT_RENDER_TARGETS_RESET"]
  const EVENT_RENDER_TARGETS_RESET u32

  @[link_name="SDL_EVENT_RENDER_DEVICE_RESET"]
  const EVENT_RENDER_DEVICE_RESET u32

  @[link_name="SDL_EVENT_RENDER_DEVICE_LOST"]
  const EVENT_RENDER_DEVICE_LOST u32

  @[link_name="SDL_EVENT_PRIVATE0"]
  const EVENT_PRIVATE0 u32

  @[link_name="SDL_EVENT_PRIVATE1"]
  const EVENT_PRIVATE1 u32

  @[link_name="SDL_EVENT_PRIVATE2"]
  const EVENT_PRIVATE2 u32

  @[link_name="SDL_EVENT_PRIVATE3"]
  const EVENT_PRIVATE3 u32

  @[link_name="SDL_EVENT_POLL_SENTINEL"]
  const EVENT_POLL_SENTINEL u32

  @[link_name="SDL_EVENT_USER"]
  const EVENT_USER u32

  @[link_name="SDL_EVENT_LAST"]
  const EVENT_LAST u32

  @[link_name="SDL_EVENT_ENUM_PADDING"]
  const EVENT_ENUM_PADDING u32

  @[link_name="SDL_ADDEVENT"]
  const ADDEVENT u32

  @[link_name="SDL_PEEKEVENT"]
  const PEEKEVENT u32

  @[link_name="SDL_GETEVENT"]
  const GETEVENT u32

  @[link_name="SDL_FOLDER_HOME"]
  const FOLDER_HOME u32

  @[link_name="SDL_FOLDER_DESKTOP"]
  const FOLDER_DESKTOP u32

  @[link_name="SDL_FOLDER_DOCUMENTS"]
  const FOLDER_DOCUMENTS u32

  @[link_name="SDL_FOLDER_DOWNLOADS"]
  const FOLDER_DOWNLOADS u32

  @[link_name="SDL_FOLDER_MUSIC"]
  const FOLDER_MUSIC u32

  @[link_name="SDL_FOLDER_PICTURES"]
  const FOLDER_PICTURES u32

  @[link_name="SDL_FOLDER_PUBLICSHARE"]
  const FOLDER_PUBLICSHARE u32

  @[link_name="SDL_FOLDER_SAVEDGAMES"]
  const FOLDER_SAVEDGAMES u32

  @[link_name="SDL_FOLDER_SCREENSHOTS"]
  const FOLDER_SCREENSHOTS u32

  @[link_name="SDL_FOLDER_TEMPLATES"]
  const FOLDER_TEMPLATES u32

  @[link_name="SDL_FOLDER_VIDEOS"]
  const FOLDER_VIDEOS u32

  @[link_name="SDL_FOLDER_COUNT"]
  const FOLDER_COUNT u32

  @[link_name="SDL_PATHTYPE_NONE"]
  const PATHTYPE_NONE u32

  @[link_name="SDL_PATHTYPE_FILE"]
  const PATHTYPE_FILE u32

  @[link_name="SDL_PATHTYPE_DIRECTORY"]
  const PATHTYPE_DIRECTORY u32

  @[link_name="SDL_PATHTYPE_OTHER"]
  const PATHTYPE_OTHER u32

  @[link_name="SDL_ENUM_CONTINUE"]
  const ENUM_CONTINUE u32

  @[link_name="SDL_ENUM_SUCCESS"]
  const ENUM_SUCCESS u32

  @[link_name="SDL_ENUM_FAILURE"]
  const ENUM_FAILURE u32

  @[link_name="SDL_GPU_PRIMITIVETYPE_TRIANGLELIST"]
  const GPU_PRIMITIVETYPE_TRIANGLELIST u32

  @[link_name="SDL_GPU_PRIMITIVETYPE_TRIANGLESTRIP"]
  const GPU_PRIMITIVETYPE_TRIANGLESTRIP u32

  @[link_name="SDL_GPU_PRIMITIVETYPE_LINELIST"]
  const GPU_PRIMITIVETYPE_LINELIST u32

  @[link_name="SDL_GPU_PRIMITIVETYPE_LINESTRIP"]
  const GPU_PRIMITIVETYPE_LINESTRIP u32

  @[link_name="SDL_GPU_PRIMITIVETYPE_POINTLIST"]
  const GPU_PRIMITIVETYPE_POINTLIST u32

  @[link_name="SDL_GPU_LOADOP_LOAD"]
  const GPU_LOADOP_LOAD u32

  @[link_name="SDL_GPU_LOADOP_CLEAR"]
  const GPU_LOADOP_CLEAR u32

  @[link_name="SDL_GPU_LOADOP_DONT_CARE"]
  const GPU_LOADOP_DONT_CARE u32

  @[link_name="SDL_GPU_STOREOP_STORE"]
  const GPU_STOREOP_STORE u32

  @[link_name="SDL_GPU_STOREOP_DONT_CARE"]
  const GPU_STOREOP_DONT_CARE u32

  @[link_name="SDL_GPU_STOREOP_RESOLVE"]
  const GPU_STOREOP_RESOLVE u32

  @[link_name="SDL_GPU_STOREOP_RESOLVE_AND_STORE"]
  const GPU_STOREOP_RESOLVE_AND_STORE u32

  @[link_name="SDL_GPU_INDEXELEMENTSIZE_16BIT"]
  const GPU_INDEXELEMENTSIZE_16BIT u32

  @[link_name="SDL_GPU_INDEXELEMENTSIZE_32BIT"]
  const GPU_INDEXELEMENTSIZE_32BIT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_INVALID"]
  const GPU_TEXTUREFORMAT_INVALID u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_A8_UNORM"]
  const GPU_TEXTUREFORMAT_A8_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8_UNORM"]
  const GPU_TEXTUREFORMAT_R8_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8_UNORM"]
  const GPU_TEXTUREFORMAT_R8G8_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_UNORM"]
  const GPU_TEXTUREFORMAT_R8G8B8A8_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16_UNORM"]
  const GPU_TEXTUREFORMAT_R16_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_UNORM"]
  const GPU_TEXTUREFORMAT_R16G16_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_UNORM"]
  const GPU_TEXTUREFORMAT_R16G16B16A16_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R10G10B10A2_UNORM"]
  const GPU_TEXTUREFORMAT_R10G10B10A2_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_B5G6R5_UNORM"]
  const GPU_TEXTUREFORMAT_B5G6R5_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_B5G5R5A1_UNORM"]
  const GPU_TEXTUREFORMAT_B5G5R5A1_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_B4G4R4A4_UNORM"]
  const GPU_TEXTUREFORMAT_B4G4R4A4_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_B8G8R8A8_UNORM"]
  const GPU_TEXTUREFORMAT_B8G8R8A8_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC1_RGBA_UNORM"]
  const GPU_TEXTUREFORMAT_BC1_RGBA_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC2_RGBA_UNORM"]
  const GPU_TEXTUREFORMAT_BC2_RGBA_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC3_RGBA_UNORM"]
  const GPU_TEXTUREFORMAT_BC3_RGBA_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC4_R_UNORM"]
  const GPU_TEXTUREFORMAT_BC4_R_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC5_RG_UNORM"]
  const GPU_TEXTUREFORMAT_BC5_RG_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC7_RGBA_UNORM"]
  const GPU_TEXTUREFORMAT_BC7_RGBA_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC6H_RGB_FLOAT"]
  const GPU_TEXTUREFORMAT_BC6H_RGB_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC6H_RGB_UFLOAT"]
  const GPU_TEXTUREFORMAT_BC6H_RGB_UFLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8_SNORM"]
  const GPU_TEXTUREFORMAT_R8_SNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8_SNORM"]
  const GPU_TEXTUREFORMAT_R8G8_SNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_SNORM"]
  const GPU_TEXTUREFORMAT_R8G8B8A8_SNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16_SNORM"]
  const GPU_TEXTUREFORMAT_R16_SNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_SNORM"]
  const GPU_TEXTUREFORMAT_R16G16_SNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_SNORM"]
  const GPU_TEXTUREFORMAT_R16G16B16A16_SNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16_FLOAT"]
  const GPU_TEXTUREFORMAT_R16_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_FLOAT"]
  const GPU_TEXTUREFORMAT_R16G16_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_FLOAT"]
  const GPU_TEXTUREFORMAT_R16G16B16A16_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R32_FLOAT"]
  const GPU_TEXTUREFORMAT_R32_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32_FLOAT"]
  const GPU_TEXTUREFORMAT_R32G32_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32B32A32_FLOAT"]
  const GPU_TEXTUREFORMAT_R32G32B32A32_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R11G11B10_UFLOAT"]
  const GPU_TEXTUREFORMAT_R11G11B10_UFLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8_UINT"]
  const GPU_TEXTUREFORMAT_R8_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8_UINT"]
  const GPU_TEXTUREFORMAT_R8G8_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_UINT"]
  const GPU_TEXTUREFORMAT_R8G8B8A8_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16_UINT"]
  const GPU_TEXTUREFORMAT_R16_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_UINT"]
  const GPU_TEXTUREFORMAT_R16G16_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_UINT"]
  const GPU_TEXTUREFORMAT_R16G16B16A16_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R32_UINT"]
  const GPU_TEXTUREFORMAT_R32_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32_UINT"]
  const GPU_TEXTUREFORMAT_R32G32_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32B32A32_UINT"]
  const GPU_TEXTUREFORMAT_R32G32B32A32_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8_INT"]
  const GPU_TEXTUREFORMAT_R8_INT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8_INT"]
  const GPU_TEXTUREFORMAT_R8G8_INT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_INT"]
  const GPU_TEXTUREFORMAT_R8G8B8A8_INT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16_INT"]
  const GPU_TEXTUREFORMAT_R16_INT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_INT"]
  const GPU_TEXTUREFORMAT_R16G16_INT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_INT"]
  const GPU_TEXTUREFORMAT_R16G16B16A16_INT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R32_INT"]
  const GPU_TEXTUREFORMAT_R32_INT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32_INT"]
  const GPU_TEXTUREFORMAT_R32G32_INT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32B32A32_INT"]
  const GPU_TEXTUREFORMAT_R32G32B32A32_INT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_R8G8B8A8_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_B8G8R8A8_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_B8G8R8A8_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC1_RGBA_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_BC1_RGBA_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC2_RGBA_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_BC2_RGBA_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC3_RGBA_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_BC3_RGBA_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_BC7_RGBA_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_BC7_RGBA_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_D16_UNORM"]
  const GPU_TEXTUREFORMAT_D16_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_D24_UNORM"]
  const GPU_TEXTUREFORMAT_D24_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_D32_FLOAT"]
  const GPU_TEXTUREFORMAT_D32_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_D24_UNORM_S8_UINT"]
  const GPU_TEXTUREFORMAT_D24_UNORM_S8_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_D32_FLOAT_S8_UINT"]
  const GPU_TEXTUREFORMAT_D32_FLOAT_S8_UINT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_4x4_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_4x4_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x4_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_5x4_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x5_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_5x5_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x5_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_6x5_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x6_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_6x6_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x5_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_8x5_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x6_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_8x6_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x8_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_8x8_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x5_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_10x5_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x6_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_10x6_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x8_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_10x8_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x10_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_10x10_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x10_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_12x10_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x12_UNORM"]
  const GPU_TEXTUREFORMAT_ASTC_12x12_UNORM u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_4x4_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_4x4_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x4_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_5x4_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x5_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_5x5_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x5_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_6x5_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x6_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_6x6_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x5_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_8x5_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x6_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_8x6_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x8_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_8x8_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x5_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_10x5_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x6_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_10x6_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x8_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_10x8_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x10_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_10x10_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x10_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_12x10_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x12_UNORM_SRGB"]
  const GPU_TEXTUREFORMAT_ASTC_12x12_UNORM_SRGB u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_4x4_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_4x4_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x4_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_5x4_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x5_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_5x5_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x5_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_6x5_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x6_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_6x6_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x5_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_8x5_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x6_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_8x6_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x8_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_8x8_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x5_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_10x5_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x6_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_10x6_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x8_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_10x8_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x10_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_10x10_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x10_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_12x10_FLOAT u32

  @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x12_FLOAT"]
  const GPU_TEXTUREFORMAT_ASTC_12x12_FLOAT u32

  @[link_name="SDL_GPU_TEXTURETYPE_2D"]
  const GPU_TEXTURETYPE_2D u32

  @[link_name="SDL_GPU_TEXTURETYPE_2D_ARRAY"]
  const GPU_TEXTURETYPE_2D_ARRAY u32

  @[link_name="SDL_GPU_TEXTURETYPE_3D"]
  const GPU_TEXTURETYPE_3D u32

  @[link_name="SDL_GPU_TEXTURETYPE_CUBE"]
  const GPU_TEXTURETYPE_CUBE u32

  @[link_name="SDL_GPU_TEXTURETYPE_CUBE_ARRAY"]
  const GPU_TEXTURETYPE_CUBE_ARRAY u32

  @[link_name="SDL_GPU_SAMPLECOUNT_1"]
  const GPU_SAMPLECOUNT_1 u32

  @[link_name="SDL_GPU_SAMPLECOUNT_2"]
  const GPU_SAMPLECOUNT_2 u32

  @[link_name="SDL_GPU_SAMPLECOUNT_4"]
  const GPU_SAMPLECOUNT_4 u32

  @[link_name="SDL_GPU_SAMPLECOUNT_8"]
  const GPU_SAMPLECOUNT_8 u32

  @[link_name="SDL_GPU_CUBEMAPFACE_POSITIVEX"]
  const GPU_CUBEMAPFACE_POSITIVEX u32

  @[link_name="SDL_GPU_CUBEMAPFACE_NEGATIVEX"]
  const GPU_CUBEMAPFACE_NEGATIVEX u32

  @[link_name="SDL_GPU_CUBEMAPFACE_POSITIVEY"]
  const GPU_CUBEMAPFACE_POSITIVEY u32

  @[link_name="SDL_GPU_CUBEMAPFACE_NEGATIVEY"]
  const GPU_CUBEMAPFACE_NEGATIVEY u32

  @[link_name="SDL_GPU_CUBEMAPFACE_POSITIVEZ"]
  const GPU_CUBEMAPFACE_POSITIVEZ u32

  @[link_name="SDL_GPU_CUBEMAPFACE_NEGATIVEZ"]
  const GPU_CUBEMAPFACE_NEGATIVEZ u32

  @[link_name="SDL_GPU_TRANSFERBUFFERUSAGE_UPLOAD"]
  const GPU_TRANSFERBUFFERUSAGE_UPLOAD u32

  @[link_name="SDL_GPU_TRANSFERBUFFERUSAGE_DOWNLOAD"]
  const GPU_TRANSFERBUFFERUSAGE_DOWNLOAD u32

  @[link_name="SDL_GPU_SHADERSTAGE_VERTEX"]
  const GPU_SHADERSTAGE_VERTEX u32

  @[link_name="SDL_GPU_SHADERSTAGE_FRAGMENT"]
  const GPU_SHADERSTAGE_FRAGMENT u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INVALID"]
  const GPU_VERTEXELEMENTFORMAT_INVALID u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INT"]
  const GPU_VERTEXELEMENTFORMAT_INT u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INT2"]
  const GPU_VERTEXELEMENTFORMAT_INT2 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INT3"]
  const GPU_VERTEXELEMENTFORMAT_INT3 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INT4"]
  const GPU_VERTEXELEMENTFORMAT_INT4 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UINT"]
  const GPU_VERTEXELEMENTFORMAT_UINT u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UINT2"]
  const GPU_VERTEXELEMENTFORMAT_UINT2 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UINT3"]
  const GPU_VERTEXELEMENTFORMAT_UINT3 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UINT4"]
  const GPU_VERTEXELEMENTFORMAT_UINT4 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_FLOAT"]
  const GPU_VERTEXELEMENTFORMAT_FLOAT u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_FLOAT2"]
  const GPU_VERTEXELEMENTFORMAT_FLOAT2 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_FLOAT3"]
  const GPU_VERTEXELEMENTFORMAT_FLOAT3 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_FLOAT4"]
  const GPU_VERTEXELEMENTFORMAT_FLOAT4 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_BYTE2"]
  const GPU_VERTEXELEMENTFORMAT_BYTE2 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_BYTE4"]
  const GPU_VERTEXELEMENTFORMAT_BYTE4 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UBYTE2"]
  const GPU_VERTEXELEMENTFORMAT_UBYTE2 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UBYTE4"]
  const GPU_VERTEXELEMENTFORMAT_UBYTE4 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_BYTE2_NORM"]
  const GPU_VERTEXELEMENTFORMAT_BYTE2_NORM u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_BYTE4_NORM"]
  const GPU_VERTEXELEMENTFORMAT_BYTE4_NORM u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UBYTE2_NORM"]
  const GPU_VERTEXELEMENTFORMAT_UBYTE2_NORM u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UBYTE4_NORM"]
  const GPU_VERTEXELEMENTFORMAT_UBYTE4_NORM u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_SHORT2"]
  const GPU_VERTEXELEMENTFORMAT_SHORT2 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_SHORT4"]
  const GPU_VERTEXELEMENTFORMAT_SHORT4 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_USHORT2"]
  const GPU_VERTEXELEMENTFORMAT_USHORT2 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_USHORT4"]
  const GPU_VERTEXELEMENTFORMAT_USHORT4 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_SHORT2_NORM"]
  const GPU_VERTEXELEMENTFORMAT_SHORT2_NORM u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_SHORT4_NORM"]
  const GPU_VERTEXELEMENTFORMAT_SHORT4_NORM u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_USHORT2_NORM"]
  const GPU_VERTEXELEMENTFORMAT_USHORT2_NORM u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_USHORT4_NORM"]
  const GPU_VERTEXELEMENTFORMAT_USHORT4_NORM u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_HALF2"]
  const GPU_VERTEXELEMENTFORMAT_HALF2 u32

  @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_HALF4"]
  const GPU_VERTEXELEMENTFORMAT_HALF4 u32

  @[link_name="SDL_GPU_VERTEXINPUTRATE_VERTEX"]
  const GPU_VERTEXINPUTRATE_VERTEX u32

  @[link_name="SDL_GPU_VERTEXINPUTRATE_INSTANCE"]
  const GPU_VERTEXINPUTRATE_INSTANCE u32

  @[link_name="SDL_GPU_FILLMODE_FILL"]
  const GPU_FILLMODE_FILL u32

  @[link_name="SDL_GPU_FILLMODE_LINE"]
  const GPU_FILLMODE_LINE u32

  @[link_name="SDL_GPU_CULLMODE_NONE"]
  const GPU_CULLMODE_NONE u32

  @[link_name="SDL_GPU_CULLMODE_FRONT"]
  const GPU_CULLMODE_FRONT u32

  @[link_name="SDL_GPU_CULLMODE_BACK"]
  const GPU_CULLMODE_BACK u32

  @[link_name="SDL_GPU_FRONTFACE_COUNTER_CLOCKWISE"]
  const GPU_FRONTFACE_COUNTER_CLOCKWISE u32

  @[link_name="SDL_GPU_FRONTFACE_CLOCKWISE"]
  const GPU_FRONTFACE_CLOCKWISE u32

  @[link_name="SDL_GPU_COMPAREOP_INVALID"]
  const GPU_COMPAREOP_INVALID u32

  @[link_name="SDL_GPU_COMPAREOP_NEVER"]
  const GPU_COMPAREOP_NEVER u32

  @[link_name="SDL_GPU_COMPAREOP_LESS"]
  const GPU_COMPAREOP_LESS u32

  @[link_name="SDL_GPU_COMPAREOP_EQUAL"]
  const GPU_COMPAREOP_EQUAL u32

  @[link_name="SDL_GPU_COMPAREOP_LESS_OR_EQUAL"]
  const GPU_COMPAREOP_LESS_OR_EQUAL u32

  @[link_name="SDL_GPU_COMPAREOP_GREATER"]
  const GPU_COMPAREOP_GREATER u32

  @[link_name="SDL_GPU_COMPAREOP_NOT_EQUAL"]
  const GPU_COMPAREOP_NOT_EQUAL u32

  @[link_name="SDL_GPU_COMPAREOP_GREATER_OR_EQUAL"]
  const GPU_COMPAREOP_GREATER_OR_EQUAL u32

  @[link_name="SDL_GPU_COMPAREOP_ALWAYS"]
  const GPU_COMPAREOP_ALWAYS u32

  @[link_name="SDL_GPU_STENCILOP_INVALID"]
  const GPU_STENCILOP_INVALID u32

  @[link_name="SDL_GPU_STENCILOP_KEEP"]
  const GPU_STENCILOP_KEEP u32

  @[link_name="SDL_GPU_STENCILOP_ZERO"]
  const GPU_STENCILOP_ZERO u32

  @[link_name="SDL_GPU_STENCILOP_REPLACE"]
  const GPU_STENCILOP_REPLACE u32

  @[link_name="SDL_GPU_STENCILOP_INCREMENT_AND_CLAMP"]
  const GPU_STENCILOP_INCREMENT_AND_CLAMP u32

  @[link_name="SDL_GPU_STENCILOP_DECREMENT_AND_CLAMP"]
  const GPU_STENCILOP_DECREMENT_AND_CLAMP u32

  @[link_name="SDL_GPU_STENCILOP_INVERT"]
  const GPU_STENCILOP_INVERT u32

  @[link_name="SDL_GPU_STENCILOP_INCREMENT_AND_WRAP"]
  const GPU_STENCILOP_INCREMENT_AND_WRAP u32

  @[link_name="SDL_GPU_STENCILOP_DECREMENT_AND_WRAP"]
  const GPU_STENCILOP_DECREMENT_AND_WRAP u32

  @[link_name="SDL_GPU_BLENDOP_INVALID"]
  const GPU_BLENDOP_INVALID u32

  @[link_name="SDL_GPU_BLENDOP_ADD"]
  const GPU_BLENDOP_ADD u32

  @[link_name="SDL_GPU_BLENDOP_SUBTRACT"]
  const GPU_BLENDOP_SUBTRACT u32

  @[link_name="SDL_GPU_BLENDOP_REVERSE_SUBTRACT"]
  const GPU_BLENDOP_REVERSE_SUBTRACT u32

  @[link_name="SDL_GPU_BLENDOP_MIN"]
  const GPU_BLENDOP_MIN u32

  @[link_name="SDL_GPU_BLENDOP_MAX"]
  const GPU_BLENDOP_MAX u32

  @[link_name="SDL_GPU_BLENDFACTOR_INVALID"]
  const GPU_BLENDFACTOR_INVALID u32

  @[link_name="SDL_GPU_BLENDFACTOR_ZERO"]
  const GPU_BLENDFACTOR_ZERO u32

  @[link_name="SDL_GPU_BLENDFACTOR_ONE"]
  const GPU_BLENDFACTOR_ONE u32

  @[link_name="SDL_GPU_BLENDFACTOR_SRC_COLOR"]
  const GPU_BLENDFACTOR_SRC_COLOR u32

  @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_SRC_COLOR"]
  const GPU_BLENDFACTOR_ONE_MINUS_SRC_COLOR u32

  @[link_name="SDL_GPU_BLENDFACTOR_DST_COLOR"]
  const GPU_BLENDFACTOR_DST_COLOR u32

  @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_DST_COLOR"]
  const GPU_BLENDFACTOR_ONE_MINUS_DST_COLOR u32

  @[link_name="SDL_GPU_BLENDFACTOR_SRC_ALPHA"]
  const GPU_BLENDFACTOR_SRC_ALPHA u32

  @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_SRC_ALPHA"]
  const GPU_BLENDFACTOR_ONE_MINUS_SRC_ALPHA u32

  @[link_name="SDL_GPU_BLENDFACTOR_DST_ALPHA"]
  const GPU_BLENDFACTOR_DST_ALPHA u32

  @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_DST_ALPHA"]
  const GPU_BLENDFACTOR_ONE_MINUS_DST_ALPHA u32

  @[link_name="SDL_GPU_BLENDFACTOR_CONSTANT_COLOR"]
  const GPU_BLENDFACTOR_CONSTANT_COLOR u32

  @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_CONSTANT_COLOR"]
  const GPU_BLENDFACTOR_ONE_MINUS_CONSTANT_COLOR u32

  @[link_name="SDL_GPU_BLENDFACTOR_SRC_ALPHA_SATURATE"]
  const GPU_BLENDFACTOR_SRC_ALPHA_SATURATE u32

  @[link_name="SDL_GPU_FILTER_NEAREST"]
  const GPU_FILTER_NEAREST u32

  @[link_name="SDL_GPU_FILTER_LINEAR"]
  const GPU_FILTER_LINEAR u32

  @[link_name="SDL_GPU_SAMPLERMIPMAPMODE_NEAREST"]
  const GPU_SAMPLERMIPMAPMODE_NEAREST u32

  @[link_name="SDL_GPU_SAMPLERMIPMAPMODE_LINEAR"]
  const GPU_SAMPLERMIPMAPMODE_LINEAR u32

  @[link_name="SDL_GPU_SAMPLERADDRESSMODE_REPEAT"]
  const GPU_SAMPLERADDRESSMODE_REPEAT u32

  @[link_name="SDL_GPU_SAMPLERADDRESSMODE_MIRRORED_REPEAT"]
  const GPU_SAMPLERADDRESSMODE_MIRRORED_REPEAT u32

  @[link_name="SDL_GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE"]
  const GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE u32

  @[link_name="SDL_GPU_PRESENTMODE_VSYNC"]
  const GPU_PRESENTMODE_VSYNC u32

  @[link_name="SDL_GPU_PRESENTMODE_IMMEDIATE"]
  const GPU_PRESENTMODE_IMMEDIATE u32

  @[link_name="SDL_GPU_PRESENTMODE_MAILBOX"]
  const GPU_PRESENTMODE_MAILBOX u32

  @[link_name="SDL_GPU_SWAPCHAINCOMPOSITION_SDR"]
  const GPU_SWAPCHAINCOMPOSITION_SDR u32

  @[link_name="SDL_GPU_SWAPCHAINCOMPOSITION_SDR_LINEAR"]
  const GPU_SWAPCHAINCOMPOSITION_SDR_LINEAR u32

  @[link_name="SDL_GPU_SWAPCHAINCOMPOSITION_HDR_EXTENDED_LINEAR"]
  const GPU_SWAPCHAINCOMPOSITION_HDR_EXTENDED_LINEAR u32

  @[link_name="SDL_GPU_SWAPCHAINCOMPOSITION_HDR10_ST2084"]
  const GPU_SWAPCHAINCOMPOSITION_HDR10_ST2084 u32

  @[link_name="SDL_HID_API_BUS_UNKNOWN"]
  const HID_API_BUS_UNKNOWN u32

  @[link_name="SDL_HID_API_BUS_USB"]
  const HID_API_BUS_USB u32

  @[link_name="SDL_HID_API_BUS_BLUETOOTH"]
  const HID_API_BUS_BLUETOOTH u32

  @[link_name="SDL_HID_API_BUS_I2C"]
  const HID_API_BUS_I2C u32

  @[link_name="SDL_HID_API_BUS_SPI"]
  const HID_API_BUS_SPI u32

  @[link_name="SDL_HINT_DEFAULT"]
  const HINT_DEFAULT u32

  @[link_name="SDL_HINT_NORMAL"]
  const HINT_NORMAL u32

  @[link_name="SDL_HINT_OVERRIDE"]
  const HINT_OVERRIDE u32

  @[link_name="SDL_APP_CONTINUE"]
  const APP_CONTINUE u32

  @[link_name="SDL_APP_SUCCESS"]
  const APP_SUCCESS u32

  @[link_name="SDL_APP_FAILURE"]
  const APP_FAILURE u32

  @[link_name="SDL_LOG_CATEGORY_APPLICATION"]
  const LOG_CATEGORY_APPLICATION u32

  @[link_name="SDL_LOG_CATEGORY_ERROR"]
  const LOG_CATEGORY_ERROR u32

  @[link_name="SDL_LOG_CATEGORY_ASSERT"]
  const LOG_CATEGORY_ASSERT u32

  @[link_name="SDL_LOG_CATEGORY_SYSTEM"]
  const LOG_CATEGORY_SYSTEM u32

  @[link_name="SDL_LOG_CATEGORY_AUDIO"]
  const LOG_CATEGORY_AUDIO u32

  @[link_name="SDL_LOG_CATEGORY_VIDEO"]
  const LOG_CATEGORY_VIDEO u32

  @[link_name="SDL_LOG_CATEGORY_RENDER"]
  const LOG_CATEGORY_RENDER u32

  @[link_name="SDL_LOG_CATEGORY_INPUT"]
  const LOG_CATEGORY_INPUT u32

  @[link_name="SDL_LOG_CATEGORY_TEST"]
  const LOG_CATEGORY_TEST u32

  @[link_name="SDL_LOG_CATEGORY_GPU"]
  const LOG_CATEGORY_GPU u32

  @[link_name="SDL_LOG_CATEGORY_RESERVED2"]
  const LOG_CATEGORY_RESERVED2 u32

  @[link_name="SDL_LOG_CATEGORY_RESERVED3"]
  const LOG_CATEGORY_RESERVED3 u32

  @[link_name="SDL_LOG_CATEGORY_RESERVED4"]
  const LOG_CATEGORY_RESERVED4 u32

  @[link_name="SDL_LOG_CATEGORY_RESERVED5"]
  const LOG_CATEGORY_RESERVED5 u32

  @[link_name="SDL_LOG_CATEGORY_RESERVED6"]
  const LOG_CATEGORY_RESERVED6 u32

  @[link_name="SDL_LOG_CATEGORY_RESERVED7"]
  const LOG_CATEGORY_RESERVED7 u32

  @[link_name="SDL_LOG_CATEGORY_RESERVED8"]
  const LOG_CATEGORY_RESERVED8 u32

  @[link_name="SDL_LOG_CATEGORY_RESERVED9"]
  const LOG_CATEGORY_RESERVED9 u32

  @[link_name="SDL_LOG_CATEGORY_RESERVED10"]
  const LOG_CATEGORY_RESERVED10 u32

  @[link_name="SDL_LOG_CATEGORY_CUSTOM"]
  const LOG_CATEGORY_CUSTOM u32

  @[link_name="SDL_LOG_PRIORITY_INVALID"]
  const LOG_PRIORITY_INVALID u32

  @[link_name="SDL_LOG_PRIORITY_TRACE"]
  const LOG_PRIORITY_TRACE u32

  @[link_name="SDL_LOG_PRIORITY_VERBOSE"]
  const LOG_PRIORITY_VERBOSE u32

  @[link_name="SDL_LOG_PRIORITY_DEBUG"]
  const LOG_PRIORITY_DEBUG u32

  @[link_name="SDL_LOG_PRIORITY_INFO"]
  const LOG_PRIORITY_INFO u32

  @[link_name="SDL_LOG_PRIORITY_WARN"]
  const LOG_PRIORITY_WARN u32

  @[link_name="SDL_LOG_PRIORITY_ERROR"]
  const LOG_PRIORITY_ERROR u32

  @[link_name="SDL_LOG_PRIORITY_CRITICAL"]
  const LOG_PRIORITY_CRITICAL u32

  @[link_name="SDL_LOG_PRIORITY_COUNT"]
  const LOG_PRIORITY_COUNT u32

  @[link_name="SDL_MESSAGEBOX_COLOR_BACKGROUND"]
  const MESSAGEBOX_COLOR_BACKGROUND u32

  @[link_name="SDL_MESSAGEBOX_COLOR_TEXT"]
  const MESSAGEBOX_COLOR_TEXT u32

  @[link_name="SDL_MESSAGEBOX_COLOR_BUTTON_BORDER"]
  const MESSAGEBOX_COLOR_BUTTON_BORDER u32

  @[link_name="SDL_MESSAGEBOX_COLOR_BUTTON_BACKGROUND"]
  const MESSAGEBOX_COLOR_BUTTON_BACKGROUND u32

  @[link_name="SDL_MESSAGEBOX_COLOR_BUTTON_SELECTED"]
  const MESSAGEBOX_COLOR_BUTTON_SELECTED u32

  @[link_name="SDL_MESSAGEBOX_COLOR_COUNT"]
  const MESSAGEBOX_COLOR_COUNT u32

  @[link_name="SDL_PROCESS_STDIO_INHERITED"]
  const PROCESS_STDIO_INHERITED u32

  @[link_name="SDL_PROCESS_STDIO_NULL"]
  const PROCESS_STDIO_NULL u32

  @[link_name="SDL_PROCESS_STDIO_APP"]
  const PROCESS_STDIO_APP u32

  @[link_name="SDL_PROCESS_STDIO_REDIRECT"]
  const PROCESS_STDIO_REDIRECT u32

  @[link_name="SDL_TEXTUREACCESS_STATIC"]
  const TEXTUREACCESS_STATIC u32

  @[link_name="SDL_TEXTUREACCESS_STREAMING"]
  const TEXTUREACCESS_STREAMING u32

  @[link_name="SDL_TEXTUREACCESS_TARGET"]
  const TEXTUREACCESS_TARGET u32

  @[link_name="SDL_LOGICAL_PRESENTATION_DISABLED"]
  const LOGICAL_PRESENTATION_DISABLED u32

  @[link_name="SDL_LOGICAL_PRESENTATION_STRETCH"]
  const LOGICAL_PRESENTATION_STRETCH u32

  @[link_name="SDL_LOGICAL_PRESENTATION_LETTERBOX"]
  const LOGICAL_PRESENTATION_LETTERBOX u32

  @[link_name="SDL_LOGICAL_PRESENTATION_OVERSCAN"]
  const LOGICAL_PRESENTATION_OVERSCAN u32

  @[link_name="SDL_LOGICAL_PRESENTATION_INTEGER_SCALE"]
  const LOGICAL_PRESENTATION_INTEGER_SCALE u32

  @[link_name="SDL_SANDBOX_NONE"]
  const SANDBOX_NONE u32

  @[link_name="SDL_SANDBOX_UNKNOWN_CONTAINER"]
  const SANDBOX_UNKNOWN_CONTAINER u32

  @[link_name="SDL_SANDBOX_FLATPAK"]
  const SANDBOX_FLATPAK u32

  @[link_name="SDL_SANDBOX_SNAP"]
  const SANDBOX_SNAP u32

  @[link_name="SDL_SANDBOX_MACOS"]
  const SANDBOX_MACOS u32

  @[link_name="SDL_DATE_FORMAT_YYYYMMDD"]
  const DATE_FORMAT_YYYYMMDD u32

  @[link_name="SDL_DATE_FORMAT_DDMMYYYY"]
  const DATE_FORMAT_DDMMYYYY u32

  @[link_name="SDL_DATE_FORMAT_MMDDYYYY"]
  const DATE_FORMAT_MMDDYYYY u32

  @[link_name="SDL_TIME_FORMAT_24HR"]
  const TIME_FORMAT_24HR u32

  @[link_name="SDL_TIME_FORMAT_12HR"]
  const TIME_FORMAT_12HR u32

  @[link_name="SDLK_0"]
  const SDLK_0 u32

  @[link_name="SDLK_1"]
  const SDLK_1 u32

  @[link_name="SDLK_2"]
  const SDLK_2 u32

  @[link_name="SDLK_3"]
  const SDLK_3 u32

  @[link_name="SDLK_4"]
  const SDLK_4 u32

  @[link_name="SDLK_5"]
  const SDLK_5 u32

  @[link_name="SDLK_6"]
  const SDLK_6 u32

  @[link_name="SDLK_7"]
  const SDLK_7 u32

  @[link_name="SDLK_8"]
  const SDLK_8 u32

  @[link_name="SDLK_9"]
  const SDLK_9 u32

  @[link_name="SDLK_A"]
  const SDLK_A u32

  @[link_name="SDLK_AC_BACK"]
  const SDLK_AC_BACK u32

  @[link_name="SDLK_AC_BOOKMARKS"]
  const SDLK_AC_BOOKMARKS u32

  @[link_name="SDLK_AC_CLOSE"]
  const SDLK_AC_CLOSE u32

  @[link_name="SDLK_AC_EXIT"]
  const SDLK_AC_EXIT u32

  @[link_name="SDLK_AC_FORWARD"]
  const SDLK_AC_FORWARD u32

  @[link_name="SDLK_AC_HOME"]
  const SDLK_AC_HOME u32

  @[link_name="SDLK_AC_NEW"]
  const SDLK_AC_NEW u32

  @[link_name="SDLK_AC_OPEN"]
  const SDLK_AC_OPEN u32

  @[link_name="SDLK_AC_PRINT"]
  const SDLK_AC_PRINT u32

  @[link_name="SDLK_AC_PROPERTIES"]
  const SDLK_AC_PROPERTIES u32

  @[link_name="SDLK_AC_REFRESH"]
  const SDLK_AC_REFRESH u32

  @[link_name="SDLK_AC_SAVE"]
  const SDLK_AC_SAVE u32

  @[link_name="SDLK_AC_SEARCH"]
  const SDLK_AC_SEARCH u32

  @[link_name="SDLK_AC_STOP"]
  const SDLK_AC_STOP u32

  @[link_name="SDLK_AGAIN"]
  const SDLK_AGAIN u32

  @[link_name="SDLK_ALTERASE"]
  const SDLK_ALTERASE u32

  @[link_name="SDLK_AMPERSAND"]
  const SDLK_AMPERSAND u32

  @[link_name="SDLK_APOSTROPHE"]
  const SDLK_APOSTROPHE u32

  @[link_name="SDLK_APPLICATION"]
  const SDLK_APPLICATION u32

  @[link_name="SDLK_ASTERISK"]
  const SDLK_ASTERISK u32

  @[link_name="SDLK_AT"]
  const SDLK_AT u32

  @[link_name="SDLK_B"]
  const SDLK_B u32

  @[link_name="SDLK_BACKSLASH"]
  const SDLK_BACKSLASH u32

  @[link_name="SDLK_BACKSPACE"]
  const SDLK_BACKSPACE u32

  @[link_name="SDLK_C"]
  const SDLK_C u32

  @[link_name="SDLK_CALL"]
  const SDLK_CALL u32

  @[link_name="SDLK_CANCEL"]
  const SDLK_CANCEL u32

  @[link_name="SDLK_CAPSLOCK"]
  const SDLK_CAPSLOCK u32

  @[link_name="SDLK_CARET"]
  const SDLK_CARET u32

  @[link_name="SDLK_CHANNEL_DECREMENT"]
  const SDLK_CHANNEL_DECREMENT u32

  @[link_name="SDLK_CHANNEL_INCREMENT"]
  const SDLK_CHANNEL_INCREMENT u32

  @[link_name="SDLK_CLEAR"]
  const SDLK_CLEAR u32

  @[link_name="SDLK_CLEARAGAIN"]
  const SDLK_CLEARAGAIN u32

  @[link_name="SDLK_COLON"]
  const SDLK_COLON u32

  @[link_name="SDLK_COMMA"]
  const SDLK_COMMA u32

  @[link_name="SDLK_COPY"]
  const SDLK_COPY u32

  @[link_name="SDLK_CRSEL"]
  const SDLK_CRSEL u32

  @[link_name="SDLK_CURRENCYSUBUNIT"]
  const SDLK_CURRENCYSUBUNIT u32

  @[link_name="SDLK_CURRENCYUNIT"]
  const SDLK_CURRENCYUNIT u32

  @[link_name="SDLK_CUT"]
  const SDLK_CUT u32

  @[link_name="SDLK_D"]
  const SDLK_D u32

  @[link_name="SDLK_DBLAPOSTROPHE"]
  const SDLK_DBLAPOSTROPHE u32

  @[link_name="SDLK_DECIMALSEPARATOR"]
  const SDLK_DECIMALSEPARATOR u32

  @[link_name="SDLK_DELETE"]
  const SDLK_DELETE u32

  @[link_name="SDLK_DOLLAR"]
  const SDLK_DOLLAR u32

  @[link_name="SDLK_DOWN"]
  const SDLK_DOWN u32

  @[link_name="SDLK_E"]
  const SDLK_E u32

  @[link_name="SDLK_END"]
  const SDLK_END u32

  @[link_name="SDLK_ENDCALL"]
  const SDLK_ENDCALL u32

  @[link_name="SDLK_EQUALS"]
  const SDLK_EQUALS u32

  @[link_name="SDLK_ESCAPE"]
  const SDLK_ESCAPE u32

  @[link_name="SDLK_EXCLAIM"]
  const SDLK_EXCLAIM u32

  @[link_name="SDLK_EXECUTE"]
  const SDLK_EXECUTE u32

  @[link_name="SDLK_EXSEL"]
  const SDLK_EXSEL u32

  @[link_name="SDLK_EXTENDED_MASK"]
  const SDLK_EXTENDED_MASK u32

  @[link_name="SDLK_F"]
  const SDLK_F u32

  @[link_name="SDLK_F1"]
  const SDLK_F1 u32

  @[link_name="SDLK_F10"]
  const SDLK_F10 u32

  @[link_name="SDLK_F11"]
  const SDLK_F11 u32

  @[link_name="SDLK_F12"]
  const SDLK_F12 u32

  @[link_name="SDLK_F13"]
  const SDLK_F13 u32

  @[link_name="SDLK_F14"]
  const SDLK_F14 u32

  @[link_name="SDLK_F15"]
  const SDLK_F15 u32

  @[link_name="SDLK_F16"]
  const SDLK_F16 u32

  @[link_name="SDLK_F17"]
  const SDLK_F17 u32

  @[link_name="SDLK_F18"]
  const SDLK_F18 u32

  @[link_name="SDLK_F19"]
  const SDLK_F19 u32

  @[link_name="SDLK_F2"]
  const SDLK_F2 u32

  @[link_name="SDLK_F20"]
  const SDLK_F20 u32

  @[link_name="SDLK_F21"]
  const SDLK_F21 u32

  @[link_name="SDLK_F22"]
  const SDLK_F22 u32

  @[link_name="SDLK_F23"]
  const SDLK_F23 u32

  @[link_name="SDLK_F24"]
  const SDLK_F24 u32

  @[link_name="SDLK_F3"]
  const SDLK_F3 u32

  @[link_name="SDLK_F4"]
  const SDLK_F4 u32

  @[link_name="SDLK_F5"]
  const SDLK_F5 u32

  @[link_name="SDLK_F6"]
  const SDLK_F6 u32

  @[link_name="SDLK_F7"]
  const SDLK_F7 u32

  @[link_name="SDLK_F8"]
  const SDLK_F8 u32

  @[link_name="SDLK_F9"]
  const SDLK_F9 u32

  @[link_name="SDLK_FIND"]
  const SDLK_FIND u32

  @[link_name="SDLK_G"]
  const SDLK_G u32

  @[link_name="SDLK_GRAVE"]
  const SDLK_GRAVE u32

  @[link_name="SDLK_GREATER"]
  const SDLK_GREATER u32

  @[link_name="SDLK_H"]
  const SDLK_H u32

  @[link_name="SDLK_HASH"]
  const SDLK_HASH u32

  @[link_name="SDLK_HELP"]
  const SDLK_HELP u32

  @[link_name="SDLK_HOME"]
  const SDLK_HOME u32

  @[link_name="SDLK_I"]
  const SDLK_I u32

  @[link_name="SDLK_INSERT"]
  const SDLK_INSERT u32

  @[link_name="SDLK_J"]
  const SDLK_J u32

  @[link_name="SDLK_K"]
  const SDLK_K u32

  @[link_name="SDLK_KP_0"]
  const SDLK_KP_0 u32

  @[link_name="SDLK_KP_00"]
  const SDLK_KP_00 u32

  @[link_name="SDLK_KP_000"]
  const SDLK_KP_000 u32

  @[link_name="SDLK_KP_1"]
  const SDLK_KP_1 u32

  @[link_name="SDLK_KP_2"]
  const SDLK_KP_2 u32

  @[link_name="SDLK_KP_3"]
  const SDLK_KP_3 u32

  @[link_name="SDLK_KP_4"]
  const SDLK_KP_4 u32

  @[link_name="SDLK_KP_5"]
  const SDLK_KP_5 u32

  @[link_name="SDLK_KP_6"]
  const SDLK_KP_6 u32

  @[link_name="SDLK_KP_7"]
  const SDLK_KP_7 u32

  @[link_name="SDLK_KP_8"]
  const SDLK_KP_8 u32

  @[link_name="SDLK_KP_9"]
  const SDLK_KP_9 u32

  @[link_name="SDLK_KP_A"]
  const SDLK_KP_A u32

  @[link_name="SDLK_KP_AMPERSAND"]
  const SDLK_KP_AMPERSAND u32

  @[link_name="SDLK_KP_AT"]
  const SDLK_KP_AT u32

  @[link_name="SDLK_KP_B"]
  const SDLK_KP_B u32

  @[link_name="SDLK_KP_BACKSPACE"]
  const SDLK_KP_BACKSPACE u32

  @[link_name="SDLK_KP_BINARY"]
  const SDLK_KP_BINARY u32

  @[link_name="SDLK_KP_C"]
  const SDLK_KP_C u32

  @[link_name="SDLK_KP_CLEAR"]
  const SDLK_KP_CLEAR u32

  @[link_name="SDLK_KP_CLEARENTRY"]
  const SDLK_KP_CLEARENTRY u32

  @[link_name="SDLK_KP_COLON"]
  const SDLK_KP_COLON u32

  @[link_name="SDLK_KP_COMMA"]
  const SDLK_KP_COMMA u32

  @[link_name="SDLK_KP_D"]
  const SDLK_KP_D u32

  @[link_name="SDLK_KP_DBLAMPERSAND"]
  const SDLK_KP_DBLAMPERSAND u32

  @[link_name="SDLK_KP_DBLVERTICALBAR"]
  const SDLK_KP_DBLVERTICALBAR u32

  @[link_name="SDLK_KP_DECIMAL"]
  const SDLK_KP_DECIMAL u32

  @[link_name="SDLK_KP_DIVIDE"]
  const SDLK_KP_DIVIDE u32

  @[link_name="SDLK_KP_E"]
  const SDLK_KP_E u32

  @[link_name="SDLK_KP_ENTER"]
  const SDLK_KP_ENTER u32

  @[link_name="SDLK_KP_EQUALS"]
  const SDLK_KP_EQUALS u32

  @[link_name="SDLK_KP_EQUALSAS400"]
  const SDLK_KP_EQUALSAS400 u32

  @[link_name="SDLK_KP_EXCLAM"]
  const SDLK_KP_EXCLAM u32

  @[link_name="SDLK_KP_F"]
  const SDLK_KP_F u32

  @[link_name="SDLK_KP_GREATER"]
  const SDLK_KP_GREATER u32

  @[link_name="SDLK_KP_HASH"]
  const SDLK_KP_HASH u32

  @[link_name="SDLK_KP_HEXADECIMAL"]
  const SDLK_KP_HEXADECIMAL u32

  @[link_name="SDLK_KP_LEFTBRACE"]
  const SDLK_KP_LEFTBRACE u32

  @[link_name="SDLK_KP_LEFTPAREN"]
  const SDLK_KP_LEFTPAREN u32

  @[link_name="SDLK_KP_LESS"]
  const SDLK_KP_LESS u32

  @[link_name="SDLK_KP_MEMADD"]
  const SDLK_KP_MEMADD u32

  @[link_name="SDLK_KP_MEMCLEAR"]
  const SDLK_KP_MEMCLEAR u32

  @[link_name="SDLK_KP_MEMDIVIDE"]
  const SDLK_KP_MEMDIVIDE u32

  @[link_name="SDLK_KP_MEMMULTIPLY"]
  const SDLK_KP_MEMMULTIPLY u32

  @[link_name="SDLK_KP_MEMRECALL"]
  const SDLK_KP_MEMRECALL u32

  @[link_name="SDLK_KP_MEMSTORE"]
  const SDLK_KP_MEMSTORE u32

  @[link_name="SDLK_KP_MEMSUBTRACT"]
  const SDLK_KP_MEMSUBTRACT u32

  @[link_name="SDLK_KP_MINUS"]
  const SDLK_KP_MINUS u32

  @[link_name="SDLK_KP_MULTIPLY"]
  const SDLK_KP_MULTIPLY u32

  @[link_name="SDLK_KP_OCTAL"]
  const SDLK_KP_OCTAL u32

  @[link_name="SDLK_KP_PERCENT"]
  const SDLK_KP_PERCENT u32

  @[link_name="SDLK_KP_PERIOD"]
  const SDLK_KP_PERIOD u32

  @[link_name="SDLK_KP_PLUS"]
  const SDLK_KP_PLUS u32

  @[link_name="SDLK_KP_PLUSMINUS"]
  const SDLK_KP_PLUSMINUS u32

  @[link_name="SDLK_KP_POWER"]
  const SDLK_KP_POWER u32

  @[link_name="SDLK_KP_RIGHTBRACE"]
  const SDLK_KP_RIGHTBRACE u32

  @[link_name="SDLK_KP_RIGHTPAREN"]
  const SDLK_KP_RIGHTPAREN u32

  @[link_name="SDLK_KP_SPACE"]
  const SDLK_KP_SPACE u32

  @[link_name="SDLK_KP_TAB"]
  const SDLK_KP_TAB u32

  @[link_name="SDLK_KP_VERTICALBAR"]
  const SDLK_KP_VERTICALBAR u32

  @[link_name="SDLK_KP_XOR"]
  const SDLK_KP_XOR u32

  @[link_name="SDLK_L"]
  const SDLK_L u32

  @[link_name="SDLK_LALT"]
  const SDLK_LALT u32

  @[link_name="SDLK_LCTRL"]
  const SDLK_LCTRL u32

  @[link_name="SDLK_LEFT"]
  const SDLK_LEFT u32

  @[link_name="SDLK_LEFTBRACE"]
  const SDLK_LEFTBRACE u32

  @[link_name="SDLK_LEFTBRACKET"]
  const SDLK_LEFTBRACKET u32

  @[link_name="SDLK_LEFTPAREN"]
  const SDLK_LEFTPAREN u32

  @[link_name="SDLK_LEFT_TAB"]
  const SDLK_LEFT_TAB u32

  @[link_name="SDLK_LESS"]
  const SDLK_LESS u32

  @[link_name="SDLK_LEVEL5_SHIFT"]
  const SDLK_LEVEL5_SHIFT u32

  @[link_name="SDLK_LGUI"]
  const SDLK_LGUI u32

  @[link_name="SDLK_LHYPER"]
  const SDLK_LHYPER u32

  @[link_name="SDLK_LMETA"]
  const SDLK_LMETA u32

  @[link_name="SDLK_LSHIFT"]
  const SDLK_LSHIFT u32

  @[link_name="SDLK_M"]
  const SDLK_M u32

  @[link_name="SDLK_MEDIA_EJECT"]
  const SDLK_MEDIA_EJECT u32

  @[link_name="SDLK_MEDIA_FAST_FORWARD"]
  const SDLK_MEDIA_FAST_FORWARD u32

  @[link_name="SDLK_MEDIA_NEXT_TRACK"]
  const SDLK_MEDIA_NEXT_TRACK u32

  @[link_name="SDLK_MEDIA_PAUSE"]
  const SDLK_MEDIA_PAUSE u32

  @[link_name="SDLK_MEDIA_PLAY"]
  const SDLK_MEDIA_PLAY u32

  @[link_name="SDLK_MEDIA_PLAY_PAUSE"]
  const SDLK_MEDIA_PLAY_PAUSE u32

  @[link_name="SDLK_MEDIA_PREVIOUS_TRACK"]
  const SDLK_MEDIA_PREVIOUS_TRACK u32

  @[link_name="SDLK_MEDIA_RECORD"]
  const SDLK_MEDIA_RECORD u32

  @[link_name="SDLK_MEDIA_REWIND"]
  const SDLK_MEDIA_REWIND u32

  @[link_name="SDLK_MEDIA_SELECT"]
  const SDLK_MEDIA_SELECT u32

  @[link_name="SDLK_MEDIA_STOP"]
  const SDLK_MEDIA_STOP u32

  @[link_name="SDLK_MENU"]
  const SDLK_MENU u32

  @[link_name="SDLK_MINUS"]
  const SDLK_MINUS u32

  @[link_name="SDLK_MODE"]
  const SDLK_MODE u32

  @[link_name="SDLK_MULTI_KEY_COMPOSE"]
  const SDLK_MULTI_KEY_COMPOSE u32

  @[link_name="SDLK_MUTE"]
  const SDLK_MUTE u32

  @[link_name="SDLK_N"]
  const SDLK_N u32

  @[link_name="SDLK_NUMLOCKCLEAR"]
  const SDLK_NUMLOCKCLEAR u32

  @[link_name="SDLK_O"]
  const SDLK_O u32

  @[link_name="SDLK_OPER"]
  const SDLK_OPER u32

  @[link_name="SDLK_OUT"]
  const SDLK_OUT u32

  @[link_name="SDLK_P"]
  const SDLK_P u32

  @[link_name="SDLK_PAGEDOWN"]
  const SDLK_PAGEDOWN u32

  @[link_name="SDLK_PAGEUP"]
  const SDLK_PAGEUP u32

  @[link_name="SDLK_PASTE"]
  const SDLK_PASTE u32

  @[link_name="SDLK_PAUSE"]
  const SDLK_PAUSE u32

  @[link_name="SDLK_PERCENT"]
  const SDLK_PERCENT u32

  @[link_name="SDLK_PERIOD"]
  const SDLK_PERIOD u32

  @[link_name="SDLK_PIPE"]
  const SDLK_PIPE u32

  @[link_name="SDLK_PLUS"]
  const SDLK_PLUS u32

  @[link_name="SDLK_PLUSMINUS"]
  const SDLK_PLUSMINUS u32

  @[link_name="SDLK_POWER"]
  const SDLK_POWER u32

  @[link_name="SDLK_PRINTSCREEN"]
  const SDLK_PRINTSCREEN u32

  @[link_name="SDLK_PRIOR"]
  const SDLK_PRIOR u32

  @[link_name="SDLK_Q"]
  const SDLK_Q u32

  @[link_name="SDLK_QUESTION"]
  const SDLK_QUESTION u32

  @[link_name="SDLK_R"]
  const SDLK_R u32

  @[link_name="SDLK_RALT"]
  const SDLK_RALT u32

  @[link_name="SDLK_RCTRL"]
  const SDLK_RCTRL u32

  @[link_name="SDLK_RETURN"]
  const SDLK_RETURN u32

  @[link_name="SDLK_RETURN2"]
  const SDLK_RETURN2 u32

  @[link_name="SDLK_RGUI"]
  const SDLK_RGUI u32

  @[link_name="SDLK_RHYPER"]
  const SDLK_RHYPER u32

  @[link_name="SDLK_RIGHT"]
  const SDLK_RIGHT u32

  @[link_name="SDLK_RIGHTBRACE"]
  const SDLK_RIGHTBRACE u32

  @[link_name="SDLK_RIGHTBRACKET"]
  const SDLK_RIGHTBRACKET u32

  @[link_name="SDLK_RIGHTPAREN"]
  const SDLK_RIGHTPAREN u32

  @[link_name="SDLK_RMETA"]
  const SDLK_RMETA u32

  @[link_name="SDLK_RSHIFT"]
  const SDLK_RSHIFT u32

  @[link_name="SDLK_S"]
  const SDLK_S u32

  @[link_name="SDLK_SCANCODE_MASK"]
  const SDLK_SCANCODE_MASK u32

  @[link_name="SDLK_SCROLLLOCK"]
  const SDLK_SCROLLLOCK u32

  @[link_name="SDLK_SELECT"]
  const SDLK_SELECT u32

  @[link_name="SDLK_SEMICOLON"]
  const SDLK_SEMICOLON u32

  @[link_name="SDLK_SEPARATOR"]
  const SDLK_SEPARATOR u32

  @[link_name="SDLK_SLASH"]
  const SDLK_SLASH u32

  @[link_name="SDLK_SLEEP"]
  const SDLK_SLEEP u32

  @[link_name="SDLK_SOFTLEFT"]
  const SDLK_SOFTLEFT u32

  @[link_name="SDLK_SOFTRIGHT"]
  const SDLK_SOFTRIGHT u32

  @[link_name="SDLK_SPACE"]
  const SDLK_SPACE u32

  @[link_name="SDLK_STOP"]
  const SDLK_STOP u32

  @[link_name="SDLK_SYSREQ"]
  const SDLK_SYSREQ u32

  @[link_name="SDLK_T"]
  const SDLK_T u32

  @[link_name="SDLK_TAB"]
  const SDLK_TAB u32

  @[link_name="SDLK_THOUSANDSSEPARATOR"]
  const SDLK_THOUSANDSSEPARATOR u32

  @[link_name="SDLK_TILDE"]
  const SDLK_TILDE u32

  @[link_name="SDLK_U"]
  const SDLK_U u32

  @[link_name="SDLK_UNDERSCORE"]
  const SDLK_UNDERSCORE u32

  @[link_name="SDLK_UNDO"]
  const SDLK_UNDO u32

  @[link_name="SDLK_UNKNOWN"]
  const SDLK_UNKNOWN u32

  @[link_name="SDLK_UP"]
  const SDLK_UP u32

  @[link_name="SDLK_V"]
  const SDLK_V u32

  @[link_name="SDLK_VOLUMEDOWN"]
  const SDLK_VOLUMEDOWN u32

  @[link_name="SDLK_VOLUMEUP"]
  const SDLK_VOLUMEUP u32

  @[link_name="SDLK_W"]
  const SDLK_W u32

  @[link_name="SDLK_WAKE"]
  const SDLK_WAKE u32

  @[link_name="SDLK_X"]
  const SDLK_X u32

  @[link_name="SDLK_Y"]
  const SDLK_Y u32

  @[link_name="SDLK_Z"]
  const SDLK_Z u32

  @[link_name="SDL_ALPHA_OPAQUE"]
  const ALPHA_OPAQUE i32

  @[link_name="SDL_ALPHA_OPAQUE_FLOAT"]
  const ALPHA_OPAQUE_FLOAT f32

  @[link_name="SDL_ALPHA_TRANSPARENT"]
  const ALPHA_TRANSPARENT i32

  @[link_name="SDL_ALPHA_TRANSPARENT_FLOAT"]
  const ALPHA_TRANSPARENT_FLOAT f32

  @[link_name="SDL_AUDIO_DEVICE_DEFAULT_PLAYBACK"]
  const AUDIO_DEVICE_DEFAULT_PLAYBACK u32

  @[link_name="SDL_AUDIO_DEVICE_DEFAULT_RECORDING"]
  const AUDIO_DEVICE_DEFAULT_RECORDING u32

  @[link_name="SDL_AUDIO_MASK_BIG_ENDIAN"]
  const AUDIO_MASK_BIG_ENDIAN u32

  @[link_name="SDL_AUDIO_MASK_BITSIZE"]
  const AUDIO_MASK_BITSIZE u32

  @[link_name="SDL_AUDIO_MASK_FLOAT"]
  const AUDIO_MASK_FLOAT u32

  @[link_name="SDL_AUDIO_MASK_SIGNED"]
  const AUDIO_MASK_SIGNED u32

  @[link_name="SDL_BIG_ENDIAN"]
  const BIG_ENDIAN i32

  @[link_name="SDL_BLENDMODE_ADD"]
  const BLENDMODE_ADD u32

  @[link_name="SDL_BLENDMODE_ADD_PREMULTIPLIED"]
  const BLENDMODE_ADD_PREMULTIPLIED u32

  @[link_name="SDL_BLENDMODE_BLEND"]
  const BLENDMODE_BLEND u32

  @[link_name="SDL_BLENDMODE_BLEND_PREMULTIPLIED"]
  const BLENDMODE_BLEND_PREMULTIPLIED u32

  @[link_name="SDL_BLENDMODE_INVALID"]
  const BLENDMODE_INVALID u32

  @[link_name="SDL_BLENDMODE_MOD"]
  const BLENDMODE_MOD u32

  @[link_name="SDL_BLENDMODE_MUL"]
  const BLENDMODE_MUL u32

  @[link_name="SDL_BLENDMODE_NONE"]
  const BLENDMODE_NONE u32

  @[link_name="SDL_BUTTON_LEFT"]
  const BUTTON_LEFT i32

  @[link_name="SDL_BUTTON_LMASK"]
  const BUTTON_LMASK u32

  @[link_name="SDL_BUTTON_MIDDLE"]
  const BUTTON_MIDDLE i32

  @[link_name="SDL_BUTTON_MMASK"]
  const BUTTON_MMASK u32

  @[link_name="SDL_BUTTON_RIGHT"]
  const BUTTON_RIGHT i32

  @[link_name="SDL_BUTTON_RMASK"]
  const BUTTON_RMASK u32

  @[link_name="SDL_BUTTON_X1"]
  const BUTTON_X1 i32

  @[link_name="SDL_BUTTON_X1MASK"]
  const BUTTON_X1MASK u32

  @[link_name="SDL_BUTTON_X2"]
  const BUTTON_X2 i32

  @[link_name="SDL_BUTTON_X2MASK"]
  const BUTTON_X2MASK u32

  @[link_name="SDL_BYTEORDER"]
  const BYTEORDER i32

  @[link_name="SDL_CACHELINE_SIZE"]
  const CACHELINE_SIZE i32

  @[link_name="SDL_FLT_EPSILON"]
  const FLT_EPSILON f32

  @[link_name="SDL_GLOB_CASEINSENSITIVE"]
  const GLOB_CASEINSENSITIVE u32

  @[link_name="SDL_GL_CONTEXT_DEBUG_FLAG"]
  const GL_CONTEXT_DEBUG_FLAG i32

  @[link_name="SDL_GL_CONTEXT_FORWARD_COMPATIBLE_FLAG"]
  const GL_CONTEXT_FORWARD_COMPATIBLE_FLAG i32

  @[link_name="SDL_GL_CONTEXT_PROFILE_COMPATIBILITY"]
  const GL_CONTEXT_PROFILE_COMPATIBILITY i32

  @[link_name="SDL_GL_CONTEXT_PROFILE_CORE"]
  const GL_CONTEXT_PROFILE_CORE i32

  @[link_name="SDL_GL_CONTEXT_PROFILE_ES"]
  const GL_CONTEXT_PROFILE_ES i32

  @[link_name="SDL_GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH"]
  const GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH i32

  @[link_name="SDL_GL_CONTEXT_RELEASE_BEHAVIOR_NONE"]
  const GL_CONTEXT_RELEASE_BEHAVIOR_NONE i32

  @[link_name="SDL_GL_CONTEXT_RESET_ISOLATION_FLAG"]
  const GL_CONTEXT_RESET_ISOLATION_FLAG i32

  @[link_name="SDL_GL_CONTEXT_RESET_LOSE_CONTEXT"]
  const GL_CONTEXT_RESET_LOSE_CONTEXT i32

  @[link_name="SDL_GL_CONTEXT_RESET_NO_NOTIFICATION"]
  const GL_CONTEXT_RESET_NO_NOTIFICATION i32

  @[link_name="SDL_GL_CONTEXT_ROBUST_ACCESS_FLAG"]
  const GL_CONTEXT_ROBUST_ACCESS_FLAG i32

  @[link_name="SDL_GPU_BUFFERUSAGE_COMPUTE_STORAGE_READ"]
  const GPU_BUFFERUSAGE_COMPUTE_STORAGE_READ u32

  @[link_name="SDL_GPU_BUFFERUSAGE_COMPUTE_STORAGE_WRITE"]
  const GPU_BUFFERUSAGE_COMPUTE_STORAGE_WRITE u32

  @[link_name="SDL_GPU_BUFFERUSAGE_GRAPHICS_STORAGE_READ"]
  const GPU_BUFFERUSAGE_GRAPHICS_STORAGE_READ u32

  @[link_name="SDL_GPU_BUFFERUSAGE_INDEX"]
  const GPU_BUFFERUSAGE_INDEX u32

  @[link_name="SDL_GPU_BUFFERUSAGE_INDIRECT"]
  const GPU_BUFFERUSAGE_INDIRECT u32

  @[link_name="SDL_GPU_BUFFERUSAGE_VERTEX"]
  const GPU_BUFFERUSAGE_VERTEX u32

  @[link_name="SDL_GPU_COLORCOMPONENT_A"]
  const GPU_COLORCOMPONENT_A u32

  @[link_name="SDL_GPU_COLORCOMPONENT_B"]
  const GPU_COLORCOMPONENT_B u32

  @[link_name="SDL_GPU_COLORCOMPONENT_G"]
  const GPU_COLORCOMPONENT_G u32

  @[link_name="SDL_GPU_COLORCOMPONENT_R"]
  const GPU_COLORCOMPONENT_R u32

  @[link_name="SDL_GPU_SHADERFORMAT_DXBC"]
  const GPU_SHADERFORMAT_DXBC u32

  @[link_name="SDL_GPU_SHADERFORMAT_DXIL"]
  const GPU_SHADERFORMAT_DXIL u32

  @[link_name="SDL_GPU_SHADERFORMAT_INVALID"]
  const GPU_SHADERFORMAT_INVALID i32

  @[link_name="SDL_GPU_SHADERFORMAT_METALLIB"]
  const GPU_SHADERFORMAT_METALLIB u32

  @[link_name="SDL_GPU_SHADERFORMAT_MSL"]
  const GPU_SHADERFORMAT_MSL u32

  @[link_name="SDL_GPU_SHADERFORMAT_PRIVATE"]
  const GPU_SHADERFORMAT_PRIVATE u32

  @[link_name="SDL_GPU_SHADERFORMAT_SPIRV"]
  const GPU_SHADERFORMAT_SPIRV u32

  @[link_name="SDL_GPU_TEXTUREUSAGE_COLOR_TARGET"]
  const GPU_TEXTUREUSAGE_COLOR_TARGET u32

  @[link_name="SDL_GPU_TEXTUREUSAGE_COMPUTE_STORAGE_READ"]
  const GPU_TEXTUREUSAGE_COMPUTE_STORAGE_READ u32

  @[link_name="SDL_GPU_TEXTUREUSAGE_COMPUTE_STORAGE_SIMULTANEOUS_READ_WRITE"]
  const GPU_TEXTUREUSAGE_COMPUTE_STORAGE_SIMULTANEOUS_READ_WRITE u32

  @[link_name="SDL_GPU_TEXTUREUSAGE_COMPUTE_STORAGE_WRITE"]
  const GPU_TEXTUREUSAGE_COMPUTE_STORAGE_WRITE u32

  @[link_name="SDL_GPU_TEXTUREUSAGE_DEPTH_STENCIL_TARGET"]
  const GPU_TEXTUREUSAGE_DEPTH_STENCIL_TARGET u32

  @[link_name="SDL_GPU_TEXTUREUSAGE_GRAPHICS_STORAGE_READ"]
  const GPU_TEXTUREUSAGE_GRAPHICS_STORAGE_READ u32

  @[link_name="SDL_GPU_TEXTUREUSAGE_SAMPLER"]
  const GPU_TEXTUREUSAGE_SAMPLER u32

  @[link_name="SDL_HAPTIC_AUTOCENTER"]
  const HAPTIC_AUTOCENTER u32

  @[link_name="SDL_HAPTIC_CARTESIAN"]
  const HAPTIC_CARTESIAN i32

  @[link_name="SDL_HAPTIC_CONSTANT"]
  const HAPTIC_CONSTANT u32

  @[link_name="SDL_HAPTIC_CUSTOM"]
  const HAPTIC_CUSTOM u32

  @[link_name="SDL_HAPTIC_DAMPER"]
  const HAPTIC_DAMPER u32

  @[link_name="SDL_HAPTIC_FRICTION"]
  const HAPTIC_FRICTION u32

  @[link_name="SDL_HAPTIC_GAIN"]
  const HAPTIC_GAIN u32

  @[link_name="SDL_HAPTIC_INERTIA"]
  const HAPTIC_INERTIA u32

  @[link_name="SDL_HAPTIC_INFINITY"]
  const HAPTIC_INFINITY u32

  @[link_name="SDL_HAPTIC_LEFTRIGHT"]
  const HAPTIC_LEFTRIGHT u32

  @[link_name="SDL_HAPTIC_PAUSE"]
  const HAPTIC_PAUSE u32

  @[link_name="SDL_HAPTIC_POLAR"]
  const HAPTIC_POLAR i32

  @[link_name="SDL_HAPTIC_RAMP"]
  const HAPTIC_RAMP u32

  @[link_name="SDL_HAPTIC_RESERVED1"]
  const HAPTIC_RESERVED1 u32

  @[link_name="SDL_HAPTIC_RESERVED2"]
  const HAPTIC_RESERVED2 u32

  @[link_name="SDL_HAPTIC_RESERVED3"]
  const HAPTIC_RESERVED3 u32

  @[link_name="SDL_HAPTIC_SAWTOOTHDOWN"]
  const HAPTIC_SAWTOOTHDOWN u32

  @[link_name="SDL_HAPTIC_SAWTOOTHUP"]
  const HAPTIC_SAWTOOTHUP u32

  @[link_name="SDL_HAPTIC_SINE"]
  const HAPTIC_SINE u32

  @[link_name="SDL_HAPTIC_SPHERICAL"]
  const HAPTIC_SPHERICAL i32

  @[link_name="SDL_HAPTIC_SPRING"]
  const HAPTIC_SPRING u32

  @[link_name="SDL_HAPTIC_SQUARE"]
  const HAPTIC_SQUARE u32

  @[link_name="SDL_HAPTIC_STATUS"]
  const HAPTIC_STATUS u32

  @[link_name="SDL_HAPTIC_STEERING_AXIS"]
  const HAPTIC_STEERING_AXIS i32

  @[link_name="SDL_HAPTIC_TRIANGLE"]
  const HAPTIC_TRIANGLE u32

  @[link_name="SDL_HAT_CENTERED"]
  const HAT_CENTERED u32

  @[link_name="SDL_HAT_DOWN"]
  const HAT_DOWN u32

  @[link_name="SDL_HAT_LEFT"]
  const HAT_LEFT u32

  @[link_name="SDL_HAT_LEFTDOWN"]
  const HAT_LEFTDOWN u32

  @[link_name="SDL_HAT_LEFTUP"]
  const HAT_LEFTUP u32

  @[link_name="SDL_HAT_RIGHT"]
  const HAT_RIGHT u32

  @[link_name="SDL_HAT_RIGHTDOWN"]
  const HAT_RIGHTDOWN u32

  @[link_name="SDL_HAT_RIGHTUP"]
  const HAT_RIGHTUP u32

  @[link_name="SDL_HAT_UP"]
  const HAT_UP u32

  @[link_name="SDL_INIT_AUDIO"]
  const INIT_AUDIO u32

  @[link_name="SDL_INIT_CAMERA"]
  const INIT_CAMERA u32

  @[link_name="SDL_INIT_EVENTS"]
  const INIT_EVENTS u32

  @[link_name="SDL_INIT_GAMEPAD"]
  const INIT_GAMEPAD u32

  @[link_name="SDL_INIT_HAPTIC"]
  const INIT_HAPTIC u32

  @[link_name="SDL_INIT_JOYSTICK"]
  const INIT_JOYSTICK u32

  @[link_name="SDL_INIT_SENSOR"]
  const INIT_SENSOR u32

  @[link_name="SDL_INIT_VIDEO"]
  const INIT_VIDEO u32

  @[link_name="SDL_JOYSTICK_AXIS_MAX"]
  const JOYSTICK_AXIS_MAX i32

  @[link_name="SDL_JOYSTICK_AXIS_MIN"]
  const JOYSTICK_AXIS_MIN i32

  @[link_name="SDL_KMOD_ALT"]
  const KMOD_ALT u32

  @[link_name="SDL_KMOD_CAPS"]
  const KMOD_CAPS u32

  @[link_name="SDL_KMOD_CTRL"]
  const KMOD_CTRL u32

  @[link_name="SDL_KMOD_GUI"]
  const KMOD_GUI u32

  @[link_name="SDL_KMOD_LALT"]
  const KMOD_LALT u32

  @[link_name="SDL_KMOD_LCTRL"]
  const KMOD_LCTRL u32

  @[link_name="SDL_KMOD_LEVEL5"]
  const KMOD_LEVEL5 u32

  @[link_name="SDL_KMOD_LGUI"]
  const KMOD_LGUI u32

  @[link_name="SDL_KMOD_LSHIFT"]
  const KMOD_LSHIFT u32

  @[link_name="SDL_KMOD_MODE"]
  const KMOD_MODE u32

  @[link_name="SDL_KMOD_NONE"]
  const KMOD_NONE u32

  @[link_name="SDL_KMOD_NUM"]
  const KMOD_NUM u32

  @[link_name="SDL_KMOD_RALT"]
  const KMOD_RALT u32

  @[link_name="SDL_KMOD_RCTRL"]
  const KMOD_RCTRL u32

  @[link_name="SDL_KMOD_RGUI"]
  const KMOD_RGUI u32

  @[link_name="SDL_KMOD_RSHIFT"]
  const KMOD_RSHIFT u32

  @[link_name="SDL_KMOD_SCROLL"]
  const KMOD_SCROLL u32

  @[link_name="SDL_KMOD_SHIFT"]
  const KMOD_SHIFT u32

  @[link_name="SDL_LIL_ENDIAN"]
  const LIL_ENDIAN i32

  @[link_name="SDL_MAJOR_VERSION"]
  const MAJOR_VERSION i32

  @[link_name="SDL_MAX_SINT16"]
  const MAX_SINT16 i16

  @[link_name="SDL_MAX_SINT32"]
  const MAX_SINT32 i32

  @[link_name="SDL_MAX_SINT64"]
  const MAX_SINT64 i64

  @[link_name="SDL_MAX_SINT8"]
  const MAX_SINT8 i8

  @[link_name="SDL_MAX_TIME"]
  const MAX_TIME i64

  @[link_name="SDL_MAX_UINT16"]
  const MAX_UINT16 u16

  @[link_name="SDL_MAX_UINT32"]
  const MAX_UINT32 u32

  @[link_name="SDL_MAX_UINT64"]
  const MAX_UINT64 u64

  @[link_name="SDL_MAX_UINT8"]
  const MAX_UINT8 u8

  @[link_name="SDL_MESSAGEBOX_BUTTONS_LEFT_TO_RIGHT"]
  const MESSAGEBOX_BUTTONS_LEFT_TO_RIGHT u32

  @[link_name="SDL_MESSAGEBOX_BUTTONS_RIGHT_TO_LEFT"]
  const MESSAGEBOX_BUTTONS_RIGHT_TO_LEFT u32

  @[link_name="SDL_MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT"]
  const MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT u32

  @[link_name="SDL_MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT"]
  const MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT u32

  @[link_name="SDL_MESSAGEBOX_ERROR"]
  const MESSAGEBOX_ERROR u32

  @[link_name="SDL_MESSAGEBOX_INFORMATION"]
  const MESSAGEBOX_INFORMATION u32

  @[link_name="SDL_MESSAGEBOX_WARNING"]
  const MESSAGEBOX_WARNING u32

  @[link_name="SDL_MICRO_VERSION"]
  const MICRO_VERSION i32

  @[link_name="SDL_MINOR_VERSION"]
  const MINOR_VERSION i32

  @[link_name="SDL_MIN_SINT16"]
  const MIN_SINT16 i16

  @[link_name="SDL_MIN_SINT32"]
  const MIN_SINT32 i32

  @[link_name="SDL_MIN_SINT64"]
  const MIN_SINT64 i64

  @[link_name="SDL_MIN_SINT8"]
  const MIN_SINT8 i8

  @[link_name="SDL_MIN_TIME"]
  const MIN_TIME i64

  @[link_name="SDL_MIN_UINT16"]
  const MIN_UINT16 u16

  @[link_name="SDL_MIN_UINT32"]
  const MIN_UINT32 u32

  @[link_name="SDL_MIN_UINT64"]
  const MIN_UINT64 u64

  @[link_name="SDL_MIN_UINT8"]
  const MIN_UINT8 u8

  @[link_name="SDL_MOUSE_TOUCHID"]
  const MOUSE_TOUCHID u64

  @[link_name="SDL_MS_PER_SECOND"]
  const MS_PER_SECOND i32

  @[link_name="SDL_NS_PER_MS"]
  const NS_PER_MS i32

  @[link_name="SDL_NS_PER_SECOND"]
  const NS_PER_SECOND i64

  @[link_name="SDL_NS_PER_US"]
  const NS_PER_US i32

  @[link_name="SDL_PEN_INPUT_BUTTON_1"]
  const PEN_INPUT_BUTTON_1 u32

  @[link_name="SDL_PEN_INPUT_BUTTON_2"]
  const PEN_INPUT_BUTTON_2 u32

  @[link_name="SDL_PEN_INPUT_BUTTON_3"]
  const PEN_INPUT_BUTTON_3 u32

  @[link_name="SDL_PEN_INPUT_BUTTON_4"]
  const PEN_INPUT_BUTTON_4 u32

  @[link_name="SDL_PEN_INPUT_BUTTON_5"]
  const PEN_INPUT_BUTTON_5 u32

  @[link_name="SDL_PEN_INPUT_DOWN"]
  const PEN_INPUT_DOWN u32

  @[link_name="SDL_PEN_INPUT_ERASER_TIP"]
  const PEN_INPUT_ERASER_TIP u32

  @[link_name="SDL_PEN_MOUSEID"]
  const PEN_MOUSEID u32

  @[link_name="SDL_PEN_TOUCHID"]
  const PEN_TOUCHID u64

  @[link_name="SDL_PI_D"]
  const PI_D f64

  @[link_name="SDL_PI_F"]
  const PI_F f32

  @[link_name="SDL_RENDERER_VSYNC_ADAPTIVE"]
  const RENDERER_VSYNC_ADAPTIVE i32

  @[link_name="SDL_RENDERER_VSYNC_DISABLED"]
  const RENDERER_VSYNC_DISABLED i32

  @[link_name="SDL_STANDARD_GRAVITY"]
  const STANDARD_GRAVITY f32

  @[link_name="SDL_SURFACE_LOCKED"]
  const SURFACE_LOCKED u32

  @[link_name="SDL_SURFACE_LOCK_NEEDED"]
  const SURFACE_LOCK_NEEDED u32

  @[link_name="SDL_SURFACE_PREALLOCATED"]
  const SURFACE_PREALLOCATED u32

  @[link_name="SDL_SURFACE_SIMD_ALIGNED"]
  const SURFACE_SIMD_ALIGNED u32

  @[link_name="SDL_TOUCH_MOUSEID"]
  const TOUCH_MOUSEID u32

  @[link_name="SDL_TRAYENTRY_BUTTON"]
  const TRAYENTRY_BUTTON u32

  @[link_name="SDL_TRAYENTRY_CHECKBOX"]
  const TRAYENTRY_CHECKBOX u32

  @[link_name="SDL_TRAYENTRY_CHECKED"]
  const TRAYENTRY_CHECKED u32

  @[link_name="SDL_TRAYENTRY_DISABLED"]
  const TRAYENTRY_DISABLED u32

  @[link_name="SDL_TRAYENTRY_SUBMENU"]
  const TRAYENTRY_SUBMENU u32

  @[link_name="SDL_US_PER_SECOND"]
  const US_PER_SECOND i32

  @[link_name="SDL_VERSION"]
  const VERSION i32

  @[link_name="SDL_WINDOWPOS_CENTERED"]
  const WINDOWPOS_CENTERED u32

  @[link_name="SDL_WINDOWPOS_CENTERED_MASK"]
  const WINDOWPOS_CENTERED_MASK u32

  @[link_name="SDL_WINDOWPOS_UNDEFINED"]
  const WINDOWPOS_UNDEFINED u32

  @[link_name="SDL_WINDOWPOS_UNDEFINED_MASK"]
  const WINDOWPOS_UNDEFINED_MASK u32

  @[link_name="SDL_WINDOW_ALWAYS_ON_TOP"]
  const WINDOW_ALWAYS_ON_TOP u64

  @[link_name="SDL_WINDOW_BORDERLESS"]
  const WINDOW_BORDERLESS u64

  @[link_name="SDL_WINDOW_EXTERNAL"]
  const WINDOW_EXTERNAL u64

  @[link_name="SDL_WINDOW_FULLSCREEN"]
  const WINDOW_FULLSCREEN u64

  @[link_name="SDL_WINDOW_HIDDEN"]
  const WINDOW_HIDDEN u64

  @[link_name="SDL_WINDOW_HIGH_PIXEL_DENSITY"]
  const WINDOW_HIGH_PIXEL_DENSITY u64

  @[link_name="SDL_WINDOW_INPUT_FOCUS"]
  const WINDOW_INPUT_FOCUS u64

  @[link_name="SDL_WINDOW_KEYBOARD_GRABBED"]
  const WINDOW_KEYBOARD_GRABBED u64

  @[link_name="SDL_WINDOW_MAXIMIZED"]
  const WINDOW_MAXIMIZED u64

  @[link_name="SDL_WINDOW_METAL"]
  const WINDOW_METAL u64

  @[link_name="SDL_WINDOW_MINIMIZED"]
  const WINDOW_MINIMIZED u64

  @[link_name="SDL_WINDOW_MODAL"]
  const WINDOW_MODAL u64

  @[link_name="SDL_WINDOW_MOUSE_CAPTURE"]
  const WINDOW_MOUSE_CAPTURE u64

  @[link_name="SDL_WINDOW_MOUSE_FOCUS"]
  const WINDOW_MOUSE_FOCUS u64

  @[link_name="SDL_WINDOW_MOUSE_GRABBED"]
  const WINDOW_MOUSE_GRABBED u64

  @[link_name="SDL_WINDOW_MOUSE_RELATIVE_MODE"]
  const WINDOW_MOUSE_RELATIVE_MODE u64

  @[link_name="SDL_WINDOW_NOT_FOCUSABLE"]
  const WINDOW_NOT_FOCUSABLE u64

  @[link_name="SDL_WINDOW_OCCLUDED"]
  const WINDOW_OCCLUDED u64

  @[link_name="SDL_WINDOW_OPENGL"]
  const WINDOW_OPENGL u64

  @[link_name="SDL_WINDOW_POPUP_MENU"]
  const WINDOW_POPUP_MENU u64

  @[link_name="SDL_WINDOW_RESIZABLE"]
  const WINDOW_RESIZABLE u64

  @[link_name="SDL_WINDOW_SURFACE_VSYNC_ADAPTIVE"]
  const WINDOW_SURFACE_VSYNC_ADAPTIVE i32

  @[link_name="SDL_WINDOW_SURFACE_VSYNC_DISABLED"]
  const WINDOW_SURFACE_VSYNC_DISABLED i32

  @[link_name="SDL_WINDOW_TOOLTIP"]
  const WINDOW_TOOLTIP u64

  @[link_name="SDL_WINDOW_TRANSPARENT"]
  const WINDOW_TRANSPARENT u64

  @[link_name="SDL_WINDOW_UTILITY"]
  const WINDOW_UTILITY u64

  @[link_name="SDL_WINDOW_VULKAN"]
  const WINDOW_VULKAN u64
}

const HINT_ALLOW_ALT_TAB_WHILE_GRABBED string = "SDL_ALLOW_ALT_TAB_WHILE_GRABBED"
const HINT_ANDROID_ALLOW_RECREATE_ACTIVITY string = "SDL_ANDROID_ALLOW_RECREATE_ACTIVITY"
const HINT_ANDROID_BLOCK_ON_PAUSE string = "SDL_ANDROID_BLOCK_ON_PAUSE"
const HINT_ANDROID_LOW_LATENCY_AUDIO string = "SDL_ANDROID_LOW_LATENCY_AUDIO"
const HINT_ANDROID_TRAP_BACK_BUTTON string = "SDL_ANDROID_TRAP_BACK_BUTTON"
const HINT_APPLE_TV_CONTROLLER_UI_EVENTS string = "SDL_APPLE_TV_CONTROLLER_UI_EVENTS"
const HINT_APPLE_TV_REMOTE_ALLOW_ROTATION string = "SDL_APPLE_TV_REMOTE_ALLOW_ROTATION"
const HINT_APP_ID string = "SDL_APP_ID"
const HINT_APP_NAME string = "SDL_APP_NAME"
const HINT_ASSERT string = "SDL_ASSERT"
const HINT_AUDIO_ALSA_DEFAULT_DEVICE string = "SDL_AUDIO_ALSA_DEFAULT_DEVICE"
const HINT_AUDIO_ALSA_DEFAULT_PLAYBACK_DEVICE string = "SDL_AUDIO_ALSA_DEFAULT_PLAYBACK_DEVICE"
const HINT_AUDIO_ALSA_DEFAULT_RECORDING_DEVICE string = "SDL_AUDIO_ALSA_DEFAULT_RECORDING_DEVICE"
const HINT_AUDIO_CATEGORY string = "SDL_AUDIO_CATEGORY"
const HINT_AUDIO_CHANNELS string = "SDL_AUDIO_CHANNELS"
const HINT_AUDIO_DEVICE_APP_ICON_NAME string = "SDL_AUDIO_DEVICE_APP_ICON_NAME"
const HINT_AUDIO_DEVICE_SAMPLE_FRAMES string = "SDL_AUDIO_DEVICE_SAMPLE_FRAMES"
const HINT_AUDIO_DEVICE_STREAM_NAME string = "SDL_AUDIO_DEVICE_STREAM_NAME"
const HINT_AUDIO_DEVICE_STREAM_ROLE string = "SDL_AUDIO_DEVICE_STREAM_ROLE"
const HINT_AUDIO_DISK_INPUT_FILE string = "SDL_AUDIO_DISK_INPUT_FILE"
const HINT_AUDIO_DISK_OUTPUT_FILE string = "SDL_AUDIO_DISK_OUTPUT_FILE"
const HINT_AUDIO_DISK_TIMESCALE string = "SDL_AUDIO_DISK_TIMESCALE"
const HINT_AUDIO_DRIVER string = "SDL_AUDIO_DRIVER"
const HINT_AUDIO_DUMMY_TIMESCALE string = "SDL_AUDIO_DUMMY_TIMESCALE"
const HINT_AUDIO_FORMAT string = "SDL_AUDIO_FORMAT"
const HINT_AUDIO_FREQUENCY string = "SDL_AUDIO_FREQUENCY"
const HINT_AUDIO_INCLUDE_MONITORS string = "SDL_AUDIO_INCLUDE_MONITORS"
const HINT_AUTO_UPDATE_JOYSTICKS string = "SDL_AUTO_UPDATE_JOYSTICKS"
const HINT_AUTO_UPDATE_SENSORS string = "SDL_AUTO_UPDATE_SENSORS"
const HINT_BMP_SAVE_LEGACY_FORMAT string = "SDL_BMP_SAVE_LEGACY_FORMAT"
const HINT_CAMERA_DRIVER string = "SDL_CAMERA_DRIVER"
const HINT_CPU_FEATURE_MASK string = "SDL_CPU_FEATURE_MASK"
const HINT_DISPLAY_USABLE_BOUNDS string = "SDL_DISPLAY_USABLE_BOUNDS"
const HINT_EGL_LIBRARY string = "SDL_EGL_LIBRARY"
const HINT_EMSCRIPTEN_ASYNCIFY string = "SDL_EMSCRIPTEN_ASYNCIFY"
const HINT_EMSCRIPTEN_CANVAS_SELECTOR string = "SDL_EMSCRIPTEN_CANVAS_SELECTOR"
const HINT_EMSCRIPTEN_KEYBOARD_ELEMENT string = "SDL_EMSCRIPTEN_KEYBOARD_ELEMENT"
const HINT_ENABLE_SCREEN_KEYBOARD string = "SDL_ENABLE_SCREEN_KEYBOARD"
const HINT_EVDEV_DEVICES string = "SDL_EVDEV_DEVICES"
const HINT_EVENT_LOGGING string = "SDL_EVENT_LOGGING"
const HINT_FILE_DIALOG_DRIVER string = "SDL_FILE_DIALOG_DRIVER"
const HINT_FORCE_RAISEWINDOW string = "SDL_FORCE_RAISEWINDOW"
const HINT_FRAMEBUFFER_ACCELERATION string = "SDL_FRAMEBUFFER_ACCELERATION"
const HINT_GAMECONTROLLERCONFIG string = "SDL_GAMECONTROLLERCONFIG"
const HINT_GAMECONTROLLERCONFIG_FILE string = "SDL_GAMECONTROLLERCONFIG_FILE"
const HINT_GAMECONTROLLERTYPE string = "SDL_GAMECONTROLLERTYPE"
const HINT_GAMECONTROLLER_IGNORE_DEVICES string = "SDL_GAMECONTROLLER_IGNORE_DEVICES"
const HINT_GAMECONTROLLER_IGNORE_DEVICES_EXCEPT string = "SDL_GAMECONTROLLER_IGNORE_DEVICES_EXCEPT"
const HINT_GAMECONTROLLER_SENSOR_FUSION string = "SDL_GAMECONTROLLER_SENSOR_FUSION"
const HINT_GDK_TEXTINPUT_DEFAULT_TEXT string = "SDL_GDK_TEXTINPUT_DEFAULT_TEXT"
const HINT_GDK_TEXTINPUT_DESCRIPTION string = "SDL_GDK_TEXTINPUT_DESCRIPTION"
const HINT_GDK_TEXTINPUT_MAX_LENGTH string = "SDL_GDK_TEXTINPUT_MAX_LENGTH"
const HINT_GDK_TEXTINPUT_SCOPE string = "SDL_GDK_TEXTINPUT_SCOPE"
const HINT_GDK_TEXTINPUT_TITLE string = "SDL_GDK_TEXTINPUT_TITLE"
const HINT_GPU_DRIVER string = "SDL_GPU_DRIVER"
const HINT_HIDAPI_ENUMERATE_ONLY_CONTROLLERS string = "SDL_HIDAPI_ENUMERATE_ONLY_CONTROLLERS"
const HINT_HIDAPI_IGNORE_DEVICES string = "SDL_HIDAPI_IGNORE_DEVICES"
const HINT_HIDAPI_LIBUSB string = "SDL_HIDAPI_LIBUSB"
const HINT_HIDAPI_LIBUSB_WHITELIST string = "SDL_HIDAPI_LIBUSB_WHITELIST"
const HINT_HIDAPI_UDEV string = "SDL_HIDAPI_UDEV"
const HINT_IME_IMPLEMENTED_UI string = "SDL_IME_IMPLEMENTED_UI"
const HINT_IOS_HIDE_HOME_INDICATOR string = "SDL_IOS_HIDE_HOME_INDICATOR"
const HINT_JOYSTICK_ALLOW_BACKGROUND_EVENTS string = "SDL_JOYSTICK_ALLOW_BACKGROUND_EVENTS"
const HINT_JOYSTICK_ARCADESTICK_DEVICES string = "SDL_JOYSTICK_ARCADESTICK_DEVICES"
const HINT_JOYSTICK_ARCADESTICK_DEVICES_EXCLUDED string = "SDL_JOYSTICK_ARCADESTICK_DEVICES_EXCLUDED"
const HINT_JOYSTICK_BLACKLIST_DEVICES string = "SDL_JOYSTICK_BLACKLIST_DEVICES"
const HINT_JOYSTICK_BLACKLIST_DEVICES_EXCLUDED string = "SDL_JOYSTICK_BLACKLIST_DEVICES_EXCLUDED"
const HINT_JOYSTICK_DEVICE string = "SDL_JOYSTICK_DEVICE"
const HINT_JOYSTICK_DIRECTINPUT string = "SDL_JOYSTICK_DIRECTINPUT"
const HINT_JOYSTICK_ENHANCED_REPORTS string = "SDL_JOYSTICK_ENHANCED_REPORTS"
const HINT_JOYSTICK_FLIGHTSTICK_DEVICES string = "SDL_JOYSTICK_FLIGHTSTICK_DEVICES"
const HINT_JOYSTICK_FLIGHTSTICK_DEVICES_EXCLUDED string = "SDL_JOYSTICK_FLIGHTSTICK_DEVICES_EXCLUDED"
const HINT_JOYSTICK_GAMECUBE_DEVICES string = "SDL_JOYSTICK_GAMECUBE_DEVICES"
const HINT_JOYSTICK_GAMECUBE_DEVICES_EXCLUDED string = "SDL_JOYSTICK_GAMECUBE_DEVICES_EXCLUDED"
const HINT_JOYSTICK_GAMEINPUT string = "SDL_JOYSTICK_GAMEINPUT"
const HINT_JOYSTICK_HAPTIC_AXES string = "SDL_JOYSTICK_HAPTIC_AXES"
const HINT_JOYSTICK_HIDAPI string = "SDL_JOYSTICK_HIDAPI"
const HINT_JOYSTICK_HIDAPI_COMBINE_JOY_CONS string = "SDL_JOYSTICK_HIDAPI_COMBINE_JOY_CONS"
const HINT_JOYSTICK_HIDAPI_GAMECUBE string = "SDL_JOYSTICK_HIDAPI_GAMECUBE"
const HINT_JOYSTICK_HIDAPI_GAMECUBE_RUMBLE_BRAKE string = "SDL_JOYSTICK_HIDAPI_GAMECUBE_RUMBLE_BRAKE"
const HINT_JOYSTICK_HIDAPI_JOYCON_HOME_LED string = "SDL_JOYSTICK_HIDAPI_JOYCON_HOME_LED"
const HINT_JOYSTICK_HIDAPI_JOY_CONS string = "SDL_JOYSTICK_HIDAPI_JOY_CONS"
const HINT_JOYSTICK_HIDAPI_LUNA string = "SDL_JOYSTICK_HIDAPI_LUNA"
const HINT_JOYSTICK_HIDAPI_NINTENDO_CLASSIC string = "SDL_JOYSTICK_HIDAPI_NINTENDO_CLASSIC"
const HINT_JOYSTICK_HIDAPI_PS3 string = "SDL_JOYSTICK_HIDAPI_PS3"
const HINT_JOYSTICK_HIDAPI_PS3_SIXAXIS_DRIVER string = "SDL_JOYSTICK_HIDAPI_PS3_SIXAXIS_DRIVER"
const HINT_JOYSTICK_HIDAPI_PS4 string = "SDL_JOYSTICK_HIDAPI_PS4"
const HINT_JOYSTICK_HIDAPI_PS4_REPORT_INTERVAL string = "SDL_JOYSTICK_HIDAPI_PS4_REPORT_INTERVAL"
const HINT_JOYSTICK_HIDAPI_PS5 string = "SDL_JOYSTICK_HIDAPI_PS5"
const HINT_JOYSTICK_HIDAPI_PS5_PLAYER_LED string = "SDL_JOYSTICK_HIDAPI_PS5_PLAYER_LED"
const HINT_JOYSTICK_HIDAPI_SHIELD string = "SDL_JOYSTICK_HIDAPI_SHIELD"
const HINT_JOYSTICK_HIDAPI_STADIA string = "SDL_JOYSTICK_HIDAPI_STADIA"
const HINT_JOYSTICK_HIDAPI_STEAM string = "SDL_JOYSTICK_HIDAPI_STEAM"
const HINT_JOYSTICK_HIDAPI_STEAMDECK string = "SDL_JOYSTICK_HIDAPI_STEAMDECK"
const HINT_JOYSTICK_HIDAPI_STEAM_HOME_LED string = "SDL_JOYSTICK_HIDAPI_STEAM_HOME_LED"
const HINT_JOYSTICK_HIDAPI_STEAM_HORI string = "SDL_JOYSTICK_HIDAPI_STEAM_HORI"
const HINT_JOYSTICK_HIDAPI_SWITCH string = "SDL_JOYSTICK_HIDAPI_SWITCH"
const HINT_JOYSTICK_HIDAPI_SWITCH_HOME_LED string = "SDL_JOYSTICK_HIDAPI_SWITCH_HOME_LED"
const HINT_JOYSTICK_HIDAPI_SWITCH_PLAYER_LED string = "SDL_JOYSTICK_HIDAPI_SWITCH_PLAYER_LED"
const HINT_JOYSTICK_HIDAPI_VERTICAL_JOY_CONS string = "SDL_JOYSTICK_HIDAPI_VERTICAL_JOY_CONS"
const HINT_JOYSTICK_HIDAPI_WII string = "SDL_JOYSTICK_HIDAPI_WII"
const HINT_JOYSTICK_HIDAPI_WII_PLAYER_LED string = "SDL_JOYSTICK_HIDAPI_WII_PLAYER_LED"
const HINT_JOYSTICK_HIDAPI_XBOX string = "SDL_JOYSTICK_HIDAPI_XBOX"
const HINT_JOYSTICK_HIDAPI_XBOX_360 string = "SDL_JOYSTICK_HIDAPI_XBOX_360"
const HINT_JOYSTICK_HIDAPI_XBOX_360_PLAYER_LED string = "SDL_JOYSTICK_HIDAPI_XBOX_360_PLAYER_LED"
const HINT_JOYSTICK_HIDAPI_XBOX_360_WIRELESS string = "SDL_JOYSTICK_HIDAPI_XBOX_360_WIRELESS"
const HINT_JOYSTICK_HIDAPI_XBOX_ONE string = "SDL_JOYSTICK_HIDAPI_XBOX_ONE"
const HINT_JOYSTICK_HIDAPI_XBOX_ONE_HOME_LED string = "SDL_JOYSTICK_HIDAPI_XBOX_ONE_HOME_LED"
const HINT_JOYSTICK_IOKIT string = "SDL_JOYSTICK_IOKIT"
const HINT_JOYSTICK_LINUX_CLASSIC string = "SDL_JOYSTICK_LINUX_CLASSIC"
const HINT_JOYSTICK_LINUX_DEADZONES string = "SDL_JOYSTICK_LINUX_DEADZONES"
const HINT_JOYSTICK_LINUX_DIGITAL_HATS string = "SDL_JOYSTICK_LINUX_DIGITAL_HATS"
const HINT_JOYSTICK_LINUX_HAT_DEADZONES string = "SDL_JOYSTICK_LINUX_HAT_DEADZONES"
const HINT_JOYSTICK_MFI string = "SDL_JOYSTICK_MFI"
const HINT_JOYSTICK_RAWINPUT string = "SDL_JOYSTICK_RAWINPUT"
const HINT_JOYSTICK_RAWINPUT_CORRELATE_XINPUT string = "SDL_JOYSTICK_RAWINPUT_CORRELATE_XINPUT"
const HINT_JOYSTICK_ROG_CHAKRAM string = "SDL_JOYSTICK_ROG_CHAKRAM"
const HINT_JOYSTICK_THREAD string = "SDL_JOYSTICK_THREAD"
const HINT_JOYSTICK_THROTTLE_DEVICES string = "SDL_JOYSTICK_THROTTLE_DEVICES"
const HINT_JOYSTICK_THROTTLE_DEVICES_EXCLUDED string = "SDL_JOYSTICK_THROTTLE_DEVICES_EXCLUDED"
const HINT_JOYSTICK_WGI string = "SDL_JOYSTICK_WGI"
const HINT_JOYSTICK_WHEEL_DEVICES string = "SDL_JOYSTICK_WHEEL_DEVICES"
const HINT_JOYSTICK_WHEEL_DEVICES_EXCLUDED string = "SDL_JOYSTICK_WHEEL_DEVICES_EXCLUDED"
const HINT_JOYSTICK_ZERO_CENTERED_DEVICES string = "SDL_JOYSTICK_ZERO_CENTERED_DEVICES"
const HINT_KEYCODE_OPTIONS string = "SDL_KEYCODE_OPTIONS"
const HINT_KMSDRM_DEVICE_INDEX string = "SDL_KMSDRM_DEVICE_INDEX"
const HINT_KMSDRM_REQUIRE_DRM_MASTER string = "SDL_KMSDRM_REQUIRE_DRM_MASTER"
const HINT_LOGGING string = "SDL_LOGGING"
const HINT_MAC_BACKGROUND_APP string = "SDL_MAC_BACKGROUND_APP"
const HINT_MAC_CTRL_CLICK_EMULATE_RIGHT_CLICK string = "SDL_MAC_CTRL_CLICK_EMULATE_RIGHT_CLICK"
const HINT_MAC_OPENGL_ASYNC_DISPATCH string = "SDL_MAC_OPENGL_ASYNC_DISPATCH"
const HINT_MAC_OPTION_AS_ALT string = "SDL_MAC_OPTION_AS_ALT"
const HINT_MAC_SCROLL_MOMENTUM string = "SDL_MAC_SCROLL_MOMENTUM"
const HINT_MAIN_CALLBACK_RATE string = "SDL_MAIN_CALLBACK_RATE"
const HINT_MOUSE_AUTO_CAPTURE string = "SDL_MOUSE_AUTO_CAPTURE"
const HINT_MOUSE_DEFAULT_SYSTEM_CURSOR string = "SDL_MOUSE_DEFAULT_SYSTEM_CURSOR"
const HINT_MOUSE_DOUBLE_CLICK_RADIUS string = "SDL_MOUSE_DOUBLE_CLICK_RADIUS"
const HINT_MOUSE_DOUBLE_CLICK_TIME string = "SDL_MOUSE_DOUBLE_CLICK_TIME"
const HINT_MOUSE_EMULATE_WARP_WITH_RELATIVE string = "SDL_MOUSE_EMULATE_WARP_WITH_RELATIVE"
const HINT_MOUSE_FOCUS_CLICKTHROUGH string = "SDL_MOUSE_FOCUS_CLICKTHROUGH"
const HINT_MOUSE_NORMAL_SPEED_SCALE string = "SDL_MOUSE_NORMAL_SPEED_SCALE"
const HINT_MOUSE_RELATIVE_CURSOR_VISIBLE string = "SDL_MOUSE_RELATIVE_CURSOR_VISIBLE"
const HINT_MOUSE_RELATIVE_MODE_CENTER string = "SDL_MOUSE_RELATIVE_MODE_CENTER"
const HINT_MOUSE_RELATIVE_SPEED_SCALE string = "SDL_MOUSE_RELATIVE_SPEED_SCALE"
const HINT_MOUSE_RELATIVE_SYSTEM_SCALE string = "SDL_MOUSE_RELATIVE_SYSTEM_SCALE"
const HINT_MOUSE_RELATIVE_WARP_MOTION string = "SDL_MOUSE_RELATIVE_WARP_MOTION"
const HINT_MOUSE_TOUCH_EVENTS string = "SDL_MOUSE_TOUCH_EVENTS"
const HINT_MUTE_CONSOLE_KEYBOARD string = "SDL_MUTE_CONSOLE_KEYBOARD"
const HINT_NO_SIGNAL_HANDLERS string = "SDL_NO_SIGNAL_HANDLERS"
const HINT_OPENGL_ES_DRIVER string = "SDL_OPENGL_ES_DRIVER"
const HINT_OPENGL_LIBRARY string = "SDL_OPENGL_LIBRARY"
const HINT_OPENVR_LIBRARY string = "SDL_OPENVR_LIBRARY"
const HINT_ORIENTATIONS string = "SDL_ORIENTATIONS"
const HINT_PEN_MOUSE_EVENTS string = "SDL_PEN_MOUSE_EVENTS"
const HINT_PEN_TOUCH_EVENTS string = "SDL_PEN_TOUCH_EVENTS"
const HINT_POLL_SENTINEL string = "SDL_POLL_SENTINEL"
const HINT_PREFERRED_LOCALES string = "SDL_PREFERRED_LOCALES"
const HINT_QUIT_ON_LAST_WINDOW_CLOSE string = "SDL_QUIT_ON_LAST_WINDOW_CLOSE"
const HINT_RENDER_DIRECT3D11_DEBUG string = "SDL_RENDER_DIRECT3D11_DEBUG"
const HINT_RENDER_DIRECT3D_THREADSAFE string = "SDL_RENDER_DIRECT3D_THREADSAFE"
const HINT_RENDER_DRIVER string = "SDL_RENDER_DRIVER"
const HINT_RENDER_GPU_DEBUG string = "SDL_RENDER_GPU_DEBUG"
const HINT_RENDER_GPU_LOW_POWER string = "SDL_RENDER_GPU_LOW_POWER"
const HINT_RENDER_LINE_METHOD string = "SDL_RENDER_LINE_METHOD"
const HINT_RENDER_METAL_PREFER_LOW_POWER_DEVICE string = "SDL_RENDER_METAL_PREFER_LOW_POWER_DEVICE"
const HINT_RENDER_VSYNC string = "SDL_RENDER_VSYNC"
const HINT_RENDER_VULKAN_DEBUG string = "SDL_RENDER_VULKAN_DEBUG"
const HINT_RETURN_KEY_HIDES_IME string = "SDL_RETURN_KEY_HIDES_IME"
const HINT_ROG_GAMEPAD_MICE string = "SDL_ROG_GAMEPAD_MICE"
const HINT_ROG_GAMEPAD_MICE_EXCLUDED string = "SDL_ROG_GAMEPAD_MICE_EXCLUDED"
const HINT_RPI_VIDEO_LAYER string = "SDL_RPI_VIDEO_LAYER"
const HINT_SCREENSAVER_INHIBIT_ACTIVITY_NAME string = "SDL_SCREENSAVER_INHIBIT_ACTIVITY_NAME"
const HINT_SHUTDOWN_DBUS_ON_QUIT string = "SDL_SHUTDOWN_DBUS_ON_QUIT"
const HINT_STORAGE_TITLE_DRIVER string = "SDL_STORAGE_TITLE_DRIVER"
const HINT_STORAGE_USER_DRIVER string = "SDL_STORAGE_USER_DRIVER"
const HINT_THREAD_FORCE_REALTIME_TIME_CRITICAL string = "SDL_THREAD_FORCE_REALTIME_TIME_CRITICAL"
const HINT_THREAD_PRIORITY_POLICY string = "SDL_THREAD_PRIORITY_POLICY"
const HINT_TIMER_RESOLUTION string = "SDL_TIMER_RESOLUTION"
const HINT_TOUCH_MOUSE_EVENTS string = "SDL_TOUCH_MOUSE_EVENTS"
const HINT_TRACKPAD_IS_TOUCH_ONLY string = "SDL_TRACKPAD_IS_TOUCH_ONLY"
const HINT_TV_REMOTE_AS_JOYSTICK string = "SDL_TV_REMOTE_AS_JOYSTICK"
const HINT_VIDEO_ALLOW_SCREENSAVER string = "SDL_VIDEO_ALLOW_SCREENSAVER"
const HINT_VIDEO_DISPLAY_PRIORITY string = "SDL_VIDEO_DISPLAY_PRIORITY"
const HINT_VIDEO_DOUBLE_BUFFER string = "SDL_VIDEO_DOUBLE_BUFFER"
const HINT_VIDEO_DRIVER string = "SDL_VIDEO_DRIVER"
const HINT_VIDEO_DUMMY_SAVE_FRAMES string = "SDL_VIDEO_DUMMY_SAVE_FRAMES"
const HINT_VIDEO_EGL_ALLOW_GETDISPLAY_FALLBACK string = "SDL_VIDEO_EGL_ALLOW_GETDISPLAY_FALLBACK"
const HINT_VIDEO_FORCE_EGL string = "SDL_VIDEO_FORCE_EGL"
const HINT_VIDEO_MAC_FULLSCREEN_MENU_VISIBILITY string = "SDL_VIDEO_MAC_FULLSCREEN_MENU_VISIBILITY"
const HINT_VIDEO_MAC_FULLSCREEN_SPACES string = "SDL_VIDEO_MAC_FULLSCREEN_SPACES"
const HINT_VIDEO_MINIMIZE_ON_FOCUS_LOSS string = "SDL_VIDEO_MINIMIZE_ON_FOCUS_LOSS"
const HINT_VIDEO_OFFSCREEN_SAVE_FRAMES string = "SDL_VIDEO_OFFSCREEN_SAVE_FRAMES"
const HINT_VIDEO_SYNC_WINDOW_OPERATIONS string = "SDL_VIDEO_SYNC_WINDOW_OPERATIONS"
const HINT_VIDEO_WAYLAND_ALLOW_LIBDECOR string = "SDL_VIDEO_WAYLAND_ALLOW_LIBDECOR"
const HINT_VIDEO_WAYLAND_MODE_EMULATION string = "SDL_VIDEO_WAYLAND_MODE_EMULATION"
const HINT_VIDEO_WAYLAND_MODE_SCALING string = "SDL_VIDEO_WAYLAND_MODE_SCALING"
const HINT_VIDEO_WAYLAND_PREFER_LIBDECOR string = "SDL_VIDEO_WAYLAND_PREFER_LIBDECOR"
const HINT_VIDEO_WAYLAND_SCALE_TO_DISPLAY string = "SDL_VIDEO_WAYLAND_SCALE_TO_DISPLAY"
const HINT_VIDEO_WIN_D3DCOMPILER string = "SDL_VIDEO_WIN_D3DCOMPILER"
const HINT_VIDEO_X11_EXTERNAL_WINDOW_INPUT string = "SDL_VIDEO_X11_EXTERNAL_WINDOW_INPUT"
const HINT_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR string = "SDL_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR"
const HINT_VIDEO_X11_NET_WM_PING string = "SDL_VIDEO_X11_NET_WM_PING"
const HINT_VIDEO_X11_NODIRECTCOLOR string = "SDL_VIDEO_X11_NODIRECTCOLOR"
const HINT_VIDEO_X11_SCALING_FACTOR string = "SDL_VIDEO_X11_SCALING_FACTOR"
const HINT_VIDEO_X11_VISUALID string = "SDL_VIDEO_X11_VISUALID"
const HINT_VIDEO_X11_WINDOW_VISUALID string = "SDL_VIDEO_X11_WINDOW_VISUALID"
const HINT_VIDEO_X11_XRANDR string = "SDL_VIDEO_X11_XRANDR"
const HINT_VITA_ENABLE_BACK_TOUCH string = "SDL_VITA_ENABLE_BACK_TOUCH"
const HINT_VITA_ENABLE_FRONT_TOUCH string = "SDL_VITA_ENABLE_FRONT_TOUCH"
const HINT_VITA_MODULE_PATH string = "SDL_VITA_MODULE_PATH"
const HINT_VITA_PVR_INIT string = "SDL_VITA_PVR_INIT"
const HINT_VITA_PVR_OPENGL string = "SDL_VITA_PVR_OPENGL"
const HINT_VITA_RESOLUTION string = "SDL_VITA_RESOLUTION"
const HINT_VITA_TOUCH_MOUSE_DEVICE string = "SDL_VITA_TOUCH_MOUSE_DEVICE"
const HINT_VULKAN_DISPLAY string = "SDL_VULKAN_DISPLAY"
const HINT_VULKAN_LIBRARY string = "SDL_VULKAN_LIBRARY"
const HINT_WAVE_CHUNK_LIMIT string = "SDL_WAVE_CHUNK_LIMIT"
const HINT_WAVE_FACT_CHUNK string = "SDL_WAVE_FACT_CHUNK"
const HINT_WAVE_RIFF_CHUNK_SIZE string = "SDL_WAVE_RIFF_CHUNK_SIZE"
const HINT_WAVE_TRUNCATION string = "SDL_WAVE_TRUNCATION"
const HINT_WINDOWS_CLOSE_ON_ALT_F4 string = "SDL_WINDOWS_CLOSE_ON_ALT_F4"
const HINT_WINDOWS_ENABLE_MENU_MNEMONICS string = "SDL_WINDOWS_ENABLE_MENU_MNEMONICS"
const HINT_WINDOWS_ENABLE_MESSAGELOOP string = "SDL_WINDOWS_ENABLE_MESSAGELOOP"
const HINT_WINDOWS_ERASE_BACKGROUND_MODE string = "SDL_WINDOWS_ERASE_BACKGROUND_MODE"
const HINT_WINDOWS_FORCE_SEMAPHORE_KERNEL string = "SDL_WINDOWS_FORCE_SEMAPHORE_KERNEL"
const HINT_WINDOWS_GAMEINPUT string = "SDL_WINDOWS_GAMEINPUT"
const HINT_WINDOWS_INTRESOURCE_ICON string = "SDL_WINDOWS_INTRESOURCE_ICON"
const HINT_WINDOWS_INTRESOURCE_ICON_SMALL string = "SDL_WINDOWS_INTRESOURCE_ICON_SMALL"
const HINT_WINDOWS_RAW_KEYBOARD string = "SDL_WINDOWS_RAW_KEYBOARD"
const HINT_WINDOWS_USE_D3D9EX string = "SDL_WINDOWS_USE_D3D9EX"
const HINT_WINDOW_ACTIVATE_WHEN_RAISED string = "SDL_WINDOW_ACTIVATE_WHEN_RAISED"
const HINT_WINDOW_ACTIVATE_WHEN_SHOWN string = "SDL_WINDOW_ACTIVATE_WHEN_SHOWN"
const HINT_WINDOW_ALLOW_TOPMOST string = "SDL_WINDOW_ALLOW_TOPMOST"
const HINT_WINDOW_FRAME_USABLE_WHILE_CURSOR_HIDDEN string = "SDL_WINDOW_FRAME_USABLE_WHILE_CURSOR_HIDDEN"
const HINT_X11_FORCE_OVERRIDE_REDIRECT string = "SDL_X11_FORCE_OVERRIDE_REDIRECT"
const HINT_X11_WINDOW_TYPE string = "SDL_X11_WINDOW_TYPE"
const HINT_X11_XCB_LIBRARY string = "SDL_X11_XCB_LIBRARY"
const HINT_XINPUT_ENABLED string = "SDL_XINPUT_ENABLED"
const PROP_APP_METADATA_COPYRIGHT_STRING string = "SDL.app.metadata.copyright"
const PROP_APP_METADATA_CREATOR_STRING string = "SDL.app.metadata.creator"
const PROP_APP_METADATA_IDENTIFIER_STRING string = "SDL.app.metadata.identifier"
const PROP_APP_METADATA_NAME_STRING string = "SDL.app.metadata.name"
const PROP_APP_METADATA_TYPE_STRING string = "SDL.app.metadata.type"
const PROP_APP_METADATA_URL_STRING string = "SDL.app.metadata.url"
const PROP_APP_METADATA_VERSION_STRING string = "SDL.app.metadata.version"
const PROP_DISPLAY_HDR_ENABLED_BOOLEAN string = "SDL.display.HDR_enabled"
const PROP_DISPLAY_KMSDRM_PANEL_ORIENTATION_NUMBER string = "SDL.display.KMSDRM.panel_orientation"
const PROP_FILE_DIALOG_ACCEPT_STRING string = "SDL.filedialog.accept"
const PROP_FILE_DIALOG_CANCEL_STRING string = "SDL.filedialog.cancel"
const PROP_FILE_DIALOG_FILTERS_POINTER string = "SDL.filedialog.filters"
const PROP_FILE_DIALOG_LOCATION_STRING string = "SDL.filedialog.location"
const PROP_FILE_DIALOG_MANY_BOOLEAN string = "SDL.filedialog.many"
const PROP_FILE_DIALOG_NFILTERS_NUMBER string = "SDL.filedialog.nfilters"
const PROP_FILE_DIALOG_TITLE_STRING string = "SDL.filedialog.title"
const PROP_FILE_DIALOG_WINDOW_POINTER string = "SDL.filedialog.window"
const PROP_GAMEPAD_CAP_MONO_LED_BOOLEAN string = "SDL.joystick.cap.mono_led"
const PROP_GAMEPAD_CAP_PLAYER_LED_BOOLEAN string = "SDL.joystick.cap.player_led"
const PROP_GAMEPAD_CAP_RGB_LED_BOOLEAN string = "SDL.joystick.cap.rgb_led"
const PROP_GAMEPAD_CAP_RUMBLE_BOOLEAN string = "SDL.joystick.cap.rumble"
const PROP_GAMEPAD_CAP_TRIGGER_RUMBLE_BOOLEAN string = "SDL.joystick.cap.trigger_rumble"
const PROP_GLOBAL_VIDEO_WAYLAND_WL_DISPLAY_POINTER string = "SDL.video.wayland.wl_display"
const PROP_GPU_BUFFER_CREATE_NAME_STRING string = "SDL.gpu.buffer.create.name"
const PROP_GPU_COMPUTEPIPELINE_CREATE_NAME_STRING string = "SDL.gpu.computepipeline.create.name"
const PROP_GPU_DEVICE_CREATE_D3D12_SEMANTIC_NAME_STRING string = "SDL.gpu.device.create.d3d12.semantic"
const PROP_GPU_DEVICE_CREATE_DEBUGMODE_BOOLEAN string = "SDL.gpu.device.create.debugmode"
const PROP_GPU_DEVICE_CREATE_NAME_STRING string = "SDL.gpu.device.create.name"
const PROP_GPU_DEVICE_CREATE_PREFERLOWPOWER_BOOLEAN string = "SDL.gpu.device.create.preferlowpower"
const PROP_GPU_DEVICE_CREATE_SHADERS_DXBC_BOOLEAN string = "SDL.gpu.device.create.shaders.dxbc"
const PROP_GPU_DEVICE_CREATE_SHADERS_DXIL_BOOLEAN string = "SDL.gpu.device.create.shaders.dxil"
const PROP_GPU_DEVICE_CREATE_SHADERS_METALLIB_BOOLEAN string = "SDL.gpu.device.create.shaders.metallib"
const PROP_GPU_DEVICE_CREATE_SHADERS_MSL_BOOLEAN string = "SDL.gpu.device.create.shaders.msl"
const PROP_GPU_DEVICE_CREATE_SHADERS_PRIVATE_BOOLEAN string = "SDL.gpu.device.create.shaders.private"
const PROP_GPU_DEVICE_CREATE_SHADERS_SPIRV_BOOLEAN string = "SDL.gpu.device.create.shaders.spirv"
const PROP_GPU_GRAPHICSPIPELINE_CREATE_NAME_STRING string = "SDL.gpu.graphicspipeline.create.name"
const PROP_GPU_SAMPLER_CREATE_NAME_STRING string = "SDL.gpu.sampler.create.name"
const PROP_GPU_SHADER_CREATE_NAME_STRING string = "SDL.gpu.shader.create.name"
const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_A_FLOAT string = "SDL.gpu.texture.create.d3d12.clear.a"
const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_B_FLOAT string = "SDL.gpu.texture.create.d3d12.clear.b"
const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_DEPTH_FLOAT string = "SDL.gpu.texture.create.d3d12.clear.depth"
const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_G_FLOAT string = "SDL.gpu.texture.create.d3d12.clear.g"
const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_R_FLOAT string = "SDL.gpu.texture.create.d3d12.clear.r"
const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_STENCIL_NUMBER string = "SDL.gpu.texture.create.d3d12.clear.stencil"
const PROP_GPU_TEXTURE_CREATE_NAME_STRING string = "SDL.gpu.texture.create.name"
const PROP_GPU_TRANSFERBUFFER_CREATE_NAME_STRING string = "SDL.gpu.transferbuffer.create.name"
const PROP_IOSTREAM_ANDROID_AASSET_POINTER string = "SDL.iostream.android.aasset"
const PROP_IOSTREAM_DYNAMIC_CHUNKSIZE_NUMBER string = "SDL.iostream.dynamic.chunksize"
const PROP_IOSTREAM_DYNAMIC_MEMORY_POINTER string = "SDL.iostream.dynamic.memory"
const PROP_IOSTREAM_FILE_DESCRIPTOR_NUMBER string = "SDL.iostream.file_descriptor"
const PROP_IOSTREAM_MEMORY_POINTER string = "SDL.iostream.memory.base"
const PROP_IOSTREAM_MEMORY_SIZE_NUMBER string = "SDL.iostream.memory.size"
const PROP_IOSTREAM_STDIO_FILE_POINTER string = "SDL.iostream.stdio.file"
const PROP_IOSTREAM_WINDOWS_HANDLE_POINTER string = "SDL.iostream.windows.handle"
const PROP_JOYSTICK_CAP_MONO_LED_BOOLEAN string = "SDL.joystick.cap.mono_led"
const PROP_JOYSTICK_CAP_PLAYER_LED_BOOLEAN string = "SDL.joystick.cap.player_led"
const PROP_JOYSTICK_CAP_RGB_LED_BOOLEAN string = "SDL.joystick.cap.rgb_led"
const PROP_JOYSTICK_CAP_RUMBLE_BOOLEAN string = "SDL.joystick.cap.rumble"
const PROP_JOYSTICK_CAP_TRIGGER_RUMBLE_BOOLEAN string = "SDL.joystick.cap.trigger_rumble"
const PROP_PROCESS_BACKGROUND_BOOLEAN string = "SDL.process.background"
const PROP_PROCESS_CREATE_ARGS_POINTER string = "SDL.process.create.args"
const PROP_PROCESS_CREATE_BACKGROUND_BOOLEAN string = "SDL.process.create.background"
const PROP_PROCESS_CREATE_ENVIRONMENT_POINTER string = "SDL.process.create.environment"
const PROP_PROCESS_CREATE_STDERR_NUMBER string = "SDL.process.create.stderr_option"
const PROP_PROCESS_CREATE_STDERR_POINTER string = "SDL.process.create.stderr_source"
const PROP_PROCESS_CREATE_STDERR_TO_STDOUT_BOOLEAN string = "SDL.process.create.stderr_to_stdout"
const PROP_PROCESS_CREATE_STDIN_NUMBER string = "SDL.process.create.stdin_option"
const PROP_PROCESS_CREATE_STDIN_POINTER string = "SDL.process.create.stdin_source"
const PROP_PROCESS_CREATE_STDOUT_NUMBER string = "SDL.process.create.stdout_option"
const PROP_PROCESS_CREATE_STDOUT_POINTER string = "SDL.process.create.stdout_source"
const PROP_PROCESS_PID_NUMBER string = "SDL.process.pid"
const PROP_PROCESS_STDERR_POINTER string = "SDL.process.stderr"
const PROP_PROCESS_STDIN_POINTER string = "SDL.process.stdin"
const PROP_PROCESS_STDOUT_POINTER string = "SDL.process.stdout"
const PROP_RENDERER_CREATE_NAME_STRING string = "SDL.renderer.create.name"
const PROP_RENDERER_CREATE_OUTPUT_COLORSPACE_NUMBER string = "SDL.renderer.create.output_colorspace"
const PROP_RENDERER_CREATE_PRESENT_VSYNC_NUMBER string = "SDL.renderer.create.present_vsync"
const PROP_RENDERER_CREATE_SURFACE_POINTER string = "SDL.renderer.create.surface"
const PROP_RENDERER_CREATE_VULKAN_DEVICE_POINTER string = "SDL.renderer.create.vulkan.device"
const PROP_RENDERER_CREATE_VULKAN_GRAPHICS_QUEUE_FAMILY_INDEX_NUMBER string = "SDL.renderer.create.vulkan.graphics_queue_family_index"
const PROP_RENDERER_CREATE_VULKAN_INSTANCE_POINTER string = "SDL.renderer.create.vulkan.instance"
const PROP_RENDERER_CREATE_VULKAN_PHYSICAL_DEVICE_POINTER string = "SDL.renderer.create.vulkan.physical_device"
const PROP_RENDERER_CREATE_VULKAN_PRESENT_QUEUE_FAMILY_INDEX_NUMBER string = "SDL.renderer.create.vulkan.present_queue_family_index"
const PROP_RENDERER_CREATE_VULKAN_SURFACE_NUMBER string = "SDL.renderer.create.vulkan.surface"
const PROP_RENDERER_CREATE_WINDOW_POINTER string = "SDL.renderer.create.window"
const PROP_RENDERER_D3D11_DEVICE_POINTER string = "SDL.renderer.d3d11.device"
const PROP_RENDERER_D3D11_SWAPCHAIN_POINTER string = "SDL.renderer.d3d11.swap_chain"
const PROP_RENDERER_D3D12_COMMAND_QUEUE_POINTER string = "SDL.renderer.d3d12.command_queue"
const PROP_RENDERER_D3D12_DEVICE_POINTER string = "SDL.renderer.d3d12.device"
const PROP_RENDERER_D3D12_SWAPCHAIN_POINTER string = "SDL.renderer.d3d12.swap_chain"
const PROP_RENDERER_D3D9_DEVICE_POINTER string = "SDL.renderer.d3d9.device"
const PROP_RENDERER_GPU_DEVICE_POINTER string = "SDL.renderer.gpu.device"
const PROP_RENDERER_HDR_ENABLED_BOOLEAN string = "SDL.renderer.HDR_enabled"
const PROP_RENDERER_HDR_HEADROOM_FLOAT string = "SDL.renderer.HDR_headroom"
const PROP_RENDERER_MAX_TEXTURE_SIZE_NUMBER string = "SDL.renderer.max_texture_size"
const PROP_RENDERER_NAME_STRING string = "SDL.renderer.name"
const PROP_RENDERER_OUTPUT_COLORSPACE_NUMBER string = "SDL.renderer.output_colorspace"
const PROP_RENDERER_SDR_WHITE_POINT_FLOAT string = "SDL.renderer.SDR_white_point"
const PROP_RENDERER_SURFACE_POINTER string = "SDL.renderer.surface"
const PROP_RENDERER_TEXTURE_FORMATS_POINTER string = "SDL.renderer.texture_formats"
const PROP_RENDERER_VSYNC_NUMBER string = "SDL.renderer.vsync"
const PROP_RENDERER_VULKAN_DEVICE_POINTER string = "SDL.renderer.vulkan.device"
const PROP_RENDERER_VULKAN_GRAPHICS_QUEUE_FAMILY_INDEX_NUMBER string = "SDL.renderer.vulkan.graphics_queue_family_index"
const PROP_RENDERER_VULKAN_INSTANCE_POINTER string = "SDL.renderer.vulkan.instance"
const PROP_RENDERER_VULKAN_PHYSICAL_DEVICE_POINTER string = "SDL.renderer.vulkan.physical_device"
const PROP_RENDERER_VULKAN_PRESENT_QUEUE_FAMILY_INDEX_NUMBER string = "SDL.renderer.vulkan.present_queue_family_index"
const PROP_RENDERER_VULKAN_SURFACE_NUMBER string = "SDL.renderer.vulkan.surface"
const PROP_RENDERER_VULKAN_SWAPCHAIN_IMAGE_COUNT_NUMBER string = "SDL.renderer.vulkan.swapchain_image_count"
const PROP_RENDERER_WINDOW_POINTER string = "SDL.renderer.window"
const PROP_SURFACE_HDR_HEADROOM_FLOAT string = "SDL.surface.HDR_headroom"
const PROP_SURFACE_HOTSPOT_X_NUMBER string = "SDL.surface.hotspot.x"
const PROP_SURFACE_HOTSPOT_Y_NUMBER string = "SDL.surface.hotspot.y"
const PROP_SURFACE_SDR_WHITE_POINT_FLOAT string = "SDL.surface.SDR_white_point"
const PROP_SURFACE_TONEMAP_OPERATOR_STRING string = "SDL.surface.tonemap"
const PROP_TEXTINPUT_ANDROID_INPUTTYPE_NUMBER string = "SDL.textinput.android.inputtype"
const PROP_TEXTINPUT_AUTOCORRECT_BOOLEAN string = "SDL.textinput.autocorrect"
const PROP_TEXTINPUT_CAPITALIZATION_NUMBER string = "SDL.textinput.capitalization"
const PROP_TEXTINPUT_MULTILINE_BOOLEAN string = "SDL.textinput.multiline"
const PROP_TEXTINPUT_TYPE_NUMBER string = "SDL.textinput.type"
const PROP_TEXTURE_ACCESS_NUMBER string = "SDL.texture.access"
const PROP_TEXTURE_COLORSPACE_NUMBER string = "SDL.texture.colorspace"
const PROP_TEXTURE_CREATE_ACCESS_NUMBER string = "SDL.texture.create.access"
const PROP_TEXTURE_CREATE_COLORSPACE_NUMBER string = "SDL.texture.create.colorspace"
const PROP_TEXTURE_CREATE_D3D11_TEXTURE_POINTER string = "SDL.texture.create.d3d11.texture"
const PROP_TEXTURE_CREATE_D3D11_TEXTURE_U_POINTER string = "SDL.texture.create.d3d11.texture_u"
const PROP_TEXTURE_CREATE_D3D11_TEXTURE_V_POINTER string = "SDL.texture.create.d3d11.texture_v"
const PROP_TEXTURE_CREATE_D3D12_TEXTURE_POINTER string = "SDL.texture.create.d3d12.texture"
const PROP_TEXTURE_CREATE_D3D12_TEXTURE_U_POINTER string = "SDL.texture.create.d3d12.texture_u"
const PROP_TEXTURE_CREATE_D3D12_TEXTURE_V_POINTER string = "SDL.texture.create.d3d12.texture_v"
const PROP_TEXTURE_CREATE_FORMAT_NUMBER string = "SDL.texture.create.format"
const PROP_TEXTURE_CREATE_HDR_HEADROOM_FLOAT string = "SDL.texture.create.HDR_headroom"
const PROP_TEXTURE_CREATE_HEIGHT_NUMBER string = "SDL.texture.create.height"
const PROP_TEXTURE_CREATE_METAL_PIXELBUFFER_POINTER string = "SDL.texture.create.metal.pixelbuffer"
const PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_NUMBER string = "SDL.texture.create.opengles2.texture"
const PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_UV_NUMBER string = "SDL.texture.create.opengles2.texture_uv"
const PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_U_NUMBER string = "SDL.texture.create.opengles2.texture_u"
const PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_V_NUMBER string = "SDL.texture.create.opengles2.texture_v"
const PROP_TEXTURE_CREATE_OPENGL_TEXTURE_NUMBER string = "SDL.texture.create.opengl.texture"
const PROP_TEXTURE_CREATE_OPENGL_TEXTURE_UV_NUMBER string = "SDL.texture.create.opengl.texture_uv"
const PROP_TEXTURE_CREATE_OPENGL_TEXTURE_U_NUMBER string = "SDL.texture.create.opengl.texture_u"
const PROP_TEXTURE_CREATE_OPENGL_TEXTURE_V_NUMBER string = "SDL.texture.create.opengl.texture_v"
const PROP_TEXTURE_CREATE_SDR_WHITE_POINT_FLOAT string = "SDL.texture.create.SDR_white_point"
const PROP_TEXTURE_CREATE_VULKAN_TEXTURE_NUMBER string = "SDL.texture.create.vulkan.texture"
const PROP_TEXTURE_CREATE_WIDTH_NUMBER string = "SDL.texture.create.width"
const PROP_TEXTURE_D3D11_TEXTURE_POINTER string = "SDL.texture.d3d11.texture"
const PROP_TEXTURE_D3D11_TEXTURE_U_POINTER string = "SDL.texture.d3d11.texture_u"
const PROP_TEXTURE_D3D11_TEXTURE_V_POINTER string = "SDL.texture.d3d11.texture_v"
const PROP_TEXTURE_D3D12_TEXTURE_POINTER string = "SDL.texture.d3d12.texture"
const PROP_TEXTURE_D3D12_TEXTURE_U_POINTER string = "SDL.texture.d3d12.texture_u"
const PROP_TEXTURE_D3D12_TEXTURE_V_POINTER string = "SDL.texture.d3d12.texture_v"
const PROP_TEXTURE_FORMAT_NUMBER string = "SDL.texture.format"
const PROP_TEXTURE_HDR_HEADROOM_FLOAT string = "SDL.texture.HDR_headroom"
const PROP_TEXTURE_HEIGHT_NUMBER string = "SDL.texture.height"
const PROP_TEXTURE_OPENGLES2_TEXTURE_NUMBER string = "SDL.texture.opengles2.texture"
const PROP_TEXTURE_OPENGLES2_TEXTURE_TARGET_NUMBER string = "SDL.texture.opengles2.target"
const PROP_TEXTURE_OPENGLES2_TEXTURE_UV_NUMBER string = "SDL.texture.opengles2.texture_uv"
const PROP_TEXTURE_OPENGLES2_TEXTURE_U_NUMBER string = "SDL.texture.opengles2.texture_u"
const PROP_TEXTURE_OPENGLES2_TEXTURE_V_NUMBER string = "SDL.texture.opengles2.texture_v"
const PROP_TEXTURE_OPENGL_TEXTURE_NUMBER string = "SDL.texture.opengl.texture"
const PROP_TEXTURE_OPENGL_TEXTURE_TARGET_NUMBER string = "SDL.texture.opengl.target"
const PROP_TEXTURE_OPENGL_TEXTURE_UV_NUMBER string = "SDL.texture.opengl.texture_uv"
const PROP_TEXTURE_OPENGL_TEXTURE_U_NUMBER string = "SDL.texture.opengl.texture_u"
const PROP_TEXTURE_OPENGL_TEXTURE_V_NUMBER string = "SDL.texture.opengl.texture_v"
const PROP_TEXTURE_OPENGL_TEX_H_FLOAT string = "SDL.texture.opengl.tex_h"
const PROP_TEXTURE_OPENGL_TEX_W_FLOAT string = "SDL.texture.opengl.tex_w"
const PROP_TEXTURE_SDR_WHITE_POINT_FLOAT string = "SDL.texture.SDR_white_point"
const PROP_TEXTURE_VULKAN_TEXTURE_NUMBER string = "SDL.texture.vulkan.texture"
const PROP_TEXTURE_WIDTH_NUMBER string = "SDL.texture.width"
const PROP_THREAD_CREATE_ENTRY_FUNCTION_POINTER string = "SDL.thread.create.entry_function"
const PROP_THREAD_CREATE_NAME_STRING string = "SDL.thread.create.name"
const PROP_THREAD_CREATE_STACKSIZE_NUMBER string = "SDL.thread.create.stacksize"
const PROP_THREAD_CREATE_USERDATA_POINTER string = "SDL.thread.create.userdata"
const PROP_WINDOW_ANDROID_SURFACE_POINTER string = "SDL.window.android.surface"
const PROP_WINDOW_ANDROID_WINDOW_POINTER string = "SDL.window.android.window"
const PROP_WINDOW_COCOA_METAL_VIEW_TAG_NUMBER string = "SDL.window.cocoa.metal_view_tag"
const PROP_WINDOW_COCOA_WINDOW_POINTER string = "SDL.window.cocoa.window"
const PROP_WINDOW_CREATE_ALWAYS_ON_TOP_BOOLEAN string = "SDL.window.create.always_on_top"
const PROP_WINDOW_CREATE_BORDERLESS_BOOLEAN string = "SDL.window.create.borderless"
const PROP_WINDOW_CREATE_COCOA_VIEW_POINTER string = "SDL.window.create.cocoa.view"
const PROP_WINDOW_CREATE_COCOA_WINDOW_POINTER string = "SDL.window.create.cocoa.window"
const PROP_WINDOW_CREATE_CONSTRAIN_POPUP_BOOLEAN string = "SDL.window.create.constrain_popup"
const PROP_WINDOW_CREATE_EXTERNAL_GRAPHICS_CONTEXT_BOOLEAN string = "SDL.window.create.external_graphics_context"
const PROP_WINDOW_CREATE_FLAGS_NUMBER string = "SDL.window.create.flags"
const PROP_WINDOW_CREATE_FOCUSABLE_BOOLEAN string = "SDL.window.create.focusable"
const PROP_WINDOW_CREATE_FULLSCREEN_BOOLEAN string = "SDL.window.create.fullscreen"
const PROP_WINDOW_CREATE_HEIGHT_NUMBER string = "SDL.window.create.height"
const PROP_WINDOW_CREATE_HIDDEN_BOOLEAN string = "SDL.window.create.hidden"
const PROP_WINDOW_CREATE_HIGH_PIXEL_DENSITY_BOOLEAN string = "SDL.window.create.high_pixel_density"
const PROP_WINDOW_CREATE_MAXIMIZED_BOOLEAN string = "SDL.window.create.maximized"
const PROP_WINDOW_CREATE_MENU_BOOLEAN string = "SDL.window.create.menu"
const PROP_WINDOW_CREATE_METAL_BOOLEAN string = "SDL.window.create.metal"
const PROP_WINDOW_CREATE_MINIMIZED_BOOLEAN string = "SDL.window.create.minimized"
const PROP_WINDOW_CREATE_MODAL_BOOLEAN string = "SDL.window.create.modal"
const PROP_WINDOW_CREATE_MOUSE_GRABBED_BOOLEAN string = "SDL.window.create.mouse_grabbed"
const PROP_WINDOW_CREATE_OPENGL_BOOLEAN string = "SDL.window.create.opengl"
const PROP_WINDOW_CREATE_PARENT_POINTER string = "SDL.window.create.parent"
const PROP_WINDOW_CREATE_RESIZABLE_BOOLEAN string = "SDL.window.create.resizable"
const PROP_WINDOW_CREATE_TITLE_STRING string = "SDL.window.create.title"
const PROP_WINDOW_CREATE_TOOLTIP_BOOLEAN string = "SDL.window.create.tooltip"
const PROP_WINDOW_CREATE_TRANSPARENT_BOOLEAN string = "SDL.window.create.transparent"
const PROP_WINDOW_CREATE_UTILITY_BOOLEAN string = "SDL.window.create.utility"
const PROP_WINDOW_CREATE_VULKAN_BOOLEAN string = "SDL.window.create.vulkan"
const PROP_WINDOW_CREATE_WAYLAND_CREATE_EGL_WINDOW_BOOLEAN string = "SDL.window.create.wayland.create_egl_window"
const PROP_WINDOW_CREATE_WAYLAND_SURFACE_ROLE_CUSTOM_BOOLEAN string = "SDL.window.create.wayland.surface_role_custom"
const PROP_WINDOW_CREATE_WAYLAND_WL_SURFACE_POINTER string = "SDL.window.create.wayland.wl_surface"
const PROP_WINDOW_CREATE_WIDTH_NUMBER string = "SDL.window.create.width"
const PROP_WINDOW_CREATE_WIN32_HWND_POINTER string = "SDL.window.create.win32.hwnd"
const PROP_WINDOW_CREATE_WIN32_PIXEL_FORMAT_HWND_POINTER string = "SDL.window.create.win32.pixel_format_hwnd"
const PROP_WINDOW_CREATE_X11_WINDOW_NUMBER string = "SDL.window.create.x11.window"
const PROP_WINDOW_CREATE_X_NUMBER string = "SDL.window.create.x"
const PROP_WINDOW_CREATE_Y_NUMBER string = "SDL.window.create.y"
const PROP_WINDOW_HDR_ENABLED_BOOLEAN string = "SDL.window.HDR_enabled"
const PROP_WINDOW_HDR_HEADROOM_FLOAT string = "SDL.window.HDR_headroom"
const PROP_WINDOW_KMSDRM_DEVICE_INDEX_NUMBER string = "SDL.window.kmsdrm.dev_index"
const PROP_WINDOW_KMSDRM_DRM_FD_NUMBER string = "SDL.window.kmsdrm.drm_fd"
const PROP_WINDOW_KMSDRM_GBM_DEVICE_POINTER string = "SDL.window.kmsdrm.gbm_dev"
const PROP_WINDOW_OPENVR_OVERLAY_ID string = "SDL.window.openvr.overlay_id"
const PROP_WINDOW_SDR_WHITE_LEVEL_FLOAT string = "SDL.window.SDR_white_level"
const PROP_WINDOW_SHAPE_POINTER string = "SDL.window.shape"
const PROP_WINDOW_UIKIT_METAL_VIEW_TAG_NUMBER string = "SDL.window.uikit.metal_view_tag"
const PROP_WINDOW_UIKIT_OPENGL_FRAMEBUFFER_NUMBER string = "SDL.window.uikit.opengl.framebuffer"
const PROP_WINDOW_UIKIT_OPENGL_RENDERBUFFER_NUMBER string = "SDL.window.uikit.opengl.renderbuffer"
const PROP_WINDOW_UIKIT_OPENGL_RESOLVE_FRAMEBUFFER_NUMBER string = "SDL.window.uikit.opengl.resolve_framebuffer"
const PROP_WINDOW_UIKIT_WINDOW_POINTER string = "SDL.window.uikit.window"
const PROP_WINDOW_VIVANTE_DISPLAY_POINTER string = "SDL.window.vivante.display"
const PROP_WINDOW_VIVANTE_SURFACE_POINTER string = "SDL.window.vivante.surface"
const PROP_WINDOW_VIVANTE_WINDOW_POINTER string = "SDL.window.vivante.window"
const PROP_WINDOW_WAYLAND_DISPLAY_POINTER string = "SDL.window.wayland.display"
const PROP_WINDOW_WAYLAND_EGL_WINDOW_POINTER string = "SDL.window.wayland.egl_window"
const PROP_WINDOW_WAYLAND_SURFACE_POINTER string = "SDL.window.wayland.surface"
const PROP_WINDOW_WAYLAND_VIEWPORT_POINTER string = "SDL.window.wayland.viewport"
const PROP_WINDOW_WAYLAND_XDG_POPUP_POINTER string = "SDL.window.wayland.xdg_popup"
const PROP_WINDOW_WAYLAND_XDG_POSITIONER_POINTER string = "SDL.window.wayland.xdg_positioner"
const PROP_WINDOW_WAYLAND_XDG_SURFACE_POINTER string = "SDL.window.wayland.xdg_surface"
const PROP_WINDOW_WAYLAND_XDG_TOPLEVEL_EXPORT_HANDLE_STRING string = "SDL.window.wayland.xdg_toplevel_export_handle"
const PROP_WINDOW_WAYLAND_XDG_TOPLEVEL_POINTER string = "SDL.window.wayland.xdg_toplevel"
const PROP_WINDOW_WIN32_HDC_POINTER string = "SDL.window.win32.hdc"
const PROP_WINDOW_WIN32_HWND_POINTER string = "SDL.window.win32.hwnd"
const PROP_WINDOW_WIN32_INSTANCE_POINTER string = "SDL.window.win32.instance"
const PROP_WINDOW_X11_DISPLAY_POINTER string = "SDL.window.x11.display"
const PROP_WINDOW_X11_SCREEN_NUMBER string = "SDL.window.x11.screen"
const PROP_WINDOW_X11_WINDOW_NUMBER string = "SDL.window.x11.window"

fn malloc(size u64) rawptr {
  return raw::malloc(size)
}

fn calloc(nmemb u64, size u64) rawptr {
  return raw::calloc(nmemb, size)
}

fn realloc(mem rawptr, size u64) rawptr {
  return raw::realloc(mem, size)
}

fn free(mem rawptr) {
  raw::free(mem)
}

fn aligned_alloc(alignment u64, size u64) rawptr {
  return raw::aligned_alloc(alignment, size)
}

fn aligned_free(mem rawptr) {
  raw::aligned_free(mem)
}

fn get_num_allocations() i32 {
  return raw::get_num_allocations()
}

fn get_environment() *Environment {
  return raw::get_environment()
}

fn create_environment(populated bool) *Environment {
  return raw::create_environment(populated)
}

fn destroy_environment(env *Environment) {
  raw::destroy_environment(env)
}

fn abs(x i32) i32 {
  return raw::abs(x)
}

fn isalpha(x i32) i32 {
  return raw::isalpha(x)
}

fn isalnum(x i32) i32 {
  return raw::isalnum(x)
}

fn isblank(x i32) i32 {
  return raw::isblank(x)
}

fn iscntrl(x i32) i32 {
  return raw::iscntrl(x)
}

fn isdigit(x i32) i32 {
  return raw::isdigit(x)
}

fn isxdigit(x i32) i32 {
  return raw::isxdigit(x)
}

fn ispunct(x i32) i32 {
  return raw::ispunct(x)
}

fn isspace(x i32) i32 {
  return raw::isspace(x)
}

fn isupper(x i32) i32 {
  return raw::isupper(x)
}

fn islower(x i32) i32 {
  return raw::islower(x)
}

fn isprint(x i32) i32 {
  return raw::isprint(x)
}

fn isgraph(x i32) i32 {
  return raw::isgraph(x)
}

fn toupper(x i32) i32 {
  return raw::toupper(x)
}

fn tolower(x i32) i32 {
  return raw::tolower(x)
}

fn crc16(crc u16, data const *Void, len u64) u16 {
  return raw::crc16(crc, data, len)
}

fn crc32(crc u32, data const *Void, len u64) u32 {
  return raw::crc32(crc, data, len)
}

fn murmur3_32(data const *Void, len u64, seed u32) u32 {
  return raw::murmur3_32(data, len, seed)
}

fn memcpy(dst rawptr, src const *Void, len u64) rawptr {
  return raw::memcpy(dst, src, len)
}

fn memmove(dst rawptr, src const *Void, len u64) rawptr {
  return raw::memmove(dst, src, len)
}

fn memset(dst rawptr, c i32, len u64) rawptr {
  return raw::memset(dst, c, len)
}

fn memset4(dst rawptr, val u32, dwords u64) rawptr {
  return raw::memset4(dst, val, dwords)
}

fn memcmp(s1 const *Void, s2 const *Void, len u64) i32 {
  return raw::memcmp(s1, s2, len)
}

fn srand(seed u64) {
  raw::srand(seed)
}

fn rand(n i32) i32 {
  return raw::rand(n)
}

fn randf() f32 {
  return raw::randf()
}

fn rand_bits() u32 {
  return raw::rand_bits()
}

fn rand_r(state *u64, n i32) i32 {
  return raw::rand_r(state, n)
}

fn randf_r(state *u64) f32 {
  return raw::randf_r(state)
}

fn rand_bits_r(state *u64) u32 {
  return raw::rand_bits_r(state)
}

fn acos(x f64) f64 {
  return raw::acos(x)
}

fn acosf(x f32) f32 {
  return raw::acosf(x)
}

fn asin(x f64) f64 {
  return raw::asin(x)
}

fn asinf(x f32) f32 {
  return raw::asinf(x)
}

fn atan(x f64) f64 {
  return raw::atan(x)
}

fn atanf(x f32) f32 {
  return raw::atanf(x)
}

fn atan2(y f64, x f64) f64 {
  return raw::atan2(y, x)
}

fn atan2f(y f32, x f32) f32 {
  return raw::atan2f(y, x)
}

fn ceil(x f64) f64 {
  return raw::ceil(x)
}

fn ceilf(x f32) f32 {
  return raw::ceilf(x)
}

fn copysign(x f64, y f64) f64 {
  return raw::copysign(x, y)
}

fn copysignf(x f32, y f32) f32 {
  return raw::copysignf(x, y)
}

fn cos(x f64) f64 {
  return raw::cos(x)
}

fn cosf(x f32) f32 {
  return raw::cosf(x)
}

fn exp(x f64) f64 {
  return raw::exp(x)
}

fn expf(x f32) f32 {
  return raw::expf(x)
}

fn fabs(x f64) f64 {
  return raw::fabs(x)
}

fn fabsf(x f32) f32 {
  return raw::fabsf(x)
}

fn floor(x f64) f64 {
  return raw::floor(x)
}

fn floorf(x f32) f32 {
  return raw::floorf(x)
}

fn trunc(x f64) f64 {
  return raw::trunc(x)
}

fn truncf(x f32) f32 {
  return raw::truncf(x)
}

fn fmod(x f64, y f64) f64 {
  return raw::fmod(x, y)
}

fn fmodf(x f32, y f32) f32 {
  return raw::fmodf(x, y)
}

fn isinf(x f64) i32 {
  return raw::isinf(x)
}

fn isinff(x f32) i32 {
  return raw::isinff(x)
}

fn isnan(x f64) i32 {
  return raw::isnan(x)
}

fn isnanf(x f32) i32 {
  return raw::isnanf(x)
}

fn log_f64(x f64) f64 {
  return raw::log_f64(x)
}

fn logf(x f32) f32 {
  return raw::logf(x)
}

fn log10(x f64) f64 {
  return raw::log10(x)
}

fn log10f(x f32) f32 {
  return raw::log10f(x)
}

fn pow(x f64, y f64) f64 {
  return raw::pow(x, y)
}

fn powf(x f32, y f32) f32 {
  return raw::powf(x, y)
}

fn round(x f64) f64 {
  return raw::round(x)
}

fn roundf(x f32) f32 {
  return raw::roundf(x)
}

fn scalbn(x f64, n i32) f64 {
  return raw::scalbn(x, n)
}

fn scalbnf(x f32, n i32) f32 {
  return raw::scalbnf(x, n)
}

fn sin(x f64) f64 {
  return raw::sin(x)
}

fn sinf(x f32) f32 {
  return raw::sinf(x)
}

fn sqrt(x f64) f64 {
  return raw::sqrt(x)
}

fn sqrtf(x f32) f32 {
  return raw::sqrtf(x)
}

fn tan(x f64) f64 {
  return raw::tan(x)
}

fn tanf(x f32) f32 {
  return raw::tanf(x)
}

fn iconv_close(cd iconv_t) i32 {
  return raw::iconv_close(cd)
}

fn get_default_assertion_handler() AssertionHandler {
  return raw::get_default_assertion_handler()
}

fn get_assertion_handler(puserdata *rawptr) AssertionHandler {
  return raw::get_assertion_handler(puserdata)
}

fn get_assertion_report() const *AssertData {
  return raw::get_assertion_report()
}

fn reset_assertion_report() {
  raw::reset_assertion_report()
}

fn get_async_io_size(asyncio *AsyncIO) i64 {
  return raw::get_async_io_size(asyncio)
}

fn read_async_io(asyncio *AsyncIO, ptr rawptr, offset u64, size u64, queue *AsyncIOQueue, userdata rawptr) bool {
  return raw::read_async_io(asyncio, ptr, offset, size, queue, userdata)
}

fn write_async_io(asyncio *AsyncIO, ptr rawptr, offset u64, size u64, queue *AsyncIOQueue, userdata rawptr) bool {
  return raw::write_async_io(asyncio, ptr, offset, size, queue, userdata)
}

fn close_async_io(asyncio *AsyncIO, flush bool, queue *AsyncIOQueue, userdata rawptr) bool {
  return raw::close_async_io(asyncio, flush, queue, userdata)
}

fn create_async_io_queue() *AsyncIOQueue {
  return raw::create_async_io_queue()
}

fn destroy_async_io_queue(queue *AsyncIOQueue) {
  raw::destroy_async_io_queue(queue)
}

fn get_async_io_result(queue *AsyncIOQueue, outcome *AsyncIOOutcome) bool {
  return raw::get_async_io_result(queue, outcome)
}

fn wait_async_io_result(queue *AsyncIOQueue, outcome *AsyncIOOutcome, timeoutMS i32) bool {
  return raw::wait_async_io_result(queue, outcome, timeoutMS)
}

fn signal_async_io_queue(queue *AsyncIOQueue) {
  raw::signal_async_io_queue(queue)
}

fn try_lock_spinlock(lock *SpinLock) bool {
  return raw::try_lock_spinlock(lock)
}

fn lock_spinlock(lock *SpinLock) {
  raw::lock_spinlock(lock)
}

fn unlock_spinlock(lock *SpinLock) {
  raw::unlock_spinlock(lock)
}

fn memory_barrier_release_function() {
  raw::memory_barrier_release_function()
}

fn memory_barrier_acquire_function() {
  raw::memory_barrier_acquire_function()
}

fn compare_and_swap_atomic_int(a *AtomicInt, oldval i32, newval i32) bool {
  return raw::compare_and_swap_atomic_int(a, oldval, newval)
}

fn set_atomic_int(a *AtomicInt, v i32) i32 {
  return raw::set_atomic_int(a, v)
}

fn get_atomic_int(a *AtomicInt) i32 {
  return raw::get_atomic_int(a)
}

fn add_atomic_int(a *AtomicInt, v i32) i32 {
  return raw::add_atomic_int(a, v)
}

fn compare_and_swap_atomic_u32(a *AtomicU32, oldval u32, newval u32) bool {
  return raw::compare_and_swap_atomic_u32(a, oldval, newval)
}

fn set_atomic_u32(a *AtomicU32, v u32) u32 {
  return raw::set_atomic_u32(a, v)
}

fn get_atomic_u32(a *AtomicU32) u32 {
  return raw::get_atomic_u32(a)
}

fn compare_and_swap_atomic_pointer(a *rawptr, oldval rawptr, newval rawptr) bool {
  return raw::compare_and_swap_atomic_pointer(a, oldval, newval)
}

fn set_atomic_pointer(a *rawptr, v rawptr) rawptr {
  return raw::set_atomic_pointer(a, v)
}

fn get_atomic_pointer(a *rawptr) rawptr {
  return raw::get_atomic_pointer(a)
}

fn swap_float(x f32) f32 {
  return raw::swap_float(x)
}

fn out_of_memory() bool {
  return raw::out_of_memory()
}

fn get_error() string {
  result := raw::get_error()
  if result == nil {
    return ""
  }
  return strings::from_c(result)
}

fn clear_error() bool {
  return raw::clear_error()
}

fn get_global_properties() PropertiesID {
  return raw::get_global_properties()
}

fn create_properties() PropertiesID {
  return raw::create_properties()
}

fn copy_properties(src PropertiesID, dst PropertiesID) bool {
  return raw::copy_properties(src, dst)
}

fn lock_properties(props PropertiesID) bool {
  return raw::lock_properties(props)
}

fn unlock_properties(props PropertiesID) {
  raw::unlock_properties(props)
}

fn destroy_properties(props PropertiesID) {
  raw::destroy_properties(props)
}

fn get_thread_name(thread *Thread) OptionalString {
  result := raw::get_thread_name(thread)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_current_thread_id() ThreadID {
  return raw::get_current_thread_id()
}

fn get_thread_id(thread *Thread) ThreadID {
  return raw::get_thread_id(thread)
}

fn set_current_thread_priority(priority ThreadPriority) bool {
  return raw::set_current_thread_priority(priority)
}

fn get_thread_state(thread *Thread) ThreadState {
  return raw::get_thread_state(thread)
}

fn detach_thread(thread *Thread) {
  raw::detach_thread(thread)
}

fn get_tls(id *TLSID) rawptr {
  return raw::get_tls(id)
}

fn cleanup_tls() {
  raw::cleanup_tls()
}

fn create_mutex() *Mutex {
  return raw::create_mutex()
}

fn lock_mutex(mutex *Mutex) {
  raw::lock_mutex(mutex)
}

fn try_lock_mutex(mutex *Mutex) bool {
  return raw::try_lock_mutex(mutex)
}

fn unlock_mutex(mutex *Mutex) {
  raw::unlock_mutex(mutex)
}

fn destroy_mutex(mutex *Mutex) {
  raw::destroy_mutex(mutex)
}

fn create_rw_lock() *RWLock {
  return raw::create_rw_lock()
}

fn lock_rw_lock_for_reading(rwlock *RWLock) {
  raw::lock_rw_lock_for_reading(rwlock)
}

fn lock_rw_lock_for_writing(rwlock *RWLock) {
  raw::lock_rw_lock_for_writing(rwlock)
}

fn try_lock_rw_lock_for_reading(rwlock *RWLock) bool {
  return raw::try_lock_rw_lock_for_reading(rwlock)
}

fn try_lock_rw_lock_for_writing(rwlock *RWLock) bool {
  return raw::try_lock_rw_lock_for_writing(rwlock)
}

fn unlock_rw_lock(rwlock *RWLock) {
  raw::unlock_rw_lock(rwlock)
}

fn destroy_rw_lock(rwlock *RWLock) {
  raw::destroy_rw_lock(rwlock)
}

fn create_semaphore(initial_value u32) *Semaphore {
  return raw::create_semaphore(initial_value)
}

fn destroy_semaphore(sem *Semaphore) {
  raw::destroy_semaphore(sem)
}

fn wait_semaphore(sem *Semaphore) {
  raw::wait_semaphore(sem)
}

fn try_wait_semaphore(sem *Semaphore) bool {
  return raw::try_wait_semaphore(sem)
}

fn wait_semaphore_timeout(sem *Semaphore, timeoutMS i32) bool {
  return raw::wait_semaphore_timeout(sem, timeoutMS)
}

fn signal_semaphore(sem *Semaphore) {
  raw::signal_semaphore(sem)
}

fn get_semaphore_value(sem *Semaphore) u32 {
  return raw::get_semaphore_value(sem)
}

fn create_condition() *Condition {
  return raw::create_condition()
}

fn destroy_condition(cond *Condition) {
  raw::destroy_condition(cond)
}

fn signal_condition(cond *Condition) {
  raw::signal_condition(cond)
}

fn broadcast_condition(cond *Condition) {
  raw::broadcast_condition(cond)
}

fn wait_condition(cond *Condition, mutex *Mutex) {
  raw::wait_condition(cond, mutex)
}

fn wait_condition_timeout(cond *Condition, mutex *Mutex, timeoutMS i32) bool {
  return raw::wait_condition_timeout(cond, mutex, timeoutMS)
}

fn should_init(state *InitState) bool {
  return raw::should_init(state)
}

fn should_quit(state *InitState) bool {
  return raw::should_quit(state)
}

fn set_initialized(state *InitState, initialized bool) {
  raw::set_initialized(state, initialized)
}

fn io_from_const_mem(mem const *Void, size u64) *IOStream {
  return raw::io_from_const_mem(mem, size)
}

fn io_from_dynamic_mem() *IOStream {
  return raw::io_from_dynamic_mem()
}

fn open_io(iface const *IOStreamInterface, userdata rawptr) *IOStream {
  return raw::open_io(iface, userdata)
}

fn close_io(context *IOStream) bool {
  return raw::close_io(context)
}

fn get_io_properties(context *IOStream) PropertiesID {
  return raw::get_io_properties(context)
}

fn get_io_status(context *IOStream) IOStatus {
  return raw::get_io_status(context)
}

fn get_io_size(context *IOStream) i64 {
  return raw::get_io_size(context)
}

fn seek_io(context *IOStream, offset i64, whence IOWhence) i64 {
  return raw::seek_io(context, offset, whence)
}

fn tell_io(context *IOStream) i64 {
  return raw::tell_io(context)
}

fn write_io(context *IOStream, ptr const *Void, size u64) u64 {
  return raw::write_io(context, ptr, size)
}

fn flush_io(context *IOStream) bool {
  return raw::flush_io(context)
}

fn save_file_io(src *IOStream, data const *Void, datasize u64, closeio bool) bool {
  return raw::save_file_io(src, data, datasize, closeio)
}

fn read_u8(src *IOStream, value *u8) bool {
  return raw::read_u8(src, value)
}

fn read_s8(src *IOStream, value *i8) bool {
  return raw::read_s8(src, value)
}

fn read_u16_le(src *IOStream, value *u16) bool {
  return raw::read_u16_le(src, value)
}

fn read_s16_le(src *IOStream, value *i16) bool {
  return raw::read_s16_le(src, value)
}

fn read_u16_be(src *IOStream, value *u16) bool {
  return raw::read_u16_be(src, value)
}

fn read_s16_be(src *IOStream, value *i16) bool {
  return raw::read_s16_be(src, value)
}

fn read_u32_le(src *IOStream, value *u32) bool {
  return raw::read_u32_le(src, value)
}

fn read_s32_le(src *IOStream, value *i32) bool {
  return raw::read_s32_le(src, value)
}

fn read_u32_be(src *IOStream, value *u32) bool {
  return raw::read_u32_be(src, value)
}

fn read_s32_be(src *IOStream, value *i32) bool {
  return raw::read_s32_be(src, value)
}

fn read_u64_le(src *IOStream, value *u64) bool {
  return raw::read_u64_le(src, value)
}

fn read_s64_le(src *IOStream, value *i64) bool {
  return raw::read_s64_le(src, value)
}

fn read_u64_be(src *IOStream, value *u64) bool {
  return raw::read_u64_be(src, value)
}

fn read_s64_be(src *IOStream, value *i64) bool {
  return raw::read_s64_be(src, value)
}

fn write_u8(dst *IOStream, value u8) bool {
  return raw::write_u8(dst, value)
}

fn write_s8(dst *IOStream, value i8) bool {
  return raw::write_s8(dst, value)
}

fn write_u16_le(dst *IOStream, value u16) bool {
  return raw::write_u16_le(dst, value)
}

fn write_s16_le(dst *IOStream, value i16) bool {
  return raw::write_s16_le(dst, value)
}

fn write_u16_be(dst *IOStream, value u16) bool {
  return raw::write_u16_be(dst, value)
}

fn write_s16_be(dst *IOStream, value i16) bool {
  return raw::write_s16_be(dst, value)
}

fn write_u32_le(dst *IOStream, value u32) bool {
  return raw::write_u32_le(dst, value)
}

fn write_s32_le(dst *IOStream, value i32) bool {
  return raw::write_s32_le(dst, value)
}

fn write_u32_be(dst *IOStream, value u32) bool {
  return raw::write_u32_be(dst, value)
}

fn write_s32_be(dst *IOStream, value i32) bool {
  return raw::write_s32_be(dst, value)
}

fn write_u64_le(dst *IOStream, value u64) bool {
  return raw::write_u64_le(dst, value)
}

fn write_s64_le(dst *IOStream, value i64) bool {
  return raw::write_s64_le(dst, value)
}

fn write_u64_be(dst *IOStream, value u64) bool {
  return raw::write_u64_be(dst, value)
}

fn write_s64_be(dst *IOStream, value i64) bool {
  return raw::write_s64_be(dst, value)
}

fn get_num_audio_drivers() i32 {
  return raw::get_num_audio_drivers()
}

fn get_audio_driver(index i32) OptionalString {
  result := raw::get_audio_driver(index)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_current_audio_driver() OptionalString {
  result := raw::get_current_audio_driver()
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_audio_device_name(devid AudioDeviceID) OptionalString {
  result := raw::get_audio_device_name(devid)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn open_audio_device(devid AudioDeviceID, spec const *AudioSpec) AudioDeviceID {
  return raw::open_audio_device(devid, spec)
}

fn is_audio_device_physical(devid AudioDeviceID) bool {
  return raw::is_audio_device_physical(devid)
}

fn is_audio_device_playback(devid AudioDeviceID) bool {
  return raw::is_audio_device_playback(devid)
}

fn pause_audio_device(devid AudioDeviceID) bool {
  return raw::pause_audio_device(devid)
}

fn resume_audio_device(devid AudioDeviceID) bool {
  return raw::resume_audio_device(devid)
}

fn audio_device_paused(devid AudioDeviceID) bool {
  return raw::audio_device_paused(devid)
}

fn get_audio_device_gain(devid AudioDeviceID) f32 {
  return raw::get_audio_device_gain(devid)
}

fn set_audio_device_gain(devid AudioDeviceID, gain f32) bool {
  return raw::set_audio_device_gain(devid, gain)
}

fn close_audio_device(devid AudioDeviceID) {
  raw::close_audio_device(devid)
}

fn bind_audio_streams(devid AudioDeviceID, streams const **AudioStream, num_streams i32) bool {
  return raw::bind_audio_streams(devid, streams, num_streams)
}

fn bind_audio_stream(devid AudioDeviceID, stream *AudioStream) bool {
  return raw::bind_audio_stream(devid, stream)
}

fn unbind_audio_streams(streams const **AudioStream, num_streams i32) {
  raw::unbind_audio_streams(streams, num_streams)
}

fn unbind_audio_stream(stream *AudioStream) {
  raw::unbind_audio_stream(stream)
}

fn get_audio_stream_device(stream *AudioStream) AudioDeviceID {
  return raw::get_audio_stream_device(stream)
}

fn create_audio_stream(src_spec const *AudioSpec, dst_spec const *AudioSpec) *AudioStream {
  return raw::create_audio_stream(src_spec, dst_spec)
}

fn get_audio_stream_properties(stream *AudioStream) PropertiesID {
  return raw::get_audio_stream_properties(stream)
}

fn get_audio_stream_format(stream *AudioStream, src_spec *AudioSpec, dst_spec *AudioSpec) bool {
  return raw::get_audio_stream_format(stream, src_spec, dst_spec)
}

fn set_audio_stream_format(stream *AudioStream, src_spec const *AudioSpec, dst_spec const *AudioSpec) bool {
  return raw::set_audio_stream_format(stream, src_spec, dst_spec)
}

fn get_audio_stream_frequency_ratio(stream *AudioStream) f32 {
  return raw::get_audio_stream_frequency_ratio(stream)
}

fn set_audio_stream_frequency_ratio(stream *AudioStream, ratio f32) bool {
  return raw::set_audio_stream_frequency_ratio(stream, ratio)
}

fn get_audio_stream_gain(stream *AudioStream) f32 {
  return raw::get_audio_stream_gain(stream)
}

fn set_audio_stream_gain(stream *AudioStream, gain f32) bool {
  return raw::set_audio_stream_gain(stream, gain)
}

fn put_audio_stream_data(stream *AudioStream, buf const *Void, len i32) bool {
  return raw::put_audio_stream_data(stream, buf, len)
}

fn get_audio_stream_data(stream *AudioStream, buf rawptr, len i32) i32 {
  return raw::get_audio_stream_data(stream, buf, len)
}

fn get_audio_stream_available(stream *AudioStream) i32 {
  return raw::get_audio_stream_available(stream)
}

fn get_audio_stream_queued(stream *AudioStream) i32 {
  return raw::get_audio_stream_queued(stream)
}

fn flush_audio_stream(stream *AudioStream) bool {
  return raw::flush_audio_stream(stream)
}

fn clear_audio_stream(stream *AudioStream) bool {
  return raw::clear_audio_stream(stream)
}

fn pause_audio_stream_device(stream *AudioStream) bool {
  return raw::pause_audio_stream_device(stream)
}

fn resume_audio_stream_device(stream *AudioStream) bool {
  return raw::resume_audio_stream_device(stream)
}

fn audio_stream_device_paused(stream *AudioStream) bool {
  return raw::audio_stream_device_paused(stream)
}

fn lock_audio_stream(stream *AudioStream) bool {
  return raw::lock_audio_stream(stream)
}

fn unlock_audio_stream(stream *AudioStream) bool {
  return raw::unlock_audio_stream(stream)
}

fn destroy_audio_stream(stream *AudioStream) {
  raw::destroy_audio_stream(stream)
}

fn load_wav_io(src *IOStream, closeio bool, spec *AudioSpec, audio_buf **u8, audio_len *u32) bool {
  return raw::load_wav_io(src, closeio, spec, audio_buf, audio_len)
}

fn mix_audio(dst *u8, src const *u8, format AudioFormat, len u32, volume f32) bool {
  return raw::mix_audio(dst, src, format, len, volume)
}

fn get_audio_format_name(format AudioFormat) string {
  result := raw::get_audio_format_name(format)
  if result == nil {
    return ""
  }
  return strings::from_c(result)
}

fn get_silence_value_for_format(format AudioFormat) i32 {
  return raw::get_silence_value_for_format(format)
}

fn most_significant_bit_index32(x u32) i32 {
  return raw::most_significant_bit_index32(x)
}

fn has_exactly_one_bit_set32(x u32) bool {
  return raw::has_exactly_one_bit_set32(x)
}

fn compose_custom_blend_mode(srcColorFactor BlendFactor, dstColorFactor BlendFactor, colorOperation BlendOperation, srcAlphaFactor BlendFactor, dstAlphaFactor BlendFactor, alphaOperation BlendOperation) BlendMode {
  return raw::compose_custom_blend_mode(srcColorFactor, dstColorFactor, colorOperation, srcAlphaFactor, dstAlphaFactor, alphaOperation)
}

fn get_pixel_format_name(format PixelFormat) string {
  result := raw::get_pixel_format_name(format)
  if result == nil {
    return ""
  }
  return strings::from_c(result)
}

fn get_pixel_format_for_masks(bpp i32, Rmask u32, Gmask u32, Bmask u32, Amask u32) PixelFormat {
  return raw::get_pixel_format_for_masks(bpp, Rmask, Gmask, Bmask, Amask)
}

fn get_pixel_format_details(format PixelFormat) const *PixelFormatDetails {
  return raw::get_pixel_format_details(format)
}

fn create_palette(ncolors i32) *Palette {
  return raw::create_palette(ncolors)
}

fn set_palette_colors(palette *Palette, colors const *Color, firstcolor i32, ncolors i32) bool {
  return raw::set_palette_colors(palette, colors, firstcolor, ncolors)
}

fn destroy_palette(palette *Palette) {
  raw::destroy_palette(palette)
}

fn map_rgb(format const *PixelFormatDetails, palette const *Palette, r u8, g u8, b u8) u32 {
  return raw::map_rgb(format, palette, r, g, b)
}

fn map_rgba(format const *PixelFormatDetails, palette const *Palette, r u8, g u8, b u8, a u8) u32 {
  return raw::map_rgba(format, palette, r, g, b, a)
}

fn get_rgb(pixel u32, format const *PixelFormatDetails, palette const *Palette, r *u8, g *u8, b *u8) {
  raw::get_rgb(pixel, format, palette, r, g, b)
}

fn get_rgba(pixel u32, format const *PixelFormatDetails, palette const *Palette, r *u8, g *u8, b *u8, a *u8) {
  raw::get_rgba(pixel, format, palette, r, g, b, a)
}

fn rect_to_f_rect(rect const *Rect, frect *FRect) {
  raw::rect_to_f_rect(rect, frect)
}

fn point_in_rect(p const *Point, r const *Rect) bool {
  return raw::point_in_rect(p, r)
}

fn rect_empty(r const *Rect) bool {
  return raw::rect_empty(r)
}

fn rects_equal(a const *Rect, b const *Rect) bool {
  return raw::rects_equal(a, b)
}

fn has_rect_intersection(A const *Rect, B const *Rect) bool {
  return raw::has_rect_intersection(A, B)
}

fn get_rect_intersection(A const *Rect, B const *Rect, result *Rect) bool {
  return raw::get_rect_intersection(A, B, result)
}

fn get_rect_union(A const *Rect, B const *Rect, result *Rect) bool {
  return raw::get_rect_union(A, B, result)
}

fn get_rect_enclosing_points(points const *Point, count i32, clip const *Rect, result *Rect) bool {
  return raw::get_rect_enclosing_points(points, count, clip, result)
}

fn point_in_rect_float(p const *FPoint, r const *FRect) bool {
  return raw::point_in_rect_float(p, r)
}

fn rect_empty_float(r const *FRect) bool {
  return raw::rect_empty_float(r)
}

fn rects_equal_epsilon(a const *FRect, b const *FRect, epsilon f32) bool {
  return raw::rects_equal_epsilon(a, b, epsilon)
}

fn rects_equal_float(a const *FRect, b const *FRect) bool {
  return raw::rects_equal_float(a, b)
}

fn has_rect_intersection_float(A const *FRect, B const *FRect) bool {
  return raw::has_rect_intersection_float(A, B)
}

fn get_rect_intersection_float(A const *FRect, B const *FRect, result *FRect) bool {
  return raw::get_rect_intersection_float(A, B, result)
}

fn get_rect_union_float(A const *FRect, B const *FRect, result *FRect) bool {
  return raw::get_rect_union_float(A, B, result)
}

fn get_rect_enclosing_points_float(points const *FPoint, count i32, clip const *FRect, result *FRect) bool {
  return raw::get_rect_enclosing_points_float(points, count, clip, result)
}

fn create_surface(width i32, height i32, format PixelFormat) *Surface {
  return raw::create_surface(width, height, format)
}

fn create_surface_from(width i32, height i32, format PixelFormat, pixels rawptr, pitch i32) *Surface {
  return raw::create_surface_from(width, height, format, pixels, pitch)
}

fn destroy_surface(surface *Surface) {
  raw::destroy_surface(surface)
}

fn get_surface_properties(surface *Surface) PropertiesID {
  return raw::get_surface_properties(surface)
}

fn set_surface_colorspace(surface *Surface, colorspace Colorspace) bool {
  return raw::set_surface_colorspace(surface, colorspace)
}

fn get_surface_colorspace(surface *Surface) Colorspace {
  return raw::get_surface_colorspace(surface)
}

fn create_surface_palette(surface *Surface) *Palette {
  return raw::create_surface_palette(surface)
}

fn set_surface_palette(surface *Surface, palette *Palette) bool {
  return raw::set_surface_palette(surface, palette)
}

fn get_surface_palette(surface *Surface) *Palette {
  return raw::get_surface_palette(surface)
}

fn add_surface_alternate_image(surface *Surface, image *Surface) bool {
  return raw::add_surface_alternate_image(surface, image)
}

fn surface_has_alternate_images(surface *Surface) bool {
  return raw::surface_has_alternate_images(surface)
}

fn remove_surface_alternate_images(surface *Surface) {
  raw::remove_surface_alternate_images(surface)
}

fn lock_surface(surface *Surface) bool {
  return raw::lock_surface(surface)
}

fn unlock_surface(surface *Surface) {
  raw::unlock_surface(surface)
}

fn load_bmp_io(src *IOStream, closeio bool) *Surface {
  return raw::load_bmp_io(src, closeio)
}

fn load_bmp(file string) (*Surface, error) {
  file_c, file_err := strings::to_c(file)
  if file_err != nil {
    return nil, file_err
  }
  defer strings::free_c(&file_c)
  result := raw::load_bmp(strings::c_view(&file_c))
  if result == nil {
    return nil, error("SDL load_bmp failed")
  }
  return result, nil
}

fn save_bmp_io(surface *Surface, dst *IOStream, closeio bool) bool {
  return raw::save_bmp_io(surface, dst, closeio)
}

fn save_bmp(surface *Surface, file string) error {
  file_c, file_err := strings::to_c(file)
  if file_err != nil {
    return file_err
  }
  defer strings::free_c(&file_c)
  if raw::save_bmp(surface, strings::c_view(&file_c)) == false {
    return error("SDL save_bmp failed")
  }
  return nil
}

fn set_surface_rle(surface *Surface, enabled bool) bool {
  return raw::set_surface_rle(surface, enabled)
}

fn surface_has_rle(surface *Surface) bool {
  return raw::surface_has_rle(surface)
}

fn set_surface_color_key(surface *Surface, enabled bool, key u32) bool {
  return raw::set_surface_color_key(surface, enabled, key)
}

fn surface_has_color_key(surface *Surface) bool {
  return raw::surface_has_color_key(surface)
}

fn get_surface_color_key(surface *Surface, key *u32) bool {
  return raw::get_surface_color_key(surface, key)
}

fn set_surface_color_mod(surface *Surface, r u8, g u8, b u8) bool {
  return raw::set_surface_color_mod(surface, r, g, b)
}

fn get_surface_color_mod(surface *Surface, r *u8, g *u8, b *u8) bool {
  return raw::get_surface_color_mod(surface, r, g, b)
}

fn set_surface_alpha_mod(surface *Surface, alpha u8) bool {
  return raw::set_surface_alpha_mod(surface, alpha)
}

fn get_surface_alpha_mod(surface *Surface, alpha *u8) bool {
  return raw::get_surface_alpha_mod(surface, alpha)
}

fn set_surface_blend_mode(surface *Surface, blendMode BlendMode) bool {
  return raw::set_surface_blend_mode(surface, blendMode)
}

fn get_surface_blend_mode(surface *Surface, blendMode *BlendMode) bool {
  return raw::get_surface_blend_mode(surface, blendMode)
}

fn set_surface_clip_rect(surface *Surface, rect const *Rect) bool {
  return raw::set_surface_clip_rect(surface, rect)
}

fn get_surface_clip_rect(surface *Surface, rect *Rect) bool {
  return raw::get_surface_clip_rect(surface, rect)
}

fn flip_surface(surface *Surface, flip FlipMode) bool {
  return raw::flip_surface(surface, flip)
}

fn duplicate_surface(surface *Surface) *Surface {
  return raw::duplicate_surface(surface)
}

fn scale_surface(surface *Surface, width i32, height i32, scaleMode ScaleMode) *Surface {
  return raw::scale_surface(surface, width, height, scaleMode)
}

fn convert_surface(surface *Surface, format PixelFormat) *Surface {
  return raw::convert_surface(surface, format)
}

fn convert_surface_and_colorspace(surface *Surface, format PixelFormat, palette *Palette, colorspace Colorspace, props PropertiesID) *Surface {
  return raw::convert_surface_and_colorspace(surface, format, palette, colorspace, props)
}

fn convert_pixels(width i32, height i32, src_format PixelFormat, src const *Void, src_pitch i32, dst_format PixelFormat, dst rawptr, dst_pitch i32) bool {
  return raw::convert_pixels(width, height, src_format, src, src_pitch, dst_format, dst, dst_pitch)
}

fn convert_pixels_and_colorspace(width i32, height i32, src_format PixelFormat, src_colorspace Colorspace, src_properties PropertiesID, src const *Void, src_pitch i32, dst_format PixelFormat, dst_colorspace Colorspace, dst_properties PropertiesID, dst rawptr, dst_pitch i32) bool {
  return raw::convert_pixels_and_colorspace(width, height, src_format, src_colorspace, src_properties, src, src_pitch, dst_format, dst_colorspace, dst_properties, dst, dst_pitch)
}

fn premultiply_alpha(width i32, height i32, src_format PixelFormat, src const *Void, src_pitch i32, dst_format PixelFormat, dst rawptr, dst_pitch i32, linear bool) bool {
  return raw::premultiply_alpha(width, height, src_format, src, src_pitch, dst_format, dst, dst_pitch, linear)
}

fn premultiply_surface_alpha(surface *Surface, linear bool) bool {
  return raw::premultiply_surface_alpha(surface, linear)
}

fn clear_surface(surface *Surface, r f32, g f32, b f32, a f32) bool {
  return raw::clear_surface(surface, r, g, b, a)
}

fn fill_surface_rect(dst *Surface, rect const *Rect, color u32) bool {
  return raw::fill_surface_rect(dst, rect, color)
}

fn fill_surface_rects(dst *Surface, rects const *Rect, count i32, color u32) bool {
  return raw::fill_surface_rects(dst, rects, count, color)
}

fn blit_surface(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect) bool {
  return raw::blit_surface(src, srcrect, dst, dstrect)
}

fn blit_surface_unchecked(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect) bool {
  return raw::blit_surface_unchecked(src, srcrect, dst, dstrect)
}

fn blit_surface_scaled(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect, scaleMode ScaleMode) bool {
  return raw::blit_surface_scaled(src, srcrect, dst, dstrect, scaleMode)
}

fn blit_surface_unchecked_scaled(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect, scaleMode ScaleMode) bool {
  return raw::blit_surface_unchecked_scaled(src, srcrect, dst, dstrect, scaleMode)
}

fn stretch_surface(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect, scaleMode ScaleMode) bool {
  return raw::stretch_surface(src, srcrect, dst, dstrect, scaleMode)
}

fn blit_surface_tiled(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect) bool {
  return raw::blit_surface_tiled(src, srcrect, dst, dstrect)
}

fn blit_surface_tiled_with_scale(src *Surface, srcrect const *Rect, scale f32, scaleMode ScaleMode, dst *Surface, dstrect const *Rect) bool {
  return raw::blit_surface_tiled_with_scale(src, srcrect, scale, scaleMode, dst, dstrect)
}

fn blit_surface9_grid(src *Surface, srcrect const *Rect, left_width i32, right_width i32, top_height i32, bottom_height i32, scale f32, scaleMode ScaleMode, dst *Surface, dstrect const *Rect) bool {
  return raw::blit_surface9_grid(src, srcrect, left_width, right_width, top_height, bottom_height, scale, scaleMode, dst, dstrect)
}

fn map_surface_rgb(surface *Surface, r u8, g u8, b u8) u32 {
  return raw::map_surface_rgb(surface, r, g, b)
}

fn map_surface_rgba(surface *Surface, r u8, g u8, b u8, a u8) u32 {
  return raw::map_surface_rgba(surface, r, g, b, a)
}

fn read_surface_pixel(surface *Surface, x i32, y i32, r *u8, g *u8, b *u8, a *u8) bool {
  return raw::read_surface_pixel(surface, x, y, r, g, b, a)
}

fn write_surface_pixel(surface *Surface, x i32, y i32, r u8, g u8, b u8, a u8) bool {
  return raw::write_surface_pixel(surface, x, y, r, g, b, a)
}

fn write_surface_pixel_float(surface *Surface, x i32, y i32, r f32, g f32, b f32, a f32) bool {
  return raw::write_surface_pixel_float(surface, x, y, r, g, b, a)
}

fn get_num_camera_drivers() i32 {
  return raw::get_num_camera_drivers()
}

fn get_camera_driver(index i32) OptionalString {
  result := raw::get_camera_driver(index)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_current_camera_driver() OptionalString {
  result := raw::get_current_camera_driver()
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_camera_name(instance_id CameraID) OptionalString {
  result := raw::get_camera_name(instance_id)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_camera_position(instance_id CameraID) CameraPosition {
  return raw::get_camera_position(instance_id)
}

fn open_camera(instance_id CameraID, spec const *CameraSpec) *Camera {
  return raw::open_camera(instance_id, spec)
}

fn get_camera_permission_state(camera *Camera) i32 {
  return raw::get_camera_permission_state(camera)
}

fn get_camera_id(camera *Camera) CameraID {
  return raw::get_camera_id(camera)
}

fn get_camera_properties(camera *Camera) PropertiesID {
  return raw::get_camera_properties(camera)
}

fn get_camera_format(camera *Camera, spec *CameraSpec) bool {
  return raw::get_camera_format(camera, spec)
}

fn acquire_camera_frame(camera *Camera, timestampNS *u64) *Surface {
  return raw::acquire_camera_frame(camera, timestampNS)
}

fn release_camera_frame(camera *Camera, frame *Surface) {
  raw::release_camera_frame(camera, frame)
}

fn close_camera(camera *Camera) {
  raw::close_camera(camera)
}

fn set_clipboard_text(text string) error {
  text_c, text_err := strings::to_c(text)
  if text_err != nil {
    return text_err
  }
  defer strings::free_c(&text_c)
  if raw::set_clipboard_text(strings::c_view(&text_c)) == false {
    return error("SDL set_clipboard_text failed")
  }
  return nil
}

fn get_clipboard_text() (strings::Owned, error) {
  result := raw::get_clipboard_text()
  if result == nil {
    return strings::Owned{text: ""}, error("SDL get_clipboard_text failed")
  }
  defer raw::free_text(result)
  return strings::clone(strings::from_c(result))
}

fn has_clipboard_text() bool {
  return raw::has_clipboard_text()
}

fn set_primary_selection_text(text string) error {
  text_c, text_err := strings::to_c(text)
  if text_err != nil {
    return text_err
  }
  defer strings::free_c(&text_c)
  if raw::set_primary_selection_text(strings::c_view(&text_c)) == false {
    return error("SDL set_primary_selection_text failed")
  }
  return nil
}

fn get_primary_selection_text() (strings::Owned, error) {
  result := raw::get_primary_selection_text()
  if result == nil {
    return strings::Owned{text: ""}, error("SDL get_primary_selection_text failed")
  }
  defer raw::free_text(result)
  return strings::clone(strings::from_c(result))
}

fn has_primary_selection_text() bool {
  return raw::has_primary_selection_text()
}

fn clear_clipboard_data() bool {
  return raw::clear_clipboard_data()
}

fn get_num_logical_cpu_cores() i32 {
  return raw::get_num_logical_cpu_cores()
}

fn get_cpu_cache_line_size() i32 {
  return raw::get_cpu_cache_line_size()
}

fn has_alti_vec() bool {
  return raw::has_alti_vec()
}

fn has_mmx() bool {
  return raw::has_mmx()
}

fn has_sse() bool {
  return raw::has_sse()
}

fn has_sse2() bool {
  return raw::has_sse2()
}

fn has_sse3() bool {
  return raw::has_sse3()
}

fn has_sse41() bool {
  return raw::has_sse41()
}

fn has_sse42() bool {
  return raw::has_sse42()
}

fn has_avx() bool {
  return raw::has_avx()
}

fn has_avx2() bool {
  return raw::has_avx2()
}

fn has_avx512_f() bool {
  return raw::has_avx512_f()
}

fn has_armsimd() bool {
  return raw::has_armsimd()
}

fn has_neon() bool {
  return raw::has_neon()
}

fn has_lsx() bool {
  return raw::has_lsx()
}

fn has_lasx() bool {
  return raw::has_lasx()
}

fn get_system_ram() i32 {
  return raw::get_system_ram()
}

fn get_simd_alignment() u64 {
  return raw::get_simd_alignment()
}

fn get_num_video_drivers() i32 {
  return raw::get_num_video_drivers()
}

fn get_video_driver(index i32) OptionalString {
  result := raw::get_video_driver(index)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_current_video_driver() OptionalString {
  result := raw::get_current_video_driver()
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_system_theme() SystemTheme {
  return raw::get_system_theme()
}

fn get_primary_display() DisplayID {
  return raw::get_primary_display()
}

fn get_display_properties(displayID DisplayID) PropertiesID {
  return raw::get_display_properties(displayID)
}

fn get_display_name(displayID DisplayID) OptionalString {
  result := raw::get_display_name(displayID)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_display_bounds(displayID DisplayID, rect *Rect) bool {
  return raw::get_display_bounds(displayID, rect)
}

fn get_display_usable_bounds(displayID DisplayID, rect *Rect) bool {
  return raw::get_display_usable_bounds(displayID, rect)
}

fn get_natural_display_orientation(displayID DisplayID) DisplayOrientation {
  return raw::get_natural_display_orientation(displayID)
}

fn get_current_display_orientation(displayID DisplayID) DisplayOrientation {
  return raw::get_current_display_orientation(displayID)
}

fn get_display_content_scale(displayID DisplayID) f32 {
  return raw::get_display_content_scale(displayID)
}

fn get_closest_fullscreen_display_mode(displayID DisplayID, w i32, h i32, refresh_rate f32, include_high_density_modes bool, closest *DisplayMode) bool {
  return raw::get_closest_fullscreen_display_mode(displayID, w, h, refresh_rate, include_high_density_modes, closest)
}

fn get_desktop_display_mode(displayID DisplayID) const *DisplayMode {
  return raw::get_desktop_display_mode(displayID)
}

fn get_current_display_mode(displayID DisplayID) const *DisplayMode {
  return raw::get_current_display_mode(displayID)
}

fn get_display_for_point(point const *Point) DisplayID {
  return raw::get_display_for_point(point)
}

fn get_display_for_rect(rect const *Rect) DisplayID {
  return raw::get_display_for_rect(rect)
}

fn get_display_for_window(window *Window) DisplayID {
  return raw::get_display_for_window(window)
}

fn get_window_pixel_density(window *Window) f32 {
  return raw::get_window_pixel_density(window)
}

fn get_window_display_scale(window *Window) f32 {
  return raw::get_window_display_scale(window)
}

fn set_window_fullscreen_mode(window *Window, mode const *DisplayMode) bool {
  return raw::set_window_fullscreen_mode(window, mode)
}

fn get_window_fullscreen_mode(window *Window) const *DisplayMode {
  return raw::get_window_fullscreen_mode(window)
}

fn get_window_pixel_format(window *Window) PixelFormat {
  return raw::get_window_pixel_format(window)
}

fn create_window(title string, w i32, h i32, flags WindowFlags) (*Window, error) {
  title_c, title_err := strings::to_c(title)
  if title_err != nil {
    return nil, title_err
  }
  defer strings::free_c(&title_c)
  result := raw::create_window(strings::c_view(&title_c), w, h, flags)
  if result == nil {
    return nil, error("SDL create_window failed")
  }
  return result, nil
}

fn create_popup_window(parent *Window, offset_x i32, offset_y i32, w i32, h i32, flags WindowFlags) *Window {
  return raw::create_popup_window(parent, offset_x, offset_y, w, h, flags)
}

fn create_window_with_properties(props PropertiesID) *Window {
  return raw::create_window_with_properties(props)
}

fn get_window_id(window *Window) WindowID {
  return raw::get_window_id(window)
}

fn get_window_from_id(id WindowID) *Window {
  return raw::get_window_from_id(id)
}

fn get_window_parent(window *Window) *Window {
  return raw::get_window_parent(window)
}

fn get_window_properties(window *Window) PropertiesID {
  return raw::get_window_properties(window)
}

fn get_window_flags(window *Window) WindowFlags {
  return raw::get_window_flags(window)
}

fn set_window_title(window *Window, title string) error {
  title_c, title_err := strings::to_c(title)
  if title_err != nil {
    return title_err
  }
  defer strings::free_c(&title_c)
  if raw::set_window_title(window, strings::c_view(&title_c)) == false {
    return error("SDL set_window_title failed")
  }
  return nil
}

fn get_window_title(window *Window) string {
  result := raw::get_window_title(window)
  if result == nil {
    return ""
  }
  return strings::from_c(result)
}

fn set_window_icon(window *Window, icon *Surface) bool {
  return raw::set_window_icon(window, icon)
}

fn set_window_position(window *Window, x i32, y i32) bool {
  return raw::set_window_position(window, x, y)
}

fn set_window_size(window *Window, w i32, h i32) bool {
  return raw::set_window_size(window, w, h)
}

fn get_window_safe_area(window *Window, rect *Rect) bool {
  return raw::get_window_safe_area(window, rect)
}

fn set_window_aspect_ratio(window *Window, min_aspect f32, max_aspect f32) bool {
  return raw::set_window_aspect_ratio(window, min_aspect, max_aspect)
}

fn set_window_minimum_size(window *Window, min_w i32, min_h i32) bool {
  return raw::set_window_minimum_size(window, min_w, min_h)
}

fn set_window_maximum_size(window *Window, max_w i32, max_h i32) bool {
  return raw::set_window_maximum_size(window, max_w, max_h)
}

fn set_window_bordered(window *Window, bordered bool) bool {
  return raw::set_window_bordered(window, bordered)
}

fn set_window_resizable(window *Window, resizable bool) bool {
  return raw::set_window_resizable(window, resizable)
}

fn set_window_always_on_top(window *Window, on_top bool) bool {
  return raw::set_window_always_on_top(window, on_top)
}

fn show_window(window *Window) bool {
  return raw::show_window(window)
}

fn hide_window(window *Window) bool {
  return raw::hide_window(window)
}

fn raise_window(window *Window) bool {
  return raw::raise_window(window)
}

fn maximize_window(window *Window) bool {
  return raw::maximize_window(window)
}

fn minimize_window(window *Window) bool {
  return raw::minimize_window(window)
}

fn restore_window(window *Window) bool {
  return raw::restore_window(window)
}

fn set_window_fullscreen(window *Window, fullscreen bool) bool {
  return raw::set_window_fullscreen(window, fullscreen)
}

fn sync_window(window *Window) bool {
  return raw::sync_window(window)
}

fn window_has_surface(window *Window) bool {
  return raw::window_has_surface(window)
}

fn get_window_surface(window *Window) *Surface {
  return raw::get_window_surface(window)
}

fn set_window_surface_v_sync(window *Window, vsync i32) bool {
  return raw::set_window_surface_v_sync(window, vsync)
}

fn update_window_surface(window *Window) bool {
  return raw::update_window_surface(window)
}

fn update_window_surface_rects(window *Window, rects const *Rect, numrects i32) bool {
  return raw::update_window_surface_rects(window, rects, numrects)
}

fn destroy_window_surface(window *Window) bool {
  return raw::destroy_window_surface(window)
}

fn set_window_keyboard_grab(window *Window, grabbed bool) bool {
  return raw::set_window_keyboard_grab(window, grabbed)
}

fn set_window_mouse_grab(window *Window, grabbed bool) bool {
  return raw::set_window_mouse_grab(window, grabbed)
}

fn get_window_keyboard_grab(window *Window) bool {
  return raw::get_window_keyboard_grab(window)
}

fn get_window_mouse_grab(window *Window) bool {
  return raw::get_window_mouse_grab(window)
}

fn get_grabbed_window() *Window {
  return raw::get_grabbed_window()
}

fn set_window_mouse_rect(window *Window, rect const *Rect) bool {
  return raw::set_window_mouse_rect(window, rect)
}

fn get_window_mouse_rect(window *Window) const *Rect {
  return raw::get_window_mouse_rect(window)
}

fn set_window_opacity(window *Window, opacity f32) bool {
  return raw::set_window_opacity(window, opacity)
}

fn get_window_opacity(window *Window) f32 {
  return raw::get_window_opacity(window)
}

fn set_window_parent(window *Window, parent *Window) bool {
  return raw::set_window_parent(window, parent)
}

fn set_window_modal(window *Window, modal bool) bool {
  return raw::set_window_modal(window, modal)
}

fn set_window_focusable(window *Window, focusable bool) bool {
  return raw::set_window_focusable(window, focusable)
}

fn show_window_system_menu(window *Window, x i32, y i32) bool {
  return raw::show_window_system_menu(window, x, y)
}

fn set_window_shape(window *Window, shape *Surface) bool {
  return raw::set_window_shape(window, shape)
}

fn flash_window(window *Window, operation FlashOperation) bool {
  return raw::flash_window(window, operation)
}

fn destroy_window(window *Window) {
  raw::destroy_window(window)
}

fn screen_saver_enabled() bool {
  return raw::screen_saver_enabled()
}

fn enable_screen_saver() bool {
  return raw::enable_screen_saver()
}

fn disable_screen_saver() bool {
  return raw::disable_screen_saver()
}

fn gl_unload_library() {
  raw::gl_unload_library()
}

fn gl_reset_attributes() {
  raw::gl_reset_attributes()
}

fn gl_set_attribute(attr GLAttr, value i32) bool {
  return raw::gl_set_attribute(attr, value)
}

fn gl_create_context(window *Window) GLContext {
  return raw::gl_create_context(window)
}

fn gl_make_current(window *Window, context GLContext) bool {
  return raw::gl_make_current(window, context)
}

fn gl_get_current_window() *Window {
  return raw::gl_get_current_window()
}

fn gl_get_current_context() GLContext {
  return raw::gl_get_current_context()
}

fn egl_get_current_display() EGLDisplay {
  return raw::egl_get_current_display()
}

fn egl_get_current_config() EGLConfig {
  return raw::egl_get_current_config()
}

fn egl_get_window_surface(window *Window) EGLSurface {
  return raw::egl_get_window_surface(window)
}

fn gl_set_swap_interval(interval i32) bool {
  return raw::gl_set_swap_interval(interval)
}

fn gl_swap_window(window *Window) bool {
  return raw::gl_swap_window(window)
}

fn gl_destroy_context(context GLContext) bool {
  return raw::gl_destroy_context(context)
}

fn get_sensor_name_for_id(instance_id SensorID) OptionalString {
  result := raw::get_sensor_name_for_id(instance_id)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_sensor_type_for_id(instance_id SensorID) SensorType {
  return raw::get_sensor_type_for_id(instance_id)
}

fn get_sensor_non_portable_type_for_id(instance_id SensorID) i32 {
  return raw::get_sensor_non_portable_type_for_id(instance_id)
}

fn open_sensor(instance_id SensorID) *Sensor {
  return raw::open_sensor(instance_id)
}

fn get_sensor_from_id(instance_id SensorID) *Sensor {
  return raw::get_sensor_from_id(instance_id)
}

fn get_sensor_properties(sensor *Sensor) PropertiesID {
  return raw::get_sensor_properties(sensor)
}

fn get_sensor_name(sensor *Sensor) OptionalString {
  result := raw::get_sensor_name(sensor)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_sensor_type(sensor *Sensor) SensorType {
  return raw::get_sensor_type(sensor)
}

fn get_sensor_non_portable_type(sensor *Sensor) i32 {
  return raw::get_sensor_non_portable_type(sensor)
}

fn get_sensor_id(sensor *Sensor) SensorID {
  return raw::get_sensor_id(sensor)
}

fn close_sensor(sensor *Sensor) {
  raw::close_sensor(sensor)
}

fn update_sensors() {
  raw::update_sensors()
}

fn lock_joysticks() {
  raw::lock_joysticks()
}

fn unlock_joysticks() {
  raw::unlock_joysticks()
}

fn has_joystick() bool {
  return raw::has_joystick()
}

fn get_joystick_name_for_id(instance_id JoystickID) OptionalString {
  result := raw::get_joystick_name_for_id(instance_id)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_joystick_path_for_id(instance_id JoystickID) OptionalString {
  result := raw::get_joystick_path_for_id(instance_id)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_joystick_player_index_for_id(instance_id JoystickID) i32 {
  return raw::get_joystick_player_index_for_id(instance_id)
}

fn get_joystick_guid_for_id(instance_id JoystickID) GUID {
  return raw::get_joystick_guid_for_id(instance_id)
}

fn get_joystick_vendor_for_id(instance_id JoystickID) u16 {
  return raw::get_joystick_vendor_for_id(instance_id)
}

fn get_joystick_product_for_id(instance_id JoystickID) u16 {
  return raw::get_joystick_product_for_id(instance_id)
}

fn get_joystick_product_version_for_id(instance_id JoystickID) u16 {
  return raw::get_joystick_product_version_for_id(instance_id)
}

fn get_joystick_type_for_id(instance_id JoystickID) JoystickType {
  return raw::get_joystick_type_for_id(instance_id)
}

fn open_joystick(instance_id JoystickID) *Joystick {
  return raw::open_joystick(instance_id)
}

fn get_joystick_from_id(instance_id JoystickID) *Joystick {
  return raw::get_joystick_from_id(instance_id)
}

fn get_joystick_from_player_index(player_index i32) *Joystick {
  return raw::get_joystick_from_player_index(player_index)
}

fn attach_virtual_joystick(desc const *VirtualJoystickDesc) JoystickID {
  return raw::attach_virtual_joystick(desc)
}

fn detach_virtual_joystick(instance_id JoystickID) bool {
  return raw::detach_virtual_joystick(instance_id)
}

fn is_joystick_virtual(instance_id JoystickID) bool {
  return raw::is_joystick_virtual(instance_id)
}

fn set_joystick_virtual_axis(joystick *Joystick, axis i32, value i16) bool {
  return raw::set_joystick_virtual_axis(joystick, axis, value)
}

fn set_joystick_virtual_ball(joystick *Joystick, ball i32, xrel i16, yrel i16) bool {
  return raw::set_joystick_virtual_ball(joystick, ball, xrel, yrel)
}

fn set_joystick_virtual_button(joystick *Joystick, button i32, down bool) bool {
  return raw::set_joystick_virtual_button(joystick, button, down)
}

fn set_joystick_virtual_hat(joystick *Joystick, hat i32, value u8) bool {
  return raw::set_joystick_virtual_hat(joystick, hat, value)
}

fn set_joystick_virtual_touchpad(joystick *Joystick, touchpad i32, finger i32, down bool, x f32, y f32, pressure f32) bool {
  return raw::set_joystick_virtual_touchpad(joystick, touchpad, finger, down, x, y, pressure)
}

fn get_joystick_properties(joystick *Joystick) PropertiesID {
  return raw::get_joystick_properties(joystick)
}

fn get_joystick_name(joystick *Joystick) OptionalString {
  result := raw::get_joystick_name(joystick)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_joystick_path(joystick *Joystick) OptionalString {
  result := raw::get_joystick_path(joystick)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_joystick_player_index(joystick *Joystick) i32 {
  return raw::get_joystick_player_index(joystick)
}

fn set_joystick_player_index(joystick *Joystick, player_index i32) bool {
  return raw::set_joystick_player_index(joystick, player_index)
}

fn get_joystick_guid(joystick *Joystick) GUID {
  return raw::get_joystick_guid(joystick)
}

fn get_joystick_vendor(joystick *Joystick) u16 {
  return raw::get_joystick_vendor(joystick)
}

fn get_joystick_product(joystick *Joystick) u16 {
  return raw::get_joystick_product(joystick)
}

fn get_joystick_product_version(joystick *Joystick) u16 {
  return raw::get_joystick_product_version(joystick)
}

fn get_joystick_firmware_version(joystick *Joystick) u16 {
  return raw::get_joystick_firmware_version(joystick)
}

fn get_joystick_serial(joystick *Joystick) OptionalString {
  result := raw::get_joystick_serial(joystick)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_joystick_type(joystick *Joystick) JoystickType {
  return raw::get_joystick_type(joystick)
}

fn get_joystick_guid_info(guid GUID, vendor *u16, product *u16, version *u16, crc16 *u16) {
  raw::get_joystick_guid_info(guid, vendor, product, version, crc16)
}

fn joystick_connected(joystick *Joystick) bool {
  return raw::joystick_connected(joystick)
}

fn get_joystick_id(joystick *Joystick) JoystickID {
  return raw::get_joystick_id(joystick)
}

fn get_num_joystick_axes(joystick *Joystick) i32 {
  return raw::get_num_joystick_axes(joystick)
}

fn get_num_joystick_balls(joystick *Joystick) i32 {
  return raw::get_num_joystick_balls(joystick)
}

fn get_num_joystick_hats(joystick *Joystick) i32 {
  return raw::get_num_joystick_hats(joystick)
}

fn get_num_joystick_buttons(joystick *Joystick) i32 {
  return raw::get_num_joystick_buttons(joystick)
}

fn set_joystick_events_enabled(enabled bool) {
  raw::set_joystick_events_enabled(enabled)
}

fn joystick_events_enabled() bool {
  return raw::joystick_events_enabled()
}

fn update_joysticks() {
  raw::update_joysticks()
}

fn get_joystick_axis(joystick *Joystick, axis i32) i16 {
  return raw::get_joystick_axis(joystick, axis)
}

fn get_joystick_axis_initial_state(joystick *Joystick, axis i32, state *i16) bool {
  return raw::get_joystick_axis_initial_state(joystick, axis, state)
}

fn get_joystick_hat(joystick *Joystick, hat i32) u8 {
  return raw::get_joystick_hat(joystick, hat)
}

fn get_joystick_button(joystick *Joystick, button i32) bool {
  return raw::get_joystick_button(joystick, button)
}

fn rumble_joystick(joystick *Joystick, low_frequency_rumble u16, high_frequency_rumble u16, duration_ms u32) bool {
  return raw::rumble_joystick(joystick, low_frequency_rumble, high_frequency_rumble, duration_ms)
}

fn rumble_joystick_triggers(joystick *Joystick, left_rumble u16, right_rumble u16, duration_ms u32) bool {
  return raw::rumble_joystick_triggers(joystick, left_rumble, right_rumble, duration_ms)
}

fn set_joystick_led(joystick *Joystick, red u8, green u8, blue u8) bool {
  return raw::set_joystick_led(joystick, red, green, blue)
}

fn send_joystick_effect(joystick *Joystick, data const *Void, size i32) bool {
  return raw::send_joystick_effect(joystick, data, size)
}

fn close_joystick(joystick *Joystick) {
  raw::close_joystick(joystick)
}

fn get_joystick_connection_state(joystick *Joystick) JoystickConnectionState {
  return raw::get_joystick_connection_state(joystick)
}

fn add_gamepad_mappings_from_io(src *IOStream, closeio bool) i32 {
  return raw::add_gamepad_mappings_from_io(src, closeio)
}

fn reload_gamepad_mappings() bool {
  return raw::reload_gamepad_mappings()
}

fn has_gamepad() bool {
  return raw::has_gamepad()
}

fn is_gamepad(instance_id JoystickID) bool {
  return raw::is_gamepad(instance_id)
}

fn get_gamepad_name_for_id(instance_id JoystickID) OptionalString {
  result := raw::get_gamepad_name_for_id(instance_id)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_gamepad_path_for_id(instance_id JoystickID) OptionalString {
  result := raw::get_gamepad_path_for_id(instance_id)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_gamepad_player_index_for_id(instance_id JoystickID) i32 {
  return raw::get_gamepad_player_index_for_id(instance_id)
}

fn get_gamepad_guid_for_id(instance_id JoystickID) GUID {
  return raw::get_gamepad_guid_for_id(instance_id)
}

fn get_gamepad_vendor_for_id(instance_id JoystickID) u16 {
  return raw::get_gamepad_vendor_for_id(instance_id)
}

fn get_gamepad_product_for_id(instance_id JoystickID) u16 {
  return raw::get_gamepad_product_for_id(instance_id)
}

fn get_gamepad_product_version_for_id(instance_id JoystickID) u16 {
  return raw::get_gamepad_product_version_for_id(instance_id)
}

fn get_gamepad_type_for_id(instance_id JoystickID) GamepadType {
  return raw::get_gamepad_type_for_id(instance_id)
}

fn get_real_gamepad_type_for_id(instance_id JoystickID) GamepadType {
  return raw::get_real_gamepad_type_for_id(instance_id)
}

fn open_gamepad(instance_id JoystickID) *Gamepad {
  return raw::open_gamepad(instance_id)
}

fn get_gamepad_from_id(instance_id JoystickID) *Gamepad {
  return raw::get_gamepad_from_id(instance_id)
}

fn get_gamepad_from_player_index(player_index i32) *Gamepad {
  return raw::get_gamepad_from_player_index(player_index)
}

fn get_gamepad_properties(gamepad *Gamepad) PropertiesID {
  return raw::get_gamepad_properties(gamepad)
}

fn get_gamepad_id(gamepad *Gamepad) JoystickID {
  return raw::get_gamepad_id(gamepad)
}

fn get_gamepad_name(gamepad *Gamepad) OptionalString {
  result := raw::get_gamepad_name(gamepad)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_gamepad_path(gamepad *Gamepad) OptionalString {
  result := raw::get_gamepad_path(gamepad)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_gamepad_type(gamepad *Gamepad) GamepadType {
  return raw::get_gamepad_type(gamepad)
}

fn get_real_gamepad_type(gamepad *Gamepad) GamepadType {
  return raw::get_real_gamepad_type(gamepad)
}

fn get_gamepad_player_index(gamepad *Gamepad) i32 {
  return raw::get_gamepad_player_index(gamepad)
}

fn set_gamepad_player_index(gamepad *Gamepad, player_index i32) bool {
  return raw::set_gamepad_player_index(gamepad, player_index)
}

fn get_gamepad_vendor(gamepad *Gamepad) u16 {
  return raw::get_gamepad_vendor(gamepad)
}

fn get_gamepad_product(gamepad *Gamepad) u16 {
  return raw::get_gamepad_product(gamepad)
}

fn get_gamepad_product_version(gamepad *Gamepad) u16 {
  return raw::get_gamepad_product_version(gamepad)
}

fn get_gamepad_firmware_version(gamepad *Gamepad) u16 {
  return raw::get_gamepad_firmware_version(gamepad)
}

fn get_gamepad_serial(gamepad *Gamepad) OptionalString {
  result := raw::get_gamepad_serial(gamepad)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_gamepad_steam_handle(gamepad *Gamepad) u64 {
  return raw::get_gamepad_steam_handle(gamepad)
}

fn get_gamepad_connection_state(gamepad *Gamepad) JoystickConnectionState {
  return raw::get_gamepad_connection_state(gamepad)
}

fn gamepad_connected(gamepad *Gamepad) bool {
  return raw::gamepad_connected(gamepad)
}

fn get_gamepad_joystick(gamepad *Gamepad) *Joystick {
  return raw::get_gamepad_joystick(gamepad)
}

fn set_gamepad_events_enabled(enabled bool) {
  raw::set_gamepad_events_enabled(enabled)
}

fn gamepad_events_enabled() bool {
  return raw::gamepad_events_enabled()
}

fn update_gamepads() {
  raw::update_gamepads()
}

fn get_gamepad_string_for_type(type_ GamepadType) OptionalString {
  result := raw::get_gamepad_string_for_type(type_)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_gamepad_string_for_axis(axis GamepadAxis) OptionalString {
  result := raw::get_gamepad_string_for_axis(axis)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn gamepad_has_axis(gamepad *Gamepad, axis GamepadAxis) bool {
  return raw::gamepad_has_axis(gamepad, axis)
}

fn get_gamepad_axis(gamepad *Gamepad, axis GamepadAxis) i16 {
  return raw::get_gamepad_axis(gamepad, axis)
}

fn get_gamepad_string_for_button(button GamepadButton) OptionalString {
  result := raw::get_gamepad_string_for_button(button)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn gamepad_has_button(gamepad *Gamepad, button GamepadButton) bool {
  return raw::gamepad_has_button(gamepad, button)
}

fn get_gamepad_button(gamepad *Gamepad, button GamepadButton) bool {
  return raw::get_gamepad_button(gamepad, button)
}

fn get_gamepad_button_label_for_type(type_ GamepadType, button GamepadButton) GamepadButtonLabel {
  return raw::get_gamepad_button_label_for_type(type_, button)
}

fn get_gamepad_button_label(gamepad *Gamepad, button GamepadButton) GamepadButtonLabel {
  return raw::get_gamepad_button_label(gamepad, button)
}

fn get_num_gamepad_touchpads(gamepad *Gamepad) i32 {
  return raw::get_num_gamepad_touchpads(gamepad)
}

fn get_num_gamepad_touchpad_fingers(gamepad *Gamepad, touchpad i32) i32 {
  return raw::get_num_gamepad_touchpad_fingers(gamepad, touchpad)
}

fn gamepad_has_sensor(gamepad *Gamepad, type_ SensorType) bool {
  return raw::gamepad_has_sensor(gamepad, type_)
}

fn set_gamepad_sensor_enabled(gamepad *Gamepad, type_ SensorType, enabled bool) bool {
  return raw::set_gamepad_sensor_enabled(gamepad, type_, enabled)
}

fn gamepad_sensor_enabled(gamepad *Gamepad, type_ SensorType) bool {
  return raw::gamepad_sensor_enabled(gamepad, type_)
}

fn get_gamepad_sensor_data_rate(gamepad *Gamepad, type_ SensorType) f32 {
  return raw::get_gamepad_sensor_data_rate(gamepad, type_)
}

fn rumble_gamepad(gamepad *Gamepad, low_frequency_rumble u16, high_frequency_rumble u16, duration_ms u32) bool {
  return raw::rumble_gamepad(gamepad, low_frequency_rumble, high_frequency_rumble, duration_ms)
}

fn rumble_gamepad_triggers(gamepad *Gamepad, left_rumble u16, right_rumble u16, duration_ms u32) bool {
  return raw::rumble_gamepad_triggers(gamepad, left_rumble, right_rumble, duration_ms)
}

fn set_gamepad_led(gamepad *Gamepad, red u8, green u8, blue u8) bool {
  return raw::set_gamepad_led(gamepad, red, green, blue)
}

fn send_gamepad_effect(gamepad *Gamepad, data const *Void, size i32) bool {
  return raw::send_gamepad_effect(gamepad, data, size)
}

fn close_gamepad(gamepad *Gamepad) {
  raw::close_gamepad(gamepad)
}

fn get_gamepad_apple_sf_symbols_name_for_button(gamepad *Gamepad, button GamepadButton) OptionalString {
  result := raw::get_gamepad_apple_sf_symbols_name_for_button(gamepad, button)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_gamepad_apple_sf_symbols_name_for_axis(gamepad *Gamepad, axis GamepadAxis) OptionalString {
  result := raw::get_gamepad_apple_sf_symbols_name_for_axis(gamepad, axis)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn has_keyboard() bool {
  return raw::has_keyboard()
}

fn get_keyboard_name_for_id(instance_id KeyboardID) OptionalString {
  result := raw::get_keyboard_name_for_id(instance_id)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_keyboard_focus() *Window {
  return raw::get_keyboard_focus()
}

fn reset_keyboard() {
  raw::reset_keyboard()
}

fn get_mod_state() Keymod {
  return raw::get_mod_state()
}

fn set_mod_state(modstate Keymod) {
  raw::set_mod_state(modstate)
}

fn get_key_from_scancode(scancode Scancode, modstate Keymod, key_event bool) Keycode {
  return raw::get_key_from_scancode(scancode, modstate, key_event)
}

fn get_scancode_from_key(key Keycode, modstate *Keymod) Scancode {
  return raw::get_scancode_from_key(key, modstate)
}

fn get_scancode_name(scancode Scancode) string {
  result := raw::get_scancode_name(scancode)
  if result == nil {
    return ""
  }
  return strings::from_c(result)
}

fn get_key_name(key Keycode) string {
  result := raw::get_key_name(key)
  if result == nil {
    return ""
  }
  return strings::from_c(result)
}

fn start_text_input(window *Window) bool {
  return raw::start_text_input(window)
}

fn start_text_input_with_properties(window *Window, props PropertiesID) bool {
  return raw::start_text_input_with_properties(window, props)
}

fn text_input_active(window *Window) bool {
  return raw::text_input_active(window)
}

fn stop_text_input(window *Window) bool {
  return raw::stop_text_input(window)
}

fn clear_composition(window *Window) bool {
  return raw::clear_composition(window)
}

fn set_text_input_area(window *Window, rect const *Rect, cursor i32) bool {
  return raw::set_text_input_area(window, rect, cursor)
}

fn has_screen_keyboard_support() bool {
  return raw::has_screen_keyboard_support()
}

fn screen_keyboard_shown(window *Window) bool {
  return raw::screen_keyboard_shown(window)
}

fn has_mouse() bool {
  return raw::has_mouse()
}

fn get_mouse_name_for_id(instance_id MouseID) OptionalString {
  result := raw::get_mouse_name_for_id(instance_id)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_mouse_focus() *Window {
  return raw::get_mouse_focus()
}

fn warp_mouse_in_window(window *Window, x f32, y f32) {
  raw::warp_mouse_in_window(window, x, y)
}

fn warp_mouse_global(x f32, y f32) bool {
  return raw::warp_mouse_global(x, y)
}

fn set_window_relative_mouse_mode(window *Window, enabled bool) bool {
  return raw::set_window_relative_mouse_mode(window, enabled)
}

fn get_window_relative_mouse_mode(window *Window) bool {
  return raw::get_window_relative_mouse_mode(window)
}

fn capture_mouse(enabled bool) bool {
  return raw::capture_mouse(enabled)
}

fn create_cursor(data const *u8, mask const *u8, w i32, h i32, hot_x i32, hot_y i32) *Cursor {
  return raw::create_cursor(data, mask, w, h, hot_x, hot_y)
}

fn create_color_cursor(surface *Surface, hot_x i32, hot_y i32) *Cursor {
  return raw::create_color_cursor(surface, hot_x, hot_y)
}

fn create_system_cursor(id SystemCursor) *Cursor {
  return raw::create_system_cursor(id)
}

fn set_cursor(cursor *Cursor) bool {
  return raw::set_cursor(cursor)
}

fn get_cursor() *Cursor {
  return raw::get_cursor()
}

fn get_default_cursor() *Cursor {
  return raw::get_default_cursor()
}

fn destroy_cursor(cursor *Cursor) {
  raw::destroy_cursor(cursor)
}

fn show_cursor() bool {
  return raw::show_cursor()
}

fn hide_cursor() bool {
  return raw::hide_cursor()
}

fn cursor_visible() bool {
  return raw::cursor_visible()
}

fn get_touch_device_name(touchID TouchID) OptionalString {
  result := raw::get_touch_device_name(touchID)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_touch_device_type(touchID TouchID) TouchDeviceType {
  return raw::get_touch_device_type(touchID)
}

fn pump_events() {
  raw::pump_events()
}

fn peep_events(events *Event, numevents i32, action EventAction, minType u32, maxType u32) i32 {
  return raw::peep_events(events, numevents, action, minType, maxType)
}

fn has_event(type_ u32) bool {
  return raw::has_event(type_)
}

fn has_events(minType u32, maxType u32) bool {
  return raw::has_events(minType, maxType)
}

fn flush_event(type_ u32) {
  raw::flush_event(type_)
}

fn flush_events(minType u32, maxType u32) {
  raw::flush_events(minType, maxType)
}

fn poll_event(event *Event) bool {
  return raw::poll_event(event)
}

fn wait_event(event *Event) bool {
  return raw::wait_event(event)
}

fn wait_event_timeout(event *Event, timeoutMS i32) bool {
  return raw::wait_event_timeout(event, timeoutMS)
}

fn push_event(event *Event) bool {
  return raw::push_event(event)
}

fn set_event_enabled(type_ u32, enabled bool) {
  raw::set_event_enabled(type_, enabled)
}

fn event_enabled(type_ u32) bool {
  return raw::event_enabled(type_)
}

fn register_events(numevents i32) u32 {
  return raw::register_events(numevents)
}

fn get_window_from_event(event const *Event) *Window {
  return raw::get_window_from_event(event)
}

fn get_base_path() OptionalString {
  result := raw::get_base_path()
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_user_folder(folder Folder) OptionalString {
  result := raw::get_user_folder(folder)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn gpu_supports_properties(props PropertiesID) bool {
  return raw::gpu_supports_properties(props)
}

fn create_gpu_device_with_properties(props PropertiesID) *GPUDevice {
  return raw::create_gpu_device_with_properties(props)
}

fn destroy_gpu_device(device *GPUDevice) {
  raw::destroy_gpu_device(device)
}

fn get_num_gpu_drivers() i32 {
  return raw::get_num_gpu_drivers()
}

fn get_gpu_driver(index i32) OptionalString {
  result := raw::get_gpu_driver(index)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_gpu_device_driver(device *GPUDevice) OptionalString {
  result := raw::get_gpu_device_driver(device)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_gpu_shader_formats(device *GPUDevice) GPUShaderFormat {
  return raw::get_gpu_shader_formats(device)
}

fn create_gpu_compute_pipeline(device *GPUDevice, createinfo const *GPUComputePipelineCreateInfo) *GPUComputePipeline {
  return raw::create_gpu_compute_pipeline(device, createinfo)
}

fn create_gpu_graphics_pipeline(device *GPUDevice, createinfo const *GPUGraphicsPipelineCreateInfo) *GPUGraphicsPipeline {
  return raw::create_gpu_graphics_pipeline(device, createinfo)
}

fn create_gpu_sampler(device *GPUDevice, createinfo const *GPUSamplerCreateInfo) *GPUSampler {
  return raw::create_gpu_sampler(device, createinfo)
}

fn create_gpu_shader(device *GPUDevice, createinfo const *GPUShaderCreateInfo) *GPUShader {
  return raw::create_gpu_shader(device, createinfo)
}

fn create_gpu_texture(device *GPUDevice, createinfo const *GPUTextureCreateInfo) *GPUTexture {
  return raw::create_gpu_texture(device, createinfo)
}

fn create_gpu_buffer(device *GPUDevice, createinfo const *GPUBufferCreateInfo) *GPUBuffer {
  return raw::create_gpu_buffer(device, createinfo)
}

fn create_gpu_transfer_buffer(device *GPUDevice, createinfo const *GPUTransferBufferCreateInfo) *GPUTransferBuffer {
  return raw::create_gpu_transfer_buffer(device, createinfo)
}

fn pop_gpu_debug_group(command_buffer *GPUCommandBuffer) {
  raw::pop_gpu_debug_group(command_buffer)
}

fn release_gpu_texture(device *GPUDevice, texture *GPUTexture) {
  raw::release_gpu_texture(device, texture)
}

fn release_gpu_sampler(device *GPUDevice, sampler *GPUSampler) {
  raw::release_gpu_sampler(device, sampler)
}

fn release_gpu_buffer(device *GPUDevice, buffer *GPUBuffer) {
  raw::release_gpu_buffer(device, buffer)
}

fn release_gpu_transfer_buffer(device *GPUDevice, transfer_buffer *GPUTransferBuffer) {
  raw::release_gpu_transfer_buffer(device, transfer_buffer)
}

fn release_gpu_compute_pipeline(device *GPUDevice, compute_pipeline *GPUComputePipeline) {
  raw::release_gpu_compute_pipeline(device, compute_pipeline)
}

fn release_gpu_shader(device *GPUDevice, shader *GPUShader) {
  raw::release_gpu_shader(device, shader)
}

fn release_gpu_graphics_pipeline(device *GPUDevice, graphics_pipeline *GPUGraphicsPipeline) {
  raw::release_gpu_graphics_pipeline(device, graphics_pipeline)
}

fn acquire_gpu_command_buffer(device *GPUDevice) *GPUCommandBuffer {
  return raw::acquire_gpu_command_buffer(device)
}

fn push_gpu_vertex_uniform_data(command_buffer *GPUCommandBuffer, slot_index u32, data const *Void, length u32) {
  raw::push_gpu_vertex_uniform_data(command_buffer, slot_index, data, length)
}

fn push_gpu_fragment_uniform_data(command_buffer *GPUCommandBuffer, slot_index u32, data const *Void, length u32) {
  raw::push_gpu_fragment_uniform_data(command_buffer, slot_index, data, length)
}

fn push_gpu_compute_uniform_data(command_buffer *GPUCommandBuffer, slot_index u32, data const *Void, length u32) {
  raw::push_gpu_compute_uniform_data(command_buffer, slot_index, data, length)
}

fn begin_gpu_render_pass(command_buffer *GPUCommandBuffer, color_target_infos const *GPUColorTargetInfo, num_color_targets u32, depth_stencil_target_info const *GPUDepthStencilTargetInfo) *GPURenderPass {
  return raw::begin_gpu_render_pass(command_buffer, color_target_infos, num_color_targets, depth_stencil_target_info)
}

fn bind_gpu_graphics_pipeline(render_pass *GPURenderPass, graphics_pipeline *GPUGraphicsPipeline) {
  raw::bind_gpu_graphics_pipeline(render_pass, graphics_pipeline)
}

fn set_gpu_viewport(render_pass *GPURenderPass, viewport const *GPUViewport) {
  raw::set_gpu_viewport(render_pass, viewport)
}

fn set_gpu_scissor(render_pass *GPURenderPass, scissor const *Rect) {
  raw::set_gpu_scissor(render_pass, scissor)
}

fn set_gpu_blend_constants(render_pass *GPURenderPass, blend_constants FColor) {
  raw::set_gpu_blend_constants(render_pass, blend_constants)
}

fn set_gpu_stencil_reference(render_pass *GPURenderPass, reference u8) {
  raw::set_gpu_stencil_reference(render_pass, reference)
}

fn bind_gpu_vertex_buffers(render_pass *GPURenderPass, first_slot u32, bindings const *GPUBufferBinding, num_bindings u32) {
  raw::bind_gpu_vertex_buffers(render_pass, first_slot, bindings, num_bindings)
}

fn bind_gpu_index_buffer(render_pass *GPURenderPass, binding const *GPUBufferBinding, index_element_size GPUIndexElementSize) {
  raw::bind_gpu_index_buffer(render_pass, binding, index_element_size)
}

fn bind_gpu_vertex_samplers(render_pass *GPURenderPass, first_slot u32, texture_sampler_bindings const *GPUTextureSamplerBinding, num_bindings u32) {
  raw::bind_gpu_vertex_samplers(render_pass, first_slot, texture_sampler_bindings, num_bindings)
}

fn bind_gpu_vertex_storage_textures(render_pass *GPURenderPass, first_slot u32, storage_textures const **GPUTexture, num_bindings u32) {
  raw::bind_gpu_vertex_storage_textures(render_pass, first_slot, storage_textures, num_bindings)
}

fn bind_gpu_vertex_storage_buffers(render_pass *GPURenderPass, first_slot u32, storage_buffers const **GPUBuffer, num_bindings u32) {
  raw::bind_gpu_vertex_storage_buffers(render_pass, first_slot, storage_buffers, num_bindings)
}

fn bind_gpu_fragment_samplers(render_pass *GPURenderPass, first_slot u32, texture_sampler_bindings const *GPUTextureSamplerBinding, num_bindings u32) {
  raw::bind_gpu_fragment_samplers(render_pass, first_slot, texture_sampler_bindings, num_bindings)
}

fn bind_gpu_fragment_storage_textures(render_pass *GPURenderPass, first_slot u32, storage_textures const **GPUTexture, num_bindings u32) {
  raw::bind_gpu_fragment_storage_textures(render_pass, first_slot, storage_textures, num_bindings)
}

fn bind_gpu_fragment_storage_buffers(render_pass *GPURenderPass, first_slot u32, storage_buffers const **GPUBuffer, num_bindings u32) {
  raw::bind_gpu_fragment_storage_buffers(render_pass, first_slot, storage_buffers, num_bindings)
}

fn draw_gpu_indexed_primitives(render_pass *GPURenderPass, num_indices u32, num_instances u32, first_index u32, vertex_offset i32, first_instance u32) {
  raw::draw_gpu_indexed_primitives(render_pass, num_indices, num_instances, first_index, vertex_offset, first_instance)
}

fn draw_gpu_primitives(render_pass *GPURenderPass, num_vertices u32, num_instances u32, first_vertex u32, first_instance u32) {
  raw::draw_gpu_primitives(render_pass, num_vertices, num_instances, first_vertex, first_instance)
}

fn draw_gpu_primitives_indirect(render_pass *GPURenderPass, buffer *GPUBuffer, offset u32, draw_count u32) {
  raw::draw_gpu_primitives_indirect(render_pass, buffer, offset, draw_count)
}

fn draw_gpu_indexed_primitives_indirect(render_pass *GPURenderPass, buffer *GPUBuffer, offset u32, draw_count u32) {
  raw::draw_gpu_indexed_primitives_indirect(render_pass, buffer, offset, draw_count)
}

fn end_gpu_render_pass(render_pass *GPURenderPass) {
  raw::end_gpu_render_pass(render_pass)
}

fn begin_gpu_compute_pass(command_buffer *GPUCommandBuffer, storage_texture_bindings const *GPUStorageTextureReadWriteBinding, num_storage_texture_bindings u32, storage_buffer_bindings const *GPUStorageBufferReadWriteBinding, num_storage_buffer_bindings u32) *GPUComputePass {
  return raw::begin_gpu_compute_pass(command_buffer, storage_texture_bindings, num_storage_texture_bindings, storage_buffer_bindings, num_storage_buffer_bindings)
}

fn bind_gpu_compute_pipeline(compute_pass *GPUComputePass, compute_pipeline *GPUComputePipeline) {
  raw::bind_gpu_compute_pipeline(compute_pass, compute_pipeline)
}

fn bind_gpu_compute_samplers(compute_pass *GPUComputePass, first_slot u32, texture_sampler_bindings const *GPUTextureSamplerBinding, num_bindings u32) {
  raw::bind_gpu_compute_samplers(compute_pass, first_slot, texture_sampler_bindings, num_bindings)
}

fn bind_gpu_compute_storage_textures(compute_pass *GPUComputePass, first_slot u32, storage_textures const **GPUTexture, num_bindings u32) {
  raw::bind_gpu_compute_storage_textures(compute_pass, first_slot, storage_textures, num_bindings)
}

fn bind_gpu_compute_storage_buffers(compute_pass *GPUComputePass, first_slot u32, storage_buffers const **GPUBuffer, num_bindings u32) {
  raw::bind_gpu_compute_storage_buffers(compute_pass, first_slot, storage_buffers, num_bindings)
}

fn dispatch_gpu_compute(compute_pass *GPUComputePass, groupcount_x u32, groupcount_y u32, groupcount_z u32) {
  raw::dispatch_gpu_compute(compute_pass, groupcount_x, groupcount_y, groupcount_z)
}

fn dispatch_gpu_compute_indirect(compute_pass *GPUComputePass, buffer *GPUBuffer, offset u32) {
  raw::dispatch_gpu_compute_indirect(compute_pass, buffer, offset)
}

fn end_gpu_compute_pass(compute_pass *GPUComputePass) {
  raw::end_gpu_compute_pass(compute_pass)
}

fn map_gpu_transfer_buffer(device *GPUDevice, transfer_buffer *GPUTransferBuffer, cycle bool) rawptr {
  return raw::map_gpu_transfer_buffer(device, transfer_buffer, cycle)
}

fn unmap_gpu_transfer_buffer(device *GPUDevice, transfer_buffer *GPUTransferBuffer) {
  raw::unmap_gpu_transfer_buffer(device, transfer_buffer)
}

fn begin_gpu_copy_pass(command_buffer *GPUCommandBuffer) *GPUCopyPass {
  return raw::begin_gpu_copy_pass(command_buffer)
}

fn upload_to_gpu_texture(copy_pass *GPUCopyPass, source const *GPUTextureTransferInfo, destination const *GPUTextureRegion, cycle bool) {
  raw::upload_to_gpu_texture(copy_pass, source, destination, cycle)
}

fn upload_to_gpu_buffer(copy_pass *GPUCopyPass, source const *GPUTransferBufferLocation, destination const *GPUBufferRegion, cycle bool) {
  raw::upload_to_gpu_buffer(copy_pass, source, destination, cycle)
}

fn copy_gpu_texture_to_texture(copy_pass *GPUCopyPass, source const *GPUTextureLocation, destination const *GPUTextureLocation, w u32, h u32, d u32, cycle bool) {
  raw::copy_gpu_texture_to_texture(copy_pass, source, destination, w, h, d, cycle)
}

fn copy_gpu_buffer_to_buffer(copy_pass *GPUCopyPass, source const *GPUBufferLocation, destination const *GPUBufferLocation, size u32, cycle bool) {
  raw::copy_gpu_buffer_to_buffer(copy_pass, source, destination, size, cycle)
}

fn download_from_gpu_texture(copy_pass *GPUCopyPass, source const *GPUTextureRegion, destination const *GPUTextureTransferInfo) {
  raw::download_from_gpu_texture(copy_pass, source, destination)
}

fn download_from_gpu_buffer(copy_pass *GPUCopyPass, source const *GPUBufferRegion, destination const *GPUTransferBufferLocation) {
  raw::download_from_gpu_buffer(copy_pass, source, destination)
}

fn end_gpu_copy_pass(copy_pass *GPUCopyPass) {
  raw::end_gpu_copy_pass(copy_pass)
}

fn generate_mipmaps_for_gpu_texture(command_buffer *GPUCommandBuffer, texture *GPUTexture) {
  raw::generate_mipmaps_for_gpu_texture(command_buffer, texture)
}

fn blit_gpu_texture(command_buffer *GPUCommandBuffer, info const *GPUBlitInfo) {
  raw::blit_gpu_texture(command_buffer, info)
}

fn window_supports_gpu_swapchain_composition(device *GPUDevice, window *Window, swapchain_composition GPUSwapchainComposition) bool {
  return raw::window_supports_gpu_swapchain_composition(device, window, swapchain_composition)
}

fn window_supports_gpu_present_mode(device *GPUDevice, window *Window, present_mode GPUPresentMode) bool {
  return raw::window_supports_gpu_present_mode(device, window, present_mode)
}

fn claim_window_for_gpu_device(device *GPUDevice, window *Window) bool {
  return raw::claim_window_for_gpu_device(device, window)
}

fn release_window_from_gpu_device(device *GPUDevice, window *Window) {
  raw::release_window_from_gpu_device(device, window)
}

fn set_gpu_swapchain_parameters(device *GPUDevice, window *Window, swapchain_composition GPUSwapchainComposition, present_mode GPUPresentMode) bool {
  return raw::set_gpu_swapchain_parameters(device, window, swapchain_composition, present_mode)
}

fn set_gpu_allowed_frames_in_flight(device *GPUDevice, allowed_frames_in_flight u32) bool {
  return raw::set_gpu_allowed_frames_in_flight(device, allowed_frames_in_flight)
}

fn get_gpu_swapchain_texture_format(device *GPUDevice, window *Window) GPUTextureFormat {
  return raw::get_gpu_swapchain_texture_format(device, window)
}

fn acquire_gpu_swapchain_texture(command_buffer *GPUCommandBuffer, window *Window, swapchain_texture **GPUTexture, swapchain_texture_width *u32, swapchain_texture_height *u32) bool {
  return raw::acquire_gpu_swapchain_texture(command_buffer, window, swapchain_texture, swapchain_texture_width, swapchain_texture_height)
}

fn wait_for_gpu_swapchain(device *GPUDevice, window *Window) bool {
  return raw::wait_for_gpu_swapchain(device, window)
}

fn wait_and_acquire_gpu_swapchain_texture(command_buffer *GPUCommandBuffer, window *Window, swapchain_texture **GPUTexture, swapchain_texture_width *u32, swapchain_texture_height *u32) bool {
  return raw::wait_and_acquire_gpu_swapchain_texture(command_buffer, window, swapchain_texture, swapchain_texture_width, swapchain_texture_height)
}

fn submit_gpu_command_buffer(command_buffer *GPUCommandBuffer) bool {
  return raw::submit_gpu_command_buffer(command_buffer)
}

fn submit_gpu_command_buffer_and_acquire_fence(command_buffer *GPUCommandBuffer) *GPUFence {
  return raw::submit_gpu_command_buffer_and_acquire_fence(command_buffer)
}

fn cancel_gpu_command_buffer(command_buffer *GPUCommandBuffer) bool {
  return raw::cancel_gpu_command_buffer(command_buffer)
}

fn wait_for_gpu_idle(device *GPUDevice) bool {
  return raw::wait_for_gpu_idle(device)
}

fn wait_for_gpu_fences(device *GPUDevice, wait_all bool, fences const **GPUFence, num_fences u32) bool {
  return raw::wait_for_gpu_fences(device, wait_all, fences, num_fences)
}

fn query_gpu_fence(device *GPUDevice, fence *GPUFence) bool {
  return raw::query_gpu_fence(device, fence)
}

fn release_gpu_fence(device *GPUDevice, fence *GPUFence) {
  raw::release_gpu_fence(device, fence)
}

fn gpu_texture_format_texel_block_size(format GPUTextureFormat) u32 {
  return raw::gpu_texture_format_texel_block_size(format)
}

fn gpu_texture_supports_format(device *GPUDevice, format GPUTextureFormat, type_ GPUTextureType, usage GPUTextureUsageFlags) bool {
  return raw::gpu_texture_supports_format(device, format, type_, usage)
}

fn gpu_texture_supports_sample_count(device *GPUDevice, format GPUTextureFormat, sample_count GPUSampleCount) bool {
  return raw::gpu_texture_supports_sample_count(device, format, sample_count)
}

fn calculate_gpu_texture_format_size(format GPUTextureFormat, width u32, height u32, depth_or_layer_count u32) u32 {
  return raw::calculate_gpu_texture_format_size(format, width, height, depth_or_layer_count)
}

fn get_haptic_name_for_id(instance_id HapticID) OptionalString {
  result := raw::get_haptic_name_for_id(instance_id)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn open_haptic(instance_id HapticID) *Haptic {
  return raw::open_haptic(instance_id)
}

fn get_haptic_from_id(instance_id HapticID) *Haptic {
  return raw::get_haptic_from_id(instance_id)
}

fn get_haptic_id(haptic *Haptic) HapticID {
  return raw::get_haptic_id(haptic)
}

fn get_haptic_name(haptic *Haptic) OptionalString {
  result := raw::get_haptic_name(haptic)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn is_mouse_haptic() bool {
  return raw::is_mouse_haptic()
}

fn open_haptic_from_mouse() *Haptic {
  return raw::open_haptic_from_mouse()
}

fn is_joystick_haptic(joystick *Joystick) bool {
  return raw::is_joystick_haptic(joystick)
}

fn open_haptic_from_joystick(joystick *Joystick) *Haptic {
  return raw::open_haptic_from_joystick(joystick)
}

fn close_haptic(haptic *Haptic) {
  raw::close_haptic(haptic)
}

fn get_max_haptic_effects(haptic *Haptic) i32 {
  return raw::get_max_haptic_effects(haptic)
}

fn get_max_haptic_effects_playing(haptic *Haptic) i32 {
  return raw::get_max_haptic_effects_playing(haptic)
}

fn get_haptic_features(haptic *Haptic) u32 {
  return raw::get_haptic_features(haptic)
}

fn get_num_haptic_axes(haptic *Haptic) i32 {
  return raw::get_num_haptic_axes(haptic)
}

fn haptic_effect_supported(haptic *Haptic, effect const *HapticEffect) bool {
  return raw::haptic_effect_supported(haptic, effect)
}

fn create_haptic_effect(haptic *Haptic, effect const *HapticEffect) i32 {
  return raw::create_haptic_effect(haptic, effect)
}

fn update_haptic_effect(haptic *Haptic, effect i32, data const *HapticEffect) bool {
  return raw::update_haptic_effect(haptic, effect, data)
}

fn run_haptic_effect(haptic *Haptic, effect i32, iterations u32) bool {
  return raw::run_haptic_effect(haptic, effect, iterations)
}

fn stop_haptic_effect(haptic *Haptic, effect i32) bool {
  return raw::stop_haptic_effect(haptic, effect)
}

fn destroy_haptic_effect(haptic *Haptic, effect i32) {
  raw::destroy_haptic_effect(haptic, effect)
}

fn get_haptic_effect_status(haptic *Haptic, effect i32) bool {
  return raw::get_haptic_effect_status(haptic, effect)
}

fn set_haptic_gain(haptic *Haptic, gain i32) bool {
  return raw::set_haptic_gain(haptic, gain)
}

fn set_haptic_autocenter(haptic *Haptic, autocenter i32) bool {
  return raw::set_haptic_autocenter(haptic, autocenter)
}

fn pause_haptic(haptic *Haptic) bool {
  return raw::pause_haptic(haptic)
}

fn resume_haptic(haptic *Haptic) bool {
  return raw::resume_haptic(haptic)
}

fn stop_haptic_effects(haptic *Haptic) bool {
  return raw::stop_haptic_effects(haptic)
}

fn haptic_rumble_supported(haptic *Haptic) bool {
  return raw::haptic_rumble_supported(haptic)
}

fn init_haptic_rumble(haptic *Haptic) bool {
  return raw::init_haptic_rumble(haptic)
}

fn play_haptic_rumble(haptic *Haptic, strength f32, length u32) bool {
  return raw::play_haptic_rumble(haptic, strength, length)
}

fn stop_haptic_rumble(haptic *Haptic) bool {
  return raw::stop_haptic_rumble(haptic)
}

fn hid_init() i32 {
  return raw::hid_init()
}

fn hid_exit() i32 {
  return raw::hid_exit()
}

fn hid_device_change_count() u32 {
  return raw::hid_device_change_count()
}

fn hid_enumerate(vendor_id u16, product_id u16) *hid_device_info {
  return raw::hid_enumerate(vendor_id, product_id)
}

fn hid_free_enumeration(devs *hid_device_info) {
  raw::hid_free_enumeration(devs)
}

fn hid_set_nonblocking(dev *hid_device, nonblock i32) i32 {
  return raw::hid_set_nonblocking(dev, nonblock)
}

fn hid_close(dev *hid_device) i32 {
  return raw::hid_close(dev)
}

fn hid_get_device_info(dev *hid_device) *hid_device_info {
  return raw::hid_get_device_info(dev)
}

fn hid_ble_scan(active bool) {
  raw::hid_ble_scan(active)
}

fn set_hint_with_priority(name string, value string, priority HintPriority) error {
  name_c, name_err := strings::to_c(name)
  if name_err != nil {
    return name_err
  }
  defer strings::free_c(&name_c)
  value_c, value_err := strings::to_c(value)
  if value_err != nil {
    return value_err
  }
  defer strings::free_c(&value_c)
  if raw::set_hint_with_priority(strings::c_view(&name_c), strings::c_view(&value_c), priority) == false {
    return error("SDL set_hint_with_priority failed")
  }
  return nil
}

fn set_hint(name string, value string) error {
  name_c, name_err := strings::to_c(name)
  if name_err != nil {
    return name_err
  }
  defer strings::free_c(&name_c)
  value_c, value_err := strings::to_c(value)
  if value_err != nil {
    return value_err
  }
  defer strings::free_c(&value_c)
  if raw::set_hint(strings::c_view(&name_c), strings::c_view(&value_c)) == false {
    return error("SDL set_hint failed")
  }
  return nil
}

fn reset_hint(name string) error {
  name_c, name_err := strings::to_c(name)
  if name_err != nil {
    return name_err
  }
  defer strings::free_c(&name_c)
  if raw::reset_hint(strings::c_view(&name_c)) == false {
    return error("SDL reset_hint failed")
  }
  return nil
}

fn reset_hints() {
  raw::reset_hints()
}

fn get_hint(name string) (OptionalString, error) {
  name_c, name_err := strings::to_c(name)
  if name_err != nil {
    return OptionalString{value: "", present: false}, name_err
  }
  defer strings::free_c(&name_c)
  result := raw::get_hint(strings::c_view(&name_c))
  if result == nil {
    return OptionalString{value: "", present: false}, nil
  }
  return OptionalString{value: strings::from_c(result), present: true}, nil
}

fn init(flags InitFlags) bool {
  return raw::init(flags)
}

fn init_sub_system(flags InitFlags) bool {
  return raw::init_sub_system(flags)
}

fn quit_sub_system(flags InitFlags) {
  raw::quit_sub_system(flags)
}

fn was_init(flags InitFlags) InitFlags {
  return raw::was_init(flags)
}

fn quit() {
  raw::quit()
}

fn is_main_thread() bool {
  return raw::is_main_thread()
}

fn unload_object(handle *SharedObject) {
  raw::unload_object(handle)
}

fn set_log_priorities(priority LogPriority) {
  raw::set_log_priorities(priority)
}

fn set_log_priority(category i32, priority LogPriority) {
  raw::set_log_priority(category, priority)
}

fn get_log_priority(category i32) LogPriority {
  return raw::get_log_priority(category)
}

fn reset_log_priorities() {
  raw::reset_log_priorities()
}

fn get_default_log_output_function() LogOutputFunction {
  return raw::get_default_log_output_function()
}

fn metal_create_view(window *Window) MetalView {
  return raw::metal_create_view(window)
}

fn metal_destroy_view(view MetalView) {
  raw::metal_destroy_view(view)
}

fn metal_get_layer(view MetalView) rawptr {
  return raw::metal_get_layer(view)
}

fn open_url(url string) error {
  url_c, url_err := strings::to_c(url)
  if url_err != nil {
    return url_err
  }
  defer strings::free_c(&url_c)
  if raw::open_url(strings::c_view(&url_c)) == false {
    return error("SDL open_url failed")
  }
  return nil
}

fn get_platform() string {
  result := raw::get_platform()
  if result == nil {
    return ""
  }
  return strings::from_c(result)
}

fn create_process_with_properties(props PropertiesID) *Process {
  return raw::create_process_with_properties(props)
}

fn get_process_properties(process *Process) PropertiesID {
  return raw::get_process_properties(process)
}

fn get_process_input(process *Process) *IOStream {
  return raw::get_process_input(process)
}

fn get_process_output(process *Process) *IOStream {
  return raw::get_process_output(process)
}

fn kill_process(process *Process, force bool) bool {
  return raw::kill_process(process, force)
}

fn destroy_process(process *Process) {
  raw::destroy_process(process)
}

fn get_num_render_drivers() i32 {
  return raw::get_num_render_drivers()
}

fn get_render_driver(index i32) OptionalString {
  result := raw::get_render_driver(index)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn create_renderer_with_properties(props PropertiesID) *Renderer {
  return raw::create_renderer_with_properties(props)
}

fn create_software_renderer(surface *Surface) *Renderer {
  return raw::create_software_renderer(surface)
}

fn get_renderer(window *Window) *Renderer {
  return raw::get_renderer(window)
}

fn get_render_window(renderer *Renderer) *Window {
  return raw::get_render_window(renderer)
}

fn get_renderer_name(renderer *Renderer) OptionalString {
  result := raw::get_renderer_name(renderer)
  if result == nil {
    return OptionalString{value: "", present: false}
  }
  return OptionalString{value: strings::from_c(result), present: true}
}

fn get_renderer_properties(renderer *Renderer) PropertiesID {
  return raw::get_renderer_properties(renderer)
}

fn create_texture(renderer *Renderer, format PixelFormat, access TextureAccess, w i32, h i32) *Texture {
  return raw::create_texture(renderer, format, access, w, h)
}

fn create_texture_from_surface(renderer *Renderer, surface *Surface) *Texture {
  return raw::create_texture_from_surface(renderer, surface)
}

fn create_texture_with_properties(renderer *Renderer, props PropertiesID) *Texture {
  return raw::create_texture_with_properties(renderer, props)
}

fn get_texture_properties(texture *Texture) PropertiesID {
  return raw::get_texture_properties(texture)
}

fn get_renderer_from_texture(texture *Texture) *Renderer {
  return raw::get_renderer_from_texture(texture)
}

fn set_texture_color_mod(texture *Texture, r u8, g u8, b u8) bool {
  return raw::set_texture_color_mod(texture, r, g, b)
}

fn set_texture_color_mod_float(texture *Texture, r f32, g f32, b f32) bool {
  return raw::set_texture_color_mod_float(texture, r, g, b)
}

fn get_texture_color_mod(texture *Texture, r *u8, g *u8, b *u8) bool {
  return raw::get_texture_color_mod(texture, r, g, b)
}

fn set_texture_alpha_mod(texture *Texture, alpha u8) bool {
  return raw::set_texture_alpha_mod(texture, alpha)
}

fn set_texture_alpha_mod_float(texture *Texture, alpha f32) bool {
  return raw::set_texture_alpha_mod_float(texture, alpha)
}

fn get_texture_alpha_mod(texture *Texture, alpha *u8) bool {
  return raw::get_texture_alpha_mod(texture, alpha)
}

fn set_texture_blend_mode(texture *Texture, blendMode BlendMode) bool {
  return raw::set_texture_blend_mode(texture, blendMode)
}

fn get_texture_blend_mode(texture *Texture, blendMode *BlendMode) bool {
  return raw::get_texture_blend_mode(texture, blendMode)
}

fn set_texture_scale_mode(texture *Texture, scaleMode ScaleMode) bool {
  return raw::set_texture_scale_mode(texture, scaleMode)
}

fn get_texture_scale_mode(texture *Texture, scaleMode *ScaleMode) bool {
  return raw::get_texture_scale_mode(texture, scaleMode)
}

fn update_texture(texture *Texture, rect const *Rect, pixels const *Void, pitch i32) bool {
  return raw::update_texture(texture, rect, pixels, pitch)
}

fn update_yuv_texture(texture *Texture, rect const *Rect, Yplane const *u8, Ypitch i32, Uplane const *u8, Upitch i32, Vplane const *u8, Vpitch i32) bool {
  return raw::update_yuv_texture(texture, rect, Yplane, Ypitch, Uplane, Upitch, Vplane, Vpitch)
}

fn update_nv_texture(texture *Texture, rect const *Rect, Yplane const *u8, Ypitch i32, UVplane const *u8, UVpitch i32) bool {
  return raw::update_nv_texture(texture, rect, Yplane, Ypitch, UVplane, UVpitch)
}

fn lock_texture_to_surface(texture *Texture, rect const *Rect, surface **Surface) bool {
  return raw::lock_texture_to_surface(texture, rect, surface)
}

fn unlock_texture(texture *Texture) {
  raw::unlock_texture(texture)
}

fn set_render_target(renderer *Renderer, texture *Texture) bool {
  return raw::set_render_target(renderer, texture)
}

fn get_render_target(renderer *Renderer) *Texture {
  return raw::get_render_target(renderer)
}

fn set_render_logical_presentation(renderer *Renderer, w i32, h i32, mode RendererLogicalPresentation) bool {
  return raw::set_render_logical_presentation(renderer, w, h, mode)
}

fn get_render_logical_presentation_rect(renderer *Renderer, rect *FRect) bool {
  return raw::get_render_logical_presentation_rect(renderer, rect)
}

fn convert_event_to_render_coordinates(renderer *Renderer, event *Event) bool {
  return raw::convert_event_to_render_coordinates(renderer, event)
}

fn set_render_viewport(renderer *Renderer, rect const *Rect) bool {
  return raw::set_render_viewport(renderer, rect)
}

fn get_render_viewport(renderer *Renderer, rect *Rect) bool {
  return raw::get_render_viewport(renderer, rect)
}

fn render_viewport_set(renderer *Renderer) bool {
  return raw::render_viewport_set(renderer)
}

fn get_render_safe_area(renderer *Renderer, rect *Rect) bool {
  return raw::get_render_safe_area(renderer, rect)
}

fn set_render_clip_rect(renderer *Renderer, rect const *Rect) bool {
  return raw::set_render_clip_rect(renderer, rect)
}

fn get_render_clip_rect(renderer *Renderer, rect *Rect) bool {
  return raw::get_render_clip_rect(renderer, rect)
}

fn render_clip_enabled(renderer *Renderer) bool {
  return raw::render_clip_enabled(renderer)
}

fn set_render_scale(renderer *Renderer, scaleX f32, scaleY f32) bool {
  return raw::set_render_scale(renderer, scaleX, scaleY)
}

fn set_render_draw_color(renderer *Renderer, r u8, g u8, b u8, a u8) bool {
  return raw::set_render_draw_color(renderer, r, g, b, a)
}

fn set_render_draw_color_float(renderer *Renderer, r f32, g f32, b f32, a f32) bool {
  return raw::set_render_draw_color_float(renderer, r, g, b, a)
}

fn get_render_draw_color(renderer *Renderer, r *u8, g *u8, b *u8, a *u8) bool {
  return raw::get_render_draw_color(renderer, r, g, b, a)
}

fn set_render_color_scale(renderer *Renderer, scale f32) bool {
  return raw::set_render_color_scale(renderer, scale)
}

fn set_render_draw_blend_mode(renderer *Renderer, blendMode BlendMode) bool {
  return raw::set_render_draw_blend_mode(renderer, blendMode)
}

fn get_render_draw_blend_mode(renderer *Renderer, blendMode *BlendMode) bool {
  return raw::get_render_draw_blend_mode(renderer, blendMode)
}

fn render_clear(renderer *Renderer) bool {
  return raw::render_clear(renderer)
}

fn render_point(renderer *Renderer, x f32, y f32) bool {
  return raw::render_point(renderer, x, y)
}

fn render_points(renderer *Renderer, points const *FPoint, count i32) bool {
  return raw::render_points(renderer, points, count)
}

fn render_line(renderer *Renderer, x1 f32, y1 f32, x2 f32, y2 f32) bool {
  return raw::render_line(renderer, x1, y1, x2, y2)
}

fn render_lines(renderer *Renderer, points const *FPoint, count i32) bool {
  return raw::render_lines(renderer, points, count)
}

fn render_rect(renderer *Renderer, rect const *FRect) bool {
  return raw::render_rect(renderer, rect)
}

fn render_rects(renderer *Renderer, rects const *FRect, count i32) bool {
  return raw::render_rects(renderer, rects, count)
}

fn render_fill_rect(renderer *Renderer, rect const *FRect) bool {
  return raw::render_fill_rect(renderer, rect)
}

fn render_fill_rects(renderer *Renderer, rects const *FRect, count i32) bool {
  return raw::render_fill_rects(renderer, rects, count)
}

fn render_texture(renderer *Renderer, texture *Texture, srcrect const *FRect, dstrect const *FRect) bool {
  return raw::render_texture(renderer, texture, srcrect, dstrect)
}

fn render_texture_rotated(renderer *Renderer, texture *Texture, srcrect const *FRect, dstrect const *FRect, angle f64, center const *FPoint, flip FlipMode) bool {
  return raw::render_texture_rotated(renderer, texture, srcrect, dstrect, angle, center, flip)
}

fn render_texture_affine(renderer *Renderer, texture *Texture, srcrect const *FRect, origin const *FPoint, right const *FPoint, down const *FPoint) bool {
  return raw::render_texture_affine(renderer, texture, srcrect, origin, right, down)
}

fn render_texture_tiled(renderer *Renderer, texture *Texture, srcrect const *FRect, scale f32, dstrect const *FRect) bool {
  return raw::render_texture_tiled(renderer, texture, srcrect, scale, dstrect)
}

fn render_texture9_grid(renderer *Renderer, texture *Texture, srcrect const *FRect, left_width f32, right_width f32, top_height f32, bottom_height f32, scale f32, dstrect const *FRect) bool {
  return raw::render_texture9_grid(renderer, texture, srcrect, left_width, right_width, top_height, bottom_height, scale, dstrect)
}

fn render_read_pixels(renderer *Renderer, rect const *Rect) *Surface {
  return raw::render_read_pixels(renderer, rect)
}

fn render_present(renderer *Renderer) bool {
  return raw::render_present(renderer)
}

fn destroy_texture(texture *Texture) {
  raw::destroy_texture(texture)
}

fn destroy_renderer(renderer *Renderer) {
  raw::destroy_renderer(renderer)
}

fn flush_renderer(renderer *Renderer) bool {
  return raw::flush_renderer(renderer)
}

fn get_render_metal_layer(renderer *Renderer) rawptr {
  return raw::get_render_metal_layer(renderer)
}

fn get_render_metal_command_encoder(renderer *Renderer) rawptr {
  return raw::get_render_metal_command_encoder(renderer)
}

fn add_vulkan_render_semaphores(renderer *Renderer, wait_stage_mask u32, wait_semaphore i64, signal_semaphore i64) bool {
  return raw::add_vulkan_render_semaphores(renderer, wait_stage_mask, wait_semaphore, signal_semaphore)
}

fn set_render_v_sync(renderer *Renderer, vsync i32) bool {
  return raw::set_render_v_sync(renderer, vsync)
}

fn render_debug_text(renderer *Renderer, x f32, y f32, str string) error {
  str_c, str_err := strings::to_c(str)
  if str_err != nil {
    return str_err
  }
  defer strings::free_c(&str_c)
  if raw::render_debug_text(renderer, x, y, strings::c_view(&str_c)) == false {
    return error("SDL render_debug_text failed")
  }
  return nil
}

fn open_storage(iface const *StorageInterface, userdata rawptr) *Storage {
  return raw::open_storage(iface, userdata)
}

fn close_storage(storage *Storage) bool {
  return raw::close_storage(storage)
}

fn storage_ready(storage *Storage) bool {
  return raw::storage_ready(storage)
}

fn get_storage_space_remaining(storage *Storage) u64 {
  return raw::get_storage_space_remaining(storage)
}

fn is_tablet() bool {
  return raw::is_tablet()
}

fn is_tv() bool {
  return raw::is_tv()
}

fn get_sandbox() Sandbox {
  return raw::get_sandbox()
}

fn on_application_will_terminate() {
  raw::on_application_will_terminate()
}

fn on_application_did_receive_memory_warning() {
  raw::on_application_did_receive_memory_warning()
}

fn on_application_will_enter_background() {
  raw::on_application_will_enter_background()
}

fn on_application_did_enter_background() {
  raw::on_application_did_enter_background()
}

fn on_application_will_enter_foreground() {
  raw::on_application_will_enter_foreground()
}

fn on_application_did_enter_foreground() {
  raw::on_application_did_enter_foreground()
}

fn get_date_time_locale_preferences(dateFormat *DateFormat, timeFormat *TimeFormat) bool {
  return raw::get_date_time_locale_preferences(dateFormat, timeFormat)
}

fn get_current_time(ticks *Time) bool {
  return raw::get_current_time(ticks)
}

fn time_to_date_time(ticks Time, dt *DateTime, localTime bool) bool {
  return raw::time_to_date_time(ticks, dt, localTime)
}

fn date_time_to_time(dt const *DateTime, ticks *Time) bool {
  return raw::date_time_to_time(dt, ticks)
}

fn time_to_windows(ticks Time, dwLowDateTime *u32, dwHighDateTime *u32) {
  raw::time_to_windows(ticks, dwLowDateTime, dwHighDateTime)
}

fn time_from_windows(dwLowDateTime u32, dwHighDateTime u32) Time {
  return raw::time_from_windows(dwLowDateTime, dwHighDateTime)
}

fn get_days_in_month(year i32, month i32) i32 {
  return raw::get_days_in_month(year, month)
}

fn get_day_of_year(year i32, month i32, day i32) i32 {
  return raw::get_day_of_year(year, month, day)
}

fn get_day_of_week(year i32, month i32, day i32) i32 {
  return raw::get_day_of_week(year, month, day)
}

fn get_ticks() u64 {
  return raw::get_ticks()
}

fn get_ticks_ns() u64 {
  return raw::get_ticks_ns()
}

fn get_performance_counter() u64 {
  return raw::get_performance_counter()
}

fn get_performance_frequency() u64 {
  return raw::get_performance_frequency()
}

fn delay(ms u32) {
  raw::delay(ms)
}

fn delay_ns(ns u64) {
  raw::delay_ns(ns)
}

fn delay_precise(ns u64) {
  raw::delay_precise(ns)
}

fn remove_timer(id TimerID) bool {
  return raw::remove_timer(id)
}

fn set_tray_icon(tray *Tray, icon *Surface) {
  raw::set_tray_icon(tray, icon)
}

fn create_tray_menu(tray *Tray) *TrayMenu {
  return raw::create_tray_menu(tray)
}

fn create_tray_submenu(entry *TrayEntry) *TrayMenu {
  return raw::create_tray_submenu(entry)
}

fn get_tray_menu(tray *Tray) *TrayMenu {
  return raw::get_tray_menu(tray)
}

fn get_tray_submenu(entry *TrayEntry) *TrayMenu {
  return raw::get_tray_submenu(entry)
}

fn remove_tray_entry(entry *TrayEntry) {
  raw::remove_tray_entry(entry)
}

fn set_tray_entry_checked(entry *TrayEntry, checked bool) {
  raw::set_tray_entry_checked(entry, checked)
}

fn get_tray_entry_checked(entry *TrayEntry) bool {
  return raw::get_tray_entry_checked(entry)
}

fn set_tray_entry_enabled(entry *TrayEntry, enabled bool) {
  raw::set_tray_entry_enabled(entry, enabled)
}

fn get_tray_entry_enabled(entry *TrayEntry) bool {
  return raw::get_tray_entry_enabled(entry)
}

fn click_tray_entry(entry *TrayEntry) {
  raw::click_tray_entry(entry)
}

fn destroy_tray(tray *Tray) {
  raw::destroy_tray(tray)
}

fn get_tray_entry_parent(entry *TrayEntry) *TrayMenu {
  return raw::get_tray_entry_parent(entry)
}

fn get_tray_menu_parent_entry(menu *TrayMenu) *TrayEntry {
  return raw::get_tray_menu_parent_entry(menu)
}

fn get_tray_menu_parent_tray(menu *TrayMenu) *Tray {
  return raw::get_tray_menu_parent_tray(menu)
}

fn update_trays() {
  raw::update_trays()
}

fn get_version() i32 {
  return raw::get_version()
}

fn get_revision() string {
  result := raw::get_revision()
  if result == nil {
    return ""
  }
  return strings::from_c(result)
}

fn set_main_ready() {
  raw::set_main_ready()
}

fn gdk_suspend_complete() {
  raw::gdk_suspend_complete()
}

fn vulkan_get_vk_get_instance_proc_addr() FunctionPointer {
  return raw::vulkan_get_vk_get_instance_proc_addr()
}

fn vulkan_unload_library() {
  raw::vulkan_unload_library()
}

fn vulkan_create_surface(window *Window, instance VkInstance, allocator const *VkAllocationCallbacks, surface *VkSurfaceKHR) bool {
  return raw::vulkan_create_surface(window, instance, allocator, surface)
}

fn vulkan_destroy_surface(instance VkInstance, surface VkSurfaceKHR, allocator const *VkAllocationCallbacks) {
  raw::vulkan_destroy_surface(instance, surface, allocator)
}

fn vulkan_get_presentation_support(instance VkInstance, physicalDevice VkPhysicalDevice, queueFamilyIndex u32) bool {
  return raw::vulkan_get_presentation_support(instance, physicalDevice, queueFamilyIndex)
}

fn create_thread_with_properties(props PropertiesID) *Thread {
  return raw::create_thread_with_properties(props)
}
