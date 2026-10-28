package sdl

use "std::libc"

@[cc_macro="SDLCALL"]
type malloc_func = fn(arg0 libc::c_size_t) rawptr

@[cc_macro="SDLCALL"]
type calloc_func = fn(arg0 libc::c_size_t, arg1 libc::c_size_t) rawptr

@[cc_macro="SDLCALL"]
type realloc_func = fn(arg0 rawptr, arg1 libc::c_size_t) rawptr

@[cc_macro="SDLCALL"]
type free_func = fn(arg0 rawptr)

@[cc_macro="SDLCALL"]
type CompareCallback = fn(arg0 const *Void, arg1 const *Void) libc::c_int

@[cc_macro="SDLCALL"]
type CompareCallback_r = fn(arg0 rawptr, arg1 const *Void, arg2 const *Void) libc::c_int

type iconv_t = *iconv_data_t

@[cc_macro="SDLCALL"]
type AssertionHandler = fn(arg0 const *AssertData, arg1 rawptr) AssertState

type SpinLock = libc::c_int

type PropertiesID = u32

@[cc_macro="SDLCALL"]
type CleanupPropertyCallback = fn(arg0 rawptr, arg1 rawptr)

@[cc_macro="SDLCALL"]
type EnumeratePropertiesCallback = fn(arg0 rawptr, arg1 PropertiesID, arg2 const cstring)

@[cc_macro="SDLCALL"]
type ThreadFunction = fn(arg0 rawptr) libc::c_int

@[cc_macro="SDLCALL"]
type FunctionPointer = fn()

type ThreadID = u64

type TLSID = AtomicInt

@[cc_macro="SDLCALL"]
type TLSDestructorCallback = fn(arg0 rawptr)

@[cc_macro="SDLCALL"]
type IOStreamInterface_sizeCallback = fn(arg0 rawptr) i64

@[cc_macro="SDLCALL"]
type IOStreamInterface_seekCallback = fn(arg0 rawptr, arg1 i64, arg2 IOWhence) i64

@[cc_macro="SDLCALL"]
type IOStreamInterface_readCallback = fn(arg0 rawptr, arg1 rawptr, arg2 libc::c_size_t, arg3 *IOStatus) libc::c_size_t

@[cc_macro="SDLCALL"]
type IOStreamInterface_writeCallback = fn(arg0 rawptr, arg1 const *Void, arg2 libc::c_size_t, arg3 *IOStatus) libc::c_size_t

@[cc_macro="SDLCALL"]
type IOStreamInterface_flushCallback = fn(arg0 rawptr, arg1 *IOStatus) bool

@[cc_macro="SDLCALL"]
type IOStreamInterface_closeCallback = fn(arg0 rawptr) bool

type AudioDeviceID = u32

@[cc_macro="SDLCALL"]
type AudioStreamCallback = fn(arg0 rawptr, arg1 *AudioStream, arg2 libc::c_int, arg3 libc::c_int)

@[cc_macro="SDLCALL"]
type AudioPostmixCallback = fn(arg0 rawptr, arg1 const *AudioSpec, arg2 *libc::c_float, arg3 libc::c_int)

type BlendMode = u32

type SurfaceFlags = u32

type CameraID = u32

@[cc_macro="SDLCALL"]
type ClipboardDataCallback = fn(arg0 rawptr, arg1 const cstring, arg2 *libc::c_size_t) const *Void

@[cc_macro="SDLCALL"]
type ClipboardCleanupCallback = fn(arg0 rawptr)

type DisplayID = u32

type WindowFlags = u64

type WindowID = u32

@[cc_macro="SDLCALL"]
type HitTest = fn(arg0 *Window, arg1 const *Point, arg2 rawptr) HitTestResult

type GLContext = *GLContextState

type EGLDisplay = rawptr

type EGLConfig = rawptr

type EGLSurface = rawptr

type EGLAttrib = libc::c_intptr_t

@[cc_macro="SDLCALL"]
type EGLAttribArrayCallback = fn(arg0 rawptr) *EGLAttrib

type EGLint = libc::c_int

@[cc_macro="SDLCALL"]
type EGLIntArrayCallback = fn(arg0 rawptr, arg1 EGLDisplay, arg2 EGLConfig) *EGLint

@[cc_macro="SDLCALL"]
type DialogFileCallback = fn(arg0 rawptr, arg1 const *const cstring, arg2 libc::c_int)

type SensorID = u32

type JoystickID = u32

@[cc_macro="SDLCALL"]
type VirtualJoystickDesc_UpdateCallback = fn(arg0 rawptr)

@[cc_macro="SDLCALL"]
type VirtualJoystickDesc_SetPlayerIndexCallback = fn(arg0 rawptr, arg1 libc::c_int)

@[cc_macro="SDLCALL"]
type VirtualJoystickDesc_RumbleCallback = fn(arg0 rawptr, arg1 u16, arg2 u16) bool

@[cc_macro="SDLCALL"]
type VirtualJoystickDesc_RumbleTriggersCallback = fn(arg0 rawptr, arg1 u16, arg2 u16) bool

@[cc_macro="SDLCALL"]
type VirtualJoystickDesc_SetLEDCallback = fn(arg0 rawptr, arg1 u8, arg2 u8, arg3 u8) bool

@[cc_macro="SDLCALL"]
type VirtualJoystickDesc_SendEffectCallback = fn(arg0 rawptr, arg1 const *Void, arg2 libc::c_int) bool

@[cc_macro="SDLCALL"]
type VirtualJoystickDesc_SetSensorsEnabledCallback = fn(arg0 rawptr, arg1 bool) bool

@[cc_macro="SDLCALL"]
type VirtualJoystickDesc_CleanupCallback = fn(arg0 rawptr)

type KeyboardID = u32

type Keymod = u16

type Keycode = u32

type MouseID = u32

type MouseButtonFlags = u32

type TouchID = u64

type FingerID = u64

type PenID = u32

type PenInputFlags = u32

@[cc_macro="SDLCALL"]
type EventFilter = fn(arg0 rawptr, arg1 *Event) bool

@[cc_macro="SDLCALL"]
type EnumerateDirectoryCallback = fn(arg0 rawptr, arg1 const cstring, arg2 const cstring) EnumerationResult

type Time = i64

type GlobFlags = u32

type GPUShaderFormat = u32

type GPUColorComponentFlags = u8

type GPUTextureUsageFlags = u32

type GPUBufferUsageFlags = u32

type HapticID = u32

@[cc_macro="SDLCALL"]
type HintCallback = fn(arg0 rawptr, arg1 const cstring, arg2 const cstring, arg3 const cstring)

type InitFlags = u32

@[cc_macro="SDLCALL"]
type MainThreadCallback = fn(arg0 rawptr)

@[cc_macro="SDLCALL"]
type LogOutputFunction = fn(arg0 rawptr, arg1 libc::c_int, arg2 LogPriority, arg3 const cstring)

type MessageBoxFlags = u32

type MessageBoxButtonFlags = u32

type MetalView = rawptr

@[cc_macro="SDLCALL"]
type StorageInterface_closeCallback = fn(arg0 rawptr) bool

@[cc_macro="SDLCALL"]
type StorageInterface_readyCallback = fn(arg0 rawptr) bool

@[cc_macro="SDLCALL"]
type StorageInterface_enumerateCallback = fn(arg0 rawptr, arg1 const cstring, arg2 EnumerateDirectoryCallback, arg3 rawptr) bool

@[cc_macro="SDLCALL"]
type StorageInterface_infoCallback = fn(arg0 rawptr, arg1 const cstring, arg2 *PathInfo) bool

@[cc_macro="SDLCALL"]
type StorageInterface_read_fileCallback = fn(arg0 rawptr, arg1 const cstring, arg2 rawptr, arg3 u64) bool

@[cc_macro="SDLCALL"]
type StorageInterface_write_fileCallback = fn(arg0 rawptr, arg1 const cstring, arg2 const *Void, arg3 u64) bool

@[cc_macro="SDLCALL"]
type StorageInterface_mkdirCallback = fn(arg0 rawptr, arg1 const cstring) bool

@[cc_macro="SDLCALL"]
type StorageInterface_removeCallback = fn(arg0 rawptr, arg1 const cstring) bool

@[cc_macro="SDLCALL"]
type StorageInterface_renameCallback = fn(arg0 rawptr, arg1 const cstring, arg2 const cstring) bool

@[cc_macro="SDLCALL"]
type StorageInterface_copyCallback = fn(arg0 rawptr, arg1 const cstring, arg2 const cstring) bool

@[cc_macro="SDLCALL"]
type StorageInterface_space_remainingCallback = fn(arg0 rawptr) u64

@[cc_macro="SDLCALL"]
type X11EventHook = fn(arg0 rawptr, arg1 *XEvent) bool

type TimerID = u32

@[cc_macro="SDLCALL"]
type TimerCallback = fn(arg0 rawptr, arg1 TimerID, arg2 u32) u32

@[cc_macro="SDLCALL"]
type NSTimerCallback = fn(arg0 rawptr, arg1 TimerID, arg2 u64) u64

type TrayEntryFlags = u32

@[cc_macro="SDLCALL"]
type TrayCallback = fn(arg0 rawptr, arg1 *TrayEntry)

@[cc_macro="SDLCALL"]
type main_func = fn(arg0 libc::c_int, arg1 *cstring) libc::c_int

@[cc_macro="SDLCALL"]
type AppInit_func = fn(arg0 *rawptr, arg1 libc::c_int, arg2 *cstring) AppResult

@[cc_macro="SDLCALL"]
type AppIterate_func = fn(arg0 rawptr) AppResult

@[cc_macro="SDLCALL"]
type AppEvent_func = fn(arg0 rawptr, arg1 *Event) AppResult

@[cc_macro="SDLCALL"]
type AppQuit_func = fn(arg0 rawptr, arg1 AppResult)

type VkInstance = *VkInstance_T

type VkSurfaceKHR = *VkSurfaceKHR_T

type VkPhysicalDevice = *VkPhysicalDevice_T

@[header="SDL3/SDL.h"]
extern {
}

@[header="SDL3/SDL_vulkan.h"]
extern {
}

@[header="SDL3/SDL_main.h", default_cc_macro="SDLCALL"]
extern {




  @[link_name="SDL_Environment"]
  opaque Environment

  @[link_name="void"]
  opaque Void



  @[link_name="struct SDL_iconv_data_t"]
  opaque iconv_data_t


  @[link_name="SDL_AssertData"]
  struct AssertData {
    always_ignore bool
    trigger_count libc::c_uint
    condition const cstring
    filename const cstring
    linenum libc::c_int
    function const cstring
    next const *AssertData
  }

  @[link_name="SDL_AssertState"]
  enum AssertState {
    @[link_name="SDL_ASSERTION_RETRY"]
    RETRY
    @[link_name="SDL_ASSERTION_BREAK"]
    BREAK
    @[link_name="SDL_ASSERTION_ABORT"]
    ABORT
    @[link_name="SDL_ASSERTION_IGNORE"]
    IGNORE
    @[link_name="SDL_ASSERTION_ALWAYS_IGNORE"]
    ALWAYS_IGNORE
  }


  @[link_name="SDL_AsyncIO"]
  opaque AsyncIO

  @[link_name="SDL_AsyncIOQueue"]
  opaque AsyncIOQueue

  @[link_name="SDL_AsyncIOTaskType"]
  enum AsyncIOTaskType {
    @[link_name="SDL_ASYNCIO_TASK_READ"]
    READ
    @[link_name="SDL_ASYNCIO_TASK_WRITE"]
    WRITE
    @[link_name="SDL_ASYNCIO_TASK_CLOSE"]
    CLOSE
  }

  @[link_name="SDL_AsyncIOResult"]
  enum AsyncIOResult {
    @[link_name="SDL_ASYNCIO_COMPLETE"]
    COMPLETE
    @[link_name="SDL_ASYNCIO_FAILURE"]
    FAILURE
    @[link_name="SDL_ASYNCIO_CANCELED"]
    CANCELED
  }

  @[link_name="SDL_AsyncIOOutcome"]
  struct AsyncIOOutcome {
    asyncio *AsyncIO
    type AsyncIOTaskType
    result AsyncIOResult
    buffer rawptr
    offset u64
    bytes_requested u64
    bytes_transferred u64
    userdata rawptr
  }


  @[link_name="SDL_AtomicInt"]
  struct AtomicInt {
    value libc::c_int
  }

  @[link_name="SDL_AtomicU32"]
  struct AtomicU32 {
    value u32
  }



  @[link_name="SDL_PropertyType"]
  enum PropertyType {
    @[link_name="SDL_PROPERTY_TYPE_INVALID"]
    INVALID
    @[link_name="SDL_PROPERTY_TYPE_POINTER"]
    POINTER
    @[link_name="SDL_PROPERTY_TYPE_STRING"]
    STRING
    @[link_name="SDL_PROPERTY_TYPE_NUMBER"]
    NUMBER
    @[link_name="SDL_PROPERTY_TYPE_FLOAT"]
    FLOAT
    @[link_name="SDL_PROPERTY_TYPE_BOOLEAN"]
    BOOLEAN
  }




  @[link_name="SDL_Thread"]
  opaque Thread


  @[link_name="SDL_ThreadPriority"]
  enum ThreadPriority {
    @[link_name="SDL_THREAD_PRIORITY_LOW"]
    LOW
    @[link_name="SDL_THREAD_PRIORITY_NORMAL"]
    NORMAL
    @[link_name="SDL_THREAD_PRIORITY_HIGH"]
    HIGH
    @[link_name="SDL_THREAD_PRIORITY_TIME_CRITICAL"]
    TIME_CRITICAL
  }

  @[link_name="SDL_ThreadState"]
  enum ThreadState {
    @[link_name="SDL_THREAD_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_THREAD_ALIVE"]
    ALIVE
    @[link_name="SDL_THREAD_DETACHED"]
    DETACHED
    @[link_name="SDL_THREAD_COMPLETE"]
    COMPLETE
  }



  @[link_name="SDL_Mutex"]
  opaque Mutex

  @[link_name="SDL_RWLock"]
  opaque RWLock

  @[link_name="SDL_Semaphore"]
  opaque Semaphore

  @[link_name="SDL_Condition"]
  opaque Condition

  @[link_name="SDL_InitState"]
  struct InitState {
    status AtomicInt
    thread ThreadID
    reserved rawptr
  }

  @[link_name="SDL_IOStream"]
  opaque IOStream

  @[link_name="SDL_IOWhence"]
  enum IOWhence {
    @[link_name="SDL_IO_SEEK_SET"]
    SET
    @[link_name="SDL_IO_SEEK_CUR"]
    CUR
    @[link_name="SDL_IO_SEEK_END"]
    END
  }

  @[link_name="SDL_IOStatus"]
  enum IOStatus {
    @[link_name="SDL_IO_STATUS_READY"]
    READY
    @[link_name="SDL_IO_STATUS_ERROR"]
    ERROR
    @[link_name="SDL_IO_STATUS_EOF"]
    EOF
    @[link_name="SDL_IO_STATUS_NOT_READY"]
    NOT_READY
    @[link_name="SDL_IO_STATUS_READONLY"]
    READONLY
    @[link_name="SDL_IO_STATUS_WRITEONLY"]
    WRITEONLY
  }

  @[link_name="SDL_IOStreamInterface"]
  struct IOStreamInterface {
    version u32
    size IOStreamInterface_sizeCallback
    seek IOStreamInterface_seekCallback
    read IOStreamInterface_readCallback
    write IOStreamInterface_writeCallback
    flush IOStreamInterface_flushCallback
    close IOStreamInterface_closeCallback
  }


  @[link_name="SDL_AudioFormat"]
  enum AudioFormat {
    @[link_name="SDL_AUDIO_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_AUDIO_U8"]
    U8
    @[link_name="SDL_AUDIO_S8"]
    S8
    @[link_name="SDL_AUDIO_S16LE"]
    S16LE
    @[link_name="SDL_AUDIO_S16BE"]
    S16BE
    @[link_name="SDL_AUDIO_S32LE"]
    S32LE
    @[link_name="SDL_AUDIO_S32BE"]
    S32BE
    @[link_name="SDL_AUDIO_F32LE"]
    F32LE
    @[link_name="SDL_AUDIO_F32BE"]
    F32BE
    @[link_name="SDL_AUDIO_S16"]
    S16
    @[link_name="SDL_AUDIO_S32"]
    S32
    @[link_name="SDL_AUDIO_F32"]
    F32
  }

  @[link_name="SDL_AudioSpec"]
  struct AudioSpec {
    format AudioFormat
    channels libc::c_int
    freq libc::c_int
  }

  @[link_name="SDL_AudioStream"]
  opaque AudioStream



  @[link_name="SDL_BlendFactor"]
  enum BlendFactor {
    @[link_name="SDL_BLENDFACTOR_ZERO"]
    ZERO
    @[link_name="SDL_BLENDFACTOR_ONE"]
    ONE
    @[link_name="SDL_BLENDFACTOR_SRC_COLOR"]
    SRC_COLOR
    @[link_name="SDL_BLENDFACTOR_ONE_MINUS_SRC_COLOR"]
    ONE_MINUS_SRC_COLOR
    @[link_name="SDL_BLENDFACTOR_SRC_ALPHA"]
    SRC_ALPHA
    @[link_name="SDL_BLENDFACTOR_ONE_MINUS_SRC_ALPHA"]
    ONE_MINUS_SRC_ALPHA
    @[link_name="SDL_BLENDFACTOR_DST_COLOR"]
    DST_COLOR
    @[link_name="SDL_BLENDFACTOR_ONE_MINUS_DST_COLOR"]
    ONE_MINUS_DST_COLOR
    @[link_name="SDL_BLENDFACTOR_DST_ALPHA"]
    DST_ALPHA
    @[link_name="SDL_BLENDFACTOR_ONE_MINUS_DST_ALPHA"]
    ONE_MINUS_DST_ALPHA
  }

  @[link_name="SDL_BlendOperation"]
  enum BlendOperation {
    @[link_name="SDL_BLENDOPERATION_ADD"]
    ADD
    @[link_name="SDL_BLENDOPERATION_SUBTRACT"]
    SUBTRACT
    @[link_name="SDL_BLENDOPERATION_REV_SUBTRACT"]
    REV_SUBTRACT
    @[link_name="SDL_BLENDOPERATION_MINIMUM"]
    MINIMUM
    @[link_name="SDL_BLENDOPERATION_MAXIMUM"]
    MAXIMUM
  }


  @[link_name="SDL_PixelFormat"]
  enum PixelFormat {
    @[link_name="SDL_PIXELFORMAT_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_PIXELFORMAT_INDEX1LSB"]
    INDEX1LSB
    @[link_name="SDL_PIXELFORMAT_INDEX1MSB"]
    INDEX1MSB
    @[link_name="SDL_PIXELFORMAT_INDEX2LSB"]
    INDEX2LSB
    @[link_name="SDL_PIXELFORMAT_INDEX2MSB"]
    INDEX2MSB
    @[link_name="SDL_PIXELFORMAT_INDEX4LSB"]
    INDEX4LSB
    @[link_name="SDL_PIXELFORMAT_INDEX4MSB"]
    INDEX4MSB
    @[link_name="SDL_PIXELFORMAT_INDEX8"]
    INDEX8
    @[link_name="SDL_PIXELFORMAT_RGB332"]
    RGB332
    @[link_name="SDL_PIXELFORMAT_XRGB4444"]
    XRGB4444
    @[link_name="SDL_PIXELFORMAT_XBGR4444"]
    XBGR4444
    @[link_name="SDL_PIXELFORMAT_XRGB1555"]
    XRGB1555
    @[link_name="SDL_PIXELFORMAT_XBGR1555"]
    XBGR1555
    @[link_name="SDL_PIXELFORMAT_ARGB4444"]
    ARGB4444
    @[link_name="SDL_PIXELFORMAT_RGBA4444"]
    RGBA4444
    @[link_name="SDL_PIXELFORMAT_ABGR4444"]
    ABGR4444
    @[link_name="SDL_PIXELFORMAT_BGRA4444"]
    BGRA4444
    @[link_name="SDL_PIXELFORMAT_ARGB1555"]
    ARGB1555
    @[link_name="SDL_PIXELFORMAT_RGBA5551"]
    RGBA5551
    @[link_name="SDL_PIXELFORMAT_ABGR1555"]
    ABGR1555
    @[link_name="SDL_PIXELFORMAT_BGRA5551"]
    BGRA5551
    @[link_name="SDL_PIXELFORMAT_RGB565"]
    RGB565
    @[link_name="SDL_PIXELFORMAT_BGR565"]
    BGR565
    @[link_name="SDL_PIXELFORMAT_RGB24"]
    RGB24
    @[link_name="SDL_PIXELFORMAT_BGR24"]
    BGR24
    @[link_name="SDL_PIXELFORMAT_XRGB8888"]
    XRGB8888
    @[link_name="SDL_PIXELFORMAT_RGBX8888"]
    RGBX8888
    @[link_name="SDL_PIXELFORMAT_XBGR8888"]
    XBGR8888
    @[link_name="SDL_PIXELFORMAT_BGRX8888"]
    BGRX8888
    @[link_name="SDL_PIXELFORMAT_ARGB8888"]
    ARGB8888
    @[link_name="SDL_PIXELFORMAT_RGBA8888"]
    RGBA8888
    @[link_name="SDL_PIXELFORMAT_ABGR8888"]
    ABGR8888
    @[link_name="SDL_PIXELFORMAT_BGRA8888"]
    BGRA8888
    @[link_name="SDL_PIXELFORMAT_XRGB2101010"]
    XRGB2101010
    @[link_name="SDL_PIXELFORMAT_XBGR2101010"]
    XBGR2101010
    @[link_name="SDL_PIXELFORMAT_ARGB2101010"]
    ARGB2101010
    @[link_name="SDL_PIXELFORMAT_ABGR2101010"]
    ABGR2101010
    @[link_name="SDL_PIXELFORMAT_RGB48"]
    RGB48
    @[link_name="SDL_PIXELFORMAT_BGR48"]
    BGR48
    @[link_name="SDL_PIXELFORMAT_RGBA64"]
    RGBA64
    @[link_name="SDL_PIXELFORMAT_ARGB64"]
    ARGB64
    @[link_name="SDL_PIXELFORMAT_BGRA64"]
    BGRA64
    @[link_name="SDL_PIXELFORMAT_ABGR64"]
    ABGR64
    @[link_name="SDL_PIXELFORMAT_RGB48_FLOAT"]
    RGB48_FLOAT
    @[link_name="SDL_PIXELFORMAT_BGR48_FLOAT"]
    BGR48_FLOAT
    @[link_name="SDL_PIXELFORMAT_RGBA64_FLOAT"]
    RGBA64_FLOAT
    @[link_name="SDL_PIXELFORMAT_ARGB64_FLOAT"]
    ARGB64_FLOAT
    @[link_name="SDL_PIXELFORMAT_BGRA64_FLOAT"]
    BGRA64_FLOAT
    @[link_name="SDL_PIXELFORMAT_ABGR64_FLOAT"]
    ABGR64_FLOAT
    @[link_name="SDL_PIXELFORMAT_RGB96_FLOAT"]
    RGB96_FLOAT
    @[link_name="SDL_PIXELFORMAT_BGR96_FLOAT"]
    BGR96_FLOAT
    @[link_name="SDL_PIXELFORMAT_RGBA128_FLOAT"]
    RGBA128_FLOAT
    @[link_name="SDL_PIXELFORMAT_ARGB128_FLOAT"]
    ARGB128_FLOAT
    @[link_name="SDL_PIXELFORMAT_BGRA128_FLOAT"]
    BGRA128_FLOAT
    @[link_name="SDL_PIXELFORMAT_ABGR128_FLOAT"]
    ABGR128_FLOAT
    @[link_name="SDL_PIXELFORMAT_YV12"]
    YV12
    @[link_name="SDL_PIXELFORMAT_IYUV"]
    IYUV
    @[link_name="SDL_PIXELFORMAT_YUY2"]
    YUY2
    @[link_name="SDL_PIXELFORMAT_UYVY"]
    UYVY
    @[link_name="SDL_PIXELFORMAT_YVYU"]
    YVYU
    @[link_name="SDL_PIXELFORMAT_NV12"]
    NV12
    @[link_name="SDL_PIXELFORMAT_NV21"]
    NV21
    @[link_name="SDL_PIXELFORMAT_P010"]
    P010
    @[link_name="SDL_PIXELFORMAT_EXTERNAL_OES"]
    EXTERNAL_OES
    @[link_name="SDL_PIXELFORMAT_MJPG"]
    MJPG
    @[link_name="SDL_PIXELFORMAT_RGBA32"]
    RGBA32
    @[link_name="SDL_PIXELFORMAT_ARGB32"]
    ARGB32
    @[link_name="SDL_PIXELFORMAT_BGRA32"]
    BGRA32
    @[link_name="SDL_PIXELFORMAT_ABGR32"]
    ABGR32
    @[link_name="SDL_PIXELFORMAT_RGBX32"]
    RGBX32
    @[link_name="SDL_PIXELFORMAT_XRGB32"]
    XRGB32
    @[link_name="SDL_PIXELFORMAT_BGRX32"]
    BGRX32
    @[link_name="SDL_PIXELFORMAT_XBGR32"]
    XBGR32
  }

  @[link_name="SDL_PixelFormatDetails"]
  struct PixelFormatDetails {
    format PixelFormat
    bits_per_pixel u8
    bytes_per_pixel u8
    padding u8[2]
    Rmask u32
    Gmask u32
    Bmask u32
    Amask u32
    Rbits u8
    Gbits u8
    Bbits u8
    Abits u8
    Rshift u8
    Gshift u8
    Bshift u8
    Ashift u8
  }

  @[link_name="SDL_Color"]
  struct Color {
    r u8
    g u8
    b u8
    a u8
  }

  @[link_name="SDL_Palette"]
  struct Palette {
    ncolors libc::c_int
    colors *Color
    version u32
    refcount libc::c_int
  }

  @[link_name="SDL_Rect"]
  struct Rect {
    x libc::c_int
    y libc::c_int
    w libc::c_int
    h libc::c_int
  }

  @[link_name="SDL_FRect"]
  struct FRect {
    x libc::c_float
    y libc::c_float
    w libc::c_float
    h libc::c_float
  }

  @[link_name="SDL_Point"]
  struct Point {
    x libc::c_int
    y libc::c_int
  }

  @[link_name="SDL_FPoint"]
  struct FPoint {
    x libc::c_float
    y libc::c_float
  }


  @[link_name="SDL_Surface"]
  struct Surface {
    flags SurfaceFlags
    format PixelFormat
    w libc::c_int
    h libc::c_int
    pitch libc::c_int
    pixels rawptr
    refcount libc::c_int
    reserved rawptr
  }

  @[link_name="SDL_Colorspace"]
  enum Colorspace {
    @[link_name="SDL_COLORSPACE_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_COLORSPACE_SRGB"]
    SRGB
    @[link_name="SDL_COLORSPACE_SRGB_LINEAR"]
    SRGB_LINEAR
    @[link_name="SDL_COLORSPACE_HDR10"]
    HDR10
    @[link_name="SDL_COLORSPACE_JPEG"]
    JPEG
    @[link_name="SDL_COLORSPACE_BT601_LIMITED"]
    BT601_LIMITED
    @[link_name="SDL_COLORSPACE_BT601_FULL"]
    BT601_FULL
    @[link_name="SDL_COLORSPACE_BT709_LIMITED"]
    BT709_LIMITED
    @[link_name="SDL_COLORSPACE_BT709_FULL"]
    BT709_FULL
    @[link_name="SDL_COLORSPACE_BT2020_LIMITED"]
    BT2020_LIMITED
    @[link_name="SDL_COLORSPACE_BT2020_FULL"]
    BT2020_FULL
    @[link_name="SDL_COLORSPACE_RGB_DEFAULT"]
    RGB_DEFAULT
    @[link_name="SDL_COLORSPACE_YUV_DEFAULT"]
    YUV_DEFAULT
  }

  @[link_name="SDL_FlipMode"]
  enum FlipMode {
    @[link_name="SDL_FLIP_NONE"]
    NONE
    @[link_name="SDL_FLIP_HORIZONTAL"]
    HORIZONTAL
    @[link_name="SDL_FLIP_VERTICAL"]
    VERTICAL
  }

  @[link_name="SDL_ScaleMode"]
  enum ScaleMode {
    @[link_name="SDL_SCALEMODE_INVALID"]
    INVALID
    @[link_name="SDL_SCALEMODE_NEAREST"]
    NEAREST
    @[link_name="SDL_SCALEMODE_LINEAR"]
    LINEAR
  }


  @[link_name="SDL_CameraSpec"]
  struct CameraSpec {
    format PixelFormat
    colorspace Colorspace
    width libc::c_int
    height libc::c_int
    framerate_numerator libc::c_int
    framerate_denominator libc::c_int
  }

  @[link_name="SDL_CameraPosition"]
  enum CameraPosition {
    @[link_name="SDL_CAMERA_POSITION_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_CAMERA_POSITION_FRONT_FACING"]
    FRONT_FACING
    @[link_name="SDL_CAMERA_POSITION_BACK_FACING"]
    BACK_FACING
  }

  @[link_name="SDL_Camera"]
  opaque Camera



  @[link_name="SDL_SystemTheme"]
  enum SystemTheme {
    @[link_name="SDL_SYSTEM_THEME_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_SYSTEM_THEME_LIGHT"]
    LIGHT
    @[link_name="SDL_SYSTEM_THEME_DARK"]
    DARK
  }


  @[link_name="SDL_DisplayOrientation"]
  enum DisplayOrientation {
    @[link_name="SDL_ORIENTATION_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_ORIENTATION_LANDSCAPE"]
    LANDSCAPE
    @[link_name="SDL_ORIENTATION_LANDSCAPE_FLIPPED"]
    LANDSCAPE_FLIPPED
    @[link_name="SDL_ORIENTATION_PORTRAIT"]
    PORTRAIT
    @[link_name="SDL_ORIENTATION_PORTRAIT_FLIPPED"]
    PORTRAIT_FLIPPED
  }

  @[link_name="SDL_DisplayModeData"]
  opaque DisplayModeData

  @[link_name="SDL_DisplayMode"]
  struct DisplayMode {
    displayID DisplayID
    format PixelFormat
    w libc::c_int
    h libc::c_int
    pixel_density libc::c_float
    refresh_rate libc::c_float
    refresh_rate_numerator libc::c_int
    refresh_rate_denominator libc::c_int
    internal *DisplayModeData
  }

  @[link_name="SDL_Window"]
  opaque Window



  @[link_name="SDL_HitTestResult"]
  enum HitTestResult {
    @[link_name="SDL_HITTEST_NORMAL"]
    NORMAL
    @[link_name="SDL_HITTEST_DRAGGABLE"]
    DRAGGABLE
    @[link_name="SDL_HITTEST_RESIZE_TOPLEFT"]
    RESIZE_TOPLEFT
    @[link_name="SDL_HITTEST_RESIZE_TOP"]
    RESIZE_TOP
    @[link_name="SDL_HITTEST_RESIZE_TOPRIGHT"]
    RESIZE_TOPRIGHT
    @[link_name="SDL_HITTEST_RESIZE_RIGHT"]
    RESIZE_RIGHT
    @[link_name="SDL_HITTEST_RESIZE_BOTTOMRIGHT"]
    RESIZE_BOTTOMRIGHT
    @[link_name="SDL_HITTEST_RESIZE_BOTTOM"]
    RESIZE_BOTTOM
    @[link_name="SDL_HITTEST_RESIZE_BOTTOMLEFT"]
    RESIZE_BOTTOMLEFT
    @[link_name="SDL_HITTEST_RESIZE_LEFT"]
    RESIZE_LEFT
  }


  @[link_name="SDL_FlashOperation"]
  enum FlashOperation {
    @[link_name="SDL_FLASH_CANCEL"]
    CANCEL
    @[link_name="SDL_FLASH_BRIEFLY"]
    BRIEFLY
    @[link_name="SDL_FLASH_UNTIL_FOCUSED"]
    UNTIL_FOCUSED
  }

  @[link_name="SDL_GLAttr"]
  enum GLAttr {
    @[link_name="SDL_GL_RED_SIZE"]
    RED_SIZE
    @[link_name="SDL_GL_GREEN_SIZE"]
    GREEN_SIZE
    @[link_name="SDL_GL_BLUE_SIZE"]
    BLUE_SIZE
    @[link_name="SDL_GL_ALPHA_SIZE"]
    ALPHA_SIZE
    @[link_name="SDL_GL_BUFFER_SIZE"]
    BUFFER_SIZE
    @[link_name="SDL_GL_DOUBLEBUFFER"]
    DOUBLEBUFFER
    @[link_name="SDL_GL_DEPTH_SIZE"]
    DEPTH_SIZE
    @[link_name="SDL_GL_STENCIL_SIZE"]
    STENCIL_SIZE
    @[link_name="SDL_GL_ACCUM_RED_SIZE"]
    ACCUM_RED_SIZE
    @[link_name="SDL_GL_ACCUM_GREEN_SIZE"]
    ACCUM_GREEN_SIZE
    @[link_name="SDL_GL_ACCUM_BLUE_SIZE"]
    ACCUM_BLUE_SIZE
    @[link_name="SDL_GL_ACCUM_ALPHA_SIZE"]
    ACCUM_ALPHA_SIZE
    @[link_name="SDL_GL_STEREO"]
    STEREO
    @[link_name="SDL_GL_MULTISAMPLEBUFFERS"]
    MULTISAMPLEBUFFERS
    @[link_name="SDL_GL_MULTISAMPLESAMPLES"]
    MULTISAMPLESAMPLES
    @[link_name="SDL_GL_ACCELERATED_VISUAL"]
    ACCELERATED_VISUAL
    @[link_name="SDL_GL_RETAINED_BACKING"]
    RETAINED_BACKING
    @[link_name="SDL_GL_CONTEXT_MAJOR_VERSION"]
    CONTEXT_MAJOR_VERSION
    @[link_name="SDL_GL_CONTEXT_MINOR_VERSION"]
    CONTEXT_MINOR_VERSION
    @[link_name="SDL_GL_CONTEXT_FLAGS"]
    CONTEXT_FLAGS
    @[link_name="SDL_GL_CONTEXT_PROFILE_MASK"]
    CONTEXT_PROFILE_MASK
    @[link_name="SDL_GL_SHARE_WITH_CURRENT_CONTEXT"]
    SHARE_WITH_CURRENT_CONTEXT
    @[link_name="SDL_GL_FRAMEBUFFER_SRGB_CAPABLE"]
    FRAMEBUFFER_SRGB_CAPABLE
    @[link_name="SDL_GL_CONTEXT_RELEASE_BEHAVIOR"]
    CONTEXT_RELEASE_BEHAVIOR
    @[link_name="SDL_GL_CONTEXT_RESET_NOTIFICATION"]
    CONTEXT_RESET_NOTIFICATION
    @[link_name="SDL_GL_CONTEXT_NO_ERROR"]
    CONTEXT_NO_ERROR
    @[link_name="SDL_GL_FLOATBUFFERS"]
    FLOATBUFFERS
    @[link_name="SDL_GL_EGL_PLATFORM"]
    EGL_PLATFORM
  }

  @[link_name="struct SDL_GLContextState"]
  opaque GLContextState










  @[link_name="SDL_DialogFileFilter"]
  struct DialogFileFilter {
    name const cstring
    pattern const cstring
  }

  @[link_name="SDL_FileDialogType"]
  enum FileDialogType {
    @[link_name="SDL_FILEDIALOG_OPENFILE"]
    OPENFILE
    @[link_name="SDL_FILEDIALOG_SAVEFILE"]
    SAVEFILE
    @[link_name="SDL_FILEDIALOG_OPENFOLDER"]
    OPENFOLDER
  }

  @[link_name="SDL_GUID"]
  struct GUID {
    data u8[16]
  }

  @[link_name="SDL_PowerState"]
  enum PowerState {
    @[link_name="SDL_POWERSTATE_ERROR"]
    ERROR
    @[link_name="SDL_POWERSTATE_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_POWERSTATE_ON_BATTERY"]
    ON_BATTERY
    @[link_name="SDL_POWERSTATE_NO_BATTERY"]
    NO_BATTERY
    @[link_name="SDL_POWERSTATE_CHARGING"]
    CHARGING
    @[link_name="SDL_POWERSTATE_CHARGED"]
    CHARGED
  }


  @[link_name="SDL_SensorType"]
  enum SensorType {
    @[link_name="SDL_SENSOR_INVALID"]
    INVALID
    @[link_name="SDL_SENSOR_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_SENSOR_ACCEL"]
    ACCEL
    @[link_name="SDL_SENSOR_GYRO"]
    GYRO
    @[link_name="SDL_SENSOR_ACCEL_L"]
    ACCEL_L
    @[link_name="SDL_SENSOR_GYRO_L"]
    GYRO_L
    @[link_name="SDL_SENSOR_ACCEL_R"]
    ACCEL_R
    @[link_name="SDL_SENSOR_GYRO_R"]
    GYRO_R
    @[link_name="SDL_SENSOR_COUNT"]
    COUNT
  }

  @[link_name="SDL_Sensor"]
  opaque Sensor


  @[link_name="SDL_JoystickType"]
  enum JoystickType {
    @[link_name="SDL_JOYSTICK_TYPE_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_JOYSTICK_TYPE_GAMEPAD"]
    GAMEPAD
    @[link_name="SDL_JOYSTICK_TYPE_WHEEL"]
    WHEEL
    @[link_name="SDL_JOYSTICK_TYPE_ARCADE_STICK"]
    ARCADE_STICK
    @[link_name="SDL_JOYSTICK_TYPE_FLIGHT_STICK"]
    FLIGHT_STICK
    @[link_name="SDL_JOYSTICK_TYPE_DANCE_PAD"]
    DANCE_PAD
    @[link_name="SDL_JOYSTICK_TYPE_GUITAR"]
    GUITAR
    @[link_name="SDL_JOYSTICK_TYPE_DRUM_KIT"]
    DRUM_KIT
    @[link_name="SDL_JOYSTICK_TYPE_ARCADE_PAD"]
    ARCADE_PAD
    @[link_name="SDL_JOYSTICK_TYPE_THROTTLE"]
    THROTTLE
    @[link_name="SDL_JOYSTICK_TYPE_COUNT"]
    COUNT
  }

  @[link_name="SDL_Joystick"]
  opaque Joystick

  @[link_name="SDL_VirtualJoystickTouchpadDesc"]
  struct VirtualJoystickTouchpadDesc {
    nfingers u16
    padding u16[3]
  }

  @[link_name="SDL_VirtualJoystickSensorDesc"]
  struct VirtualJoystickSensorDesc {
    type SensorType
    rate libc::c_float
  }

  @[link_name="SDL_VirtualJoystickDesc"]
  struct VirtualJoystickDesc {
    version u32
    type u16
    padding u16
    vendor_id u16
    product_id u16
    naxes u16
    nbuttons u16
    nballs u16
    nhats u16
    ntouchpads u16
    nsensors u16
    padding2 u16[2]
    button_mask u32
    axis_mask u32
    name const cstring
    touchpads const *VirtualJoystickTouchpadDesc
    sensors const *VirtualJoystickSensorDesc
    userdata rawptr
    Update VirtualJoystickDesc_UpdateCallback
    SetPlayerIndex VirtualJoystickDesc_SetPlayerIndexCallback
    Rumble VirtualJoystickDesc_RumbleCallback
    RumbleTriggers VirtualJoystickDesc_RumbleTriggersCallback
    SetLED VirtualJoystickDesc_SetLEDCallback
    SendEffect VirtualJoystickDesc_SendEffectCallback
    SetSensorsEnabled VirtualJoystickDesc_SetSensorsEnabledCallback
    Cleanup VirtualJoystickDesc_CleanupCallback
  }

  @[link_name="SDL_JoystickConnectionState"]
  enum JoystickConnectionState {
    @[link_name="SDL_JOYSTICK_CONNECTION_INVALID"]
    INVALID
    @[link_name="SDL_JOYSTICK_CONNECTION_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_JOYSTICK_CONNECTION_WIRED"]
    WIRED
    @[link_name="SDL_JOYSTICK_CONNECTION_WIRELESS"]
    WIRELESS
  }

  @[link_name="SDL_Gamepad"]
  opaque Gamepad

  @[link_name="SDL_GamepadType"]
  enum GamepadType {
    @[link_name="SDL_GAMEPAD_TYPE_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_GAMEPAD_TYPE_STANDARD"]
    STANDARD
    @[link_name="SDL_GAMEPAD_TYPE_XBOX360"]
    XBOX360
    @[link_name="SDL_GAMEPAD_TYPE_XBOXONE"]
    XBOXONE
    @[link_name="SDL_GAMEPAD_TYPE_PS3"]
    PS3
    @[link_name="SDL_GAMEPAD_TYPE_PS4"]
    PS4
    @[link_name="SDL_GAMEPAD_TYPE_PS5"]
    PS5
    @[link_name="SDL_GAMEPAD_TYPE_NINTENDO_SWITCH_PRO"]
    NINTENDO_SWITCH_PRO
    @[link_name="SDL_GAMEPAD_TYPE_NINTENDO_SWITCH_JOYCON_LEFT"]
    NINTENDO_SWITCH_JOYCON_LEFT
    @[link_name="SDL_GAMEPAD_TYPE_NINTENDO_SWITCH_JOYCON_RIGHT"]
    NINTENDO_SWITCH_JOYCON_RIGHT
    @[link_name="SDL_GAMEPAD_TYPE_NINTENDO_SWITCH_JOYCON_PAIR"]
    NINTENDO_SWITCH_JOYCON_PAIR
    @[link_name="SDL_GAMEPAD_TYPE_COUNT"]
    COUNT
  }

  @[link_name="SDL_GamepadBindingType"]
  enum GamepadBindingType {
    @[link_name="SDL_GAMEPAD_BINDTYPE_NONE"]
    NONE
    @[link_name="SDL_GAMEPAD_BINDTYPE_BUTTON"]
    BUTTON
    @[link_name="SDL_GAMEPAD_BINDTYPE_AXIS"]
    AXIS
    @[link_name="SDL_GAMEPAD_BINDTYPE_HAT"]
    HAT
  }

  @[link_name="__typeof__(((__typeof__(((SDL_GamepadBinding *)0)->input) *)0)->axis)"]
  struct GamepadBindingInputAxis {
    axis libc::c_int
    axis_min libc::c_int
    axis_max libc::c_int
  }

  @[link_name="__typeof__(((__typeof__(((SDL_GamepadBinding *)0)->input) *)0)->hat)"]
  struct GamepadBindingInputHat {
    hat libc::c_int
    hat_mask libc::c_int
  }

  @[link_name="__typeof__(((SDL_GamepadBinding *)0)->input)"]
  union GamepadBindingInput {
    button libc::c_int
    axis GamepadBindingInputAxis
    hat GamepadBindingInputHat
  }

  @[link_name="SDL_GamepadButton"]
  enum GamepadButton {
    @[link_name="SDL_GAMEPAD_BUTTON_INVALID"]
    INVALID
    @[link_name="SDL_GAMEPAD_BUTTON_SOUTH"]
    SOUTH
    @[link_name="SDL_GAMEPAD_BUTTON_EAST"]
    EAST
    @[link_name="SDL_GAMEPAD_BUTTON_WEST"]
    WEST
    @[link_name="SDL_GAMEPAD_BUTTON_NORTH"]
    NORTH
    @[link_name="SDL_GAMEPAD_BUTTON_BACK"]
    BACK
    @[link_name="SDL_GAMEPAD_BUTTON_GUIDE"]
    GUIDE
    @[link_name="SDL_GAMEPAD_BUTTON_START"]
    START
    @[link_name="SDL_GAMEPAD_BUTTON_LEFT_STICK"]
    LEFT_STICK
    @[link_name="SDL_GAMEPAD_BUTTON_RIGHT_STICK"]
    RIGHT_STICK
    @[link_name="SDL_GAMEPAD_BUTTON_LEFT_SHOULDER"]
    LEFT_SHOULDER
    @[link_name="SDL_GAMEPAD_BUTTON_RIGHT_SHOULDER"]
    RIGHT_SHOULDER
    @[link_name="SDL_GAMEPAD_BUTTON_DPAD_UP"]
    DPAD_UP
    @[link_name="SDL_GAMEPAD_BUTTON_DPAD_DOWN"]
    DPAD_DOWN
    @[link_name="SDL_GAMEPAD_BUTTON_DPAD_LEFT"]
    DPAD_LEFT
    @[link_name="SDL_GAMEPAD_BUTTON_DPAD_RIGHT"]
    DPAD_RIGHT
    @[link_name="SDL_GAMEPAD_BUTTON_MISC1"]
    MISC1
    @[link_name="SDL_GAMEPAD_BUTTON_RIGHT_PADDLE1"]
    RIGHT_PADDLE1
    @[link_name="SDL_GAMEPAD_BUTTON_LEFT_PADDLE1"]
    LEFT_PADDLE1
    @[link_name="SDL_GAMEPAD_BUTTON_RIGHT_PADDLE2"]
    RIGHT_PADDLE2
    @[link_name="SDL_GAMEPAD_BUTTON_LEFT_PADDLE2"]
    LEFT_PADDLE2
    @[link_name="SDL_GAMEPAD_BUTTON_TOUCHPAD"]
    TOUCHPAD
    @[link_name="SDL_GAMEPAD_BUTTON_MISC2"]
    MISC2
    @[link_name="SDL_GAMEPAD_BUTTON_MISC3"]
    MISC3
    @[link_name="SDL_GAMEPAD_BUTTON_MISC4"]
    MISC4
    @[link_name="SDL_GAMEPAD_BUTTON_MISC5"]
    MISC5
    @[link_name="SDL_GAMEPAD_BUTTON_MISC6"]
    MISC6
    @[link_name="SDL_GAMEPAD_BUTTON_COUNT"]
    COUNT
  }

  @[link_name="SDL_GamepadAxis"]
  enum GamepadAxis {
    @[link_name="SDL_GAMEPAD_AXIS_INVALID"]
    INVALID
    @[link_name="SDL_GAMEPAD_AXIS_LEFTX"]
    LEFTX
    @[link_name="SDL_GAMEPAD_AXIS_LEFTY"]
    LEFTY
    @[link_name="SDL_GAMEPAD_AXIS_RIGHTX"]
    RIGHTX
    @[link_name="SDL_GAMEPAD_AXIS_RIGHTY"]
    RIGHTY
    @[link_name="SDL_GAMEPAD_AXIS_LEFT_TRIGGER"]
    LEFT_TRIGGER
    @[link_name="SDL_GAMEPAD_AXIS_RIGHT_TRIGGER"]
    RIGHT_TRIGGER
    @[link_name="SDL_GAMEPAD_AXIS_COUNT"]
    COUNT
  }

  @[link_name="__typeof__(((__typeof__(((SDL_GamepadBinding *)0)->output) *)0)->axis)"]
  struct GamepadBindingOutputAxis {
    axis GamepadAxis
    axis_min libc::c_int
    axis_max libc::c_int
  }

  @[link_name="__typeof__(((SDL_GamepadBinding *)0)->output)"]
  union GamepadBindingOutput {
    button GamepadButton
    axis GamepadBindingOutputAxis
  }

  @[link_name="SDL_GamepadBinding"]
  struct GamepadBinding {
    input_type GamepadBindingType
    input GamepadBindingInput
    output_type GamepadBindingType
    output GamepadBindingOutput
  }

  @[link_name="SDL_GamepadButtonLabel"]
  enum GamepadButtonLabel {
    @[link_name="SDL_GAMEPAD_BUTTON_LABEL_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_GAMEPAD_BUTTON_LABEL_A"]
    A
    @[link_name="SDL_GAMEPAD_BUTTON_LABEL_B"]
    B
    @[link_name="SDL_GAMEPAD_BUTTON_LABEL_X"]
    X
    @[link_name="SDL_GAMEPAD_BUTTON_LABEL_Y"]
    Y
    @[link_name="SDL_GAMEPAD_BUTTON_LABEL_CROSS"]
    CROSS
    @[link_name="SDL_GAMEPAD_BUTTON_LABEL_CIRCLE"]
    CIRCLE
    @[link_name="SDL_GAMEPAD_BUTTON_LABEL_SQUARE"]
    SQUARE
    @[link_name="SDL_GAMEPAD_BUTTON_LABEL_TRIANGLE"]
    TRIANGLE
  }



  @[link_name="SDL_Scancode"]
  enum Scancode {
    @[link_name="SDL_SCANCODE_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_SCANCODE_A"]
    A
    @[link_name="SDL_SCANCODE_B"]
    B
    @[link_name="SDL_SCANCODE_C"]
    C
    @[link_name="SDL_SCANCODE_D"]
    D
    @[link_name="SDL_SCANCODE_E"]
    E
    @[link_name="SDL_SCANCODE_F"]
    F
    @[link_name="SDL_SCANCODE_G"]
    G
    @[link_name="SDL_SCANCODE_H"]
    H
    @[link_name="SDL_SCANCODE_I"]
    I
    @[link_name="SDL_SCANCODE_J"]
    J
    @[link_name="SDL_SCANCODE_K"]
    K
    @[link_name="SDL_SCANCODE_L"]
    L
    @[link_name="SDL_SCANCODE_M"]
    M
    @[link_name="SDL_SCANCODE_N"]
    N
    @[link_name="SDL_SCANCODE_O"]
    O
    @[link_name="SDL_SCANCODE_P"]
    P
    @[link_name="SDL_SCANCODE_Q"]
    Q
    @[link_name="SDL_SCANCODE_R"]
    R
    @[link_name="SDL_SCANCODE_S"]
    S
    @[link_name="SDL_SCANCODE_T"]
    T
    @[link_name="SDL_SCANCODE_U"]
    U
    @[link_name="SDL_SCANCODE_V"]
    V
    @[link_name="SDL_SCANCODE_W"]
    W
    @[link_name="SDL_SCANCODE_X"]
    X
    @[link_name="SDL_SCANCODE_Y"]
    Y
    @[link_name="SDL_SCANCODE_Z"]
    Z
    @[link_name="SDL_SCANCODE_1"]
    SCANCODE_1
    @[link_name="SDL_SCANCODE_2"]
    SCANCODE_2
    @[link_name="SDL_SCANCODE_3"]
    SCANCODE_3
    @[link_name="SDL_SCANCODE_4"]
    SCANCODE_4
    @[link_name="SDL_SCANCODE_5"]
    SCANCODE_5
    @[link_name="SDL_SCANCODE_6"]
    SCANCODE_6
    @[link_name="SDL_SCANCODE_7"]
    SCANCODE_7
    @[link_name="SDL_SCANCODE_8"]
    SCANCODE_8
    @[link_name="SDL_SCANCODE_9"]
    SCANCODE_9
    @[link_name="SDL_SCANCODE_0"]
    SCANCODE_0
    @[link_name="SDL_SCANCODE_RETURN"]
    RETURN
    @[link_name="SDL_SCANCODE_ESCAPE"]
    ESCAPE
    @[link_name="SDL_SCANCODE_BACKSPACE"]
    BACKSPACE
    @[link_name="SDL_SCANCODE_TAB"]
    TAB
    @[link_name="SDL_SCANCODE_SPACE"]
    SPACE
    @[link_name="SDL_SCANCODE_MINUS"]
    MINUS
    @[link_name="SDL_SCANCODE_EQUALS"]
    EQUALS
    @[link_name="SDL_SCANCODE_LEFTBRACKET"]
    LEFTBRACKET
    @[link_name="SDL_SCANCODE_RIGHTBRACKET"]
    RIGHTBRACKET
    @[link_name="SDL_SCANCODE_BACKSLASH"]
    BACKSLASH
    @[link_name="SDL_SCANCODE_NONUSHASH"]
    NONUSHASH
    @[link_name="SDL_SCANCODE_SEMICOLON"]
    SEMICOLON
    @[link_name="SDL_SCANCODE_APOSTROPHE"]
    APOSTROPHE
    @[link_name="SDL_SCANCODE_GRAVE"]
    GRAVE
    @[link_name="SDL_SCANCODE_COMMA"]
    COMMA
    @[link_name="SDL_SCANCODE_PERIOD"]
    PERIOD
    @[link_name="SDL_SCANCODE_SLASH"]
    SLASH
    @[link_name="SDL_SCANCODE_CAPSLOCK"]
    CAPSLOCK
    @[link_name="SDL_SCANCODE_F1"]
    F1
    @[link_name="SDL_SCANCODE_F2"]
    F2
    @[link_name="SDL_SCANCODE_F3"]
    F3
    @[link_name="SDL_SCANCODE_F4"]
    F4
    @[link_name="SDL_SCANCODE_F5"]
    F5
    @[link_name="SDL_SCANCODE_F6"]
    F6
    @[link_name="SDL_SCANCODE_F7"]
    F7
    @[link_name="SDL_SCANCODE_F8"]
    F8
    @[link_name="SDL_SCANCODE_F9"]
    F9
    @[link_name="SDL_SCANCODE_F10"]
    F10
    @[link_name="SDL_SCANCODE_F11"]
    F11
    @[link_name="SDL_SCANCODE_F12"]
    F12
    @[link_name="SDL_SCANCODE_PRINTSCREEN"]
    PRINTSCREEN
    @[link_name="SDL_SCANCODE_SCROLLLOCK"]
    SCROLLLOCK
    @[link_name="SDL_SCANCODE_PAUSE"]
    PAUSE
    @[link_name="SDL_SCANCODE_INSERT"]
    INSERT
    @[link_name="SDL_SCANCODE_HOME"]
    HOME
    @[link_name="SDL_SCANCODE_PAGEUP"]
    PAGEUP
    @[link_name="SDL_SCANCODE_DELETE"]
    DELETE
    @[link_name="SDL_SCANCODE_END"]
    END
    @[link_name="SDL_SCANCODE_PAGEDOWN"]
    PAGEDOWN
    @[link_name="SDL_SCANCODE_RIGHT"]
    RIGHT
    @[link_name="SDL_SCANCODE_LEFT"]
    LEFT
    @[link_name="SDL_SCANCODE_DOWN"]
    DOWN
    @[link_name="SDL_SCANCODE_UP"]
    UP
    @[link_name="SDL_SCANCODE_NUMLOCKCLEAR"]
    NUMLOCKCLEAR
    @[link_name="SDL_SCANCODE_KP_DIVIDE"]
    KP_DIVIDE
    @[link_name="SDL_SCANCODE_KP_MULTIPLY"]
    KP_MULTIPLY
    @[link_name="SDL_SCANCODE_KP_MINUS"]
    KP_MINUS
    @[link_name="SDL_SCANCODE_KP_PLUS"]
    KP_PLUS
    @[link_name="SDL_SCANCODE_KP_ENTER"]
    KP_ENTER
    @[link_name="SDL_SCANCODE_KP_1"]
    KP_1
    @[link_name="SDL_SCANCODE_KP_2"]
    KP_2
    @[link_name="SDL_SCANCODE_KP_3"]
    KP_3
    @[link_name="SDL_SCANCODE_KP_4"]
    KP_4
    @[link_name="SDL_SCANCODE_KP_5"]
    KP_5
    @[link_name="SDL_SCANCODE_KP_6"]
    KP_6
    @[link_name="SDL_SCANCODE_KP_7"]
    KP_7
    @[link_name="SDL_SCANCODE_KP_8"]
    KP_8
    @[link_name="SDL_SCANCODE_KP_9"]
    KP_9
    @[link_name="SDL_SCANCODE_KP_0"]
    KP_0
    @[link_name="SDL_SCANCODE_KP_PERIOD"]
    KP_PERIOD
    @[link_name="SDL_SCANCODE_NONUSBACKSLASH"]
    NONUSBACKSLASH
    @[link_name="SDL_SCANCODE_APPLICATION"]
    APPLICATION
    @[link_name="SDL_SCANCODE_POWER"]
    POWER
    @[link_name="SDL_SCANCODE_KP_EQUALS"]
    KP_EQUALS
    @[link_name="SDL_SCANCODE_F13"]
    F13
    @[link_name="SDL_SCANCODE_F14"]
    F14
    @[link_name="SDL_SCANCODE_F15"]
    F15
    @[link_name="SDL_SCANCODE_F16"]
    F16
    @[link_name="SDL_SCANCODE_F17"]
    F17
    @[link_name="SDL_SCANCODE_F18"]
    F18
    @[link_name="SDL_SCANCODE_F19"]
    F19
    @[link_name="SDL_SCANCODE_F20"]
    F20
    @[link_name="SDL_SCANCODE_F21"]
    F21
    @[link_name="SDL_SCANCODE_F22"]
    F22
    @[link_name="SDL_SCANCODE_F23"]
    F23
    @[link_name="SDL_SCANCODE_F24"]
    F24
    @[link_name="SDL_SCANCODE_EXECUTE"]
    EXECUTE
    @[link_name="SDL_SCANCODE_HELP"]
    HELP
    @[link_name="SDL_SCANCODE_MENU"]
    MENU
    @[link_name="SDL_SCANCODE_SELECT"]
    SELECT
    @[link_name="SDL_SCANCODE_STOP"]
    STOP
    @[link_name="SDL_SCANCODE_AGAIN"]
    AGAIN
    @[link_name="SDL_SCANCODE_UNDO"]
    UNDO
    @[link_name="SDL_SCANCODE_CUT"]
    CUT
    @[link_name="SDL_SCANCODE_COPY"]
    COPY
    @[link_name="SDL_SCANCODE_PASTE"]
    PASTE
    @[link_name="SDL_SCANCODE_FIND"]
    FIND
    @[link_name="SDL_SCANCODE_MUTE"]
    MUTE
    @[link_name="SDL_SCANCODE_VOLUMEUP"]
    VOLUMEUP
    @[link_name="SDL_SCANCODE_VOLUMEDOWN"]
    VOLUMEDOWN
    @[link_name="SDL_SCANCODE_KP_COMMA"]
    KP_COMMA
    @[link_name="SDL_SCANCODE_KP_EQUALSAS400"]
    KP_EQUALSAS400
    @[link_name="SDL_SCANCODE_INTERNATIONAL1"]
    INTERNATIONAL1
    @[link_name="SDL_SCANCODE_INTERNATIONAL2"]
    INTERNATIONAL2
    @[link_name="SDL_SCANCODE_INTERNATIONAL3"]
    INTERNATIONAL3
    @[link_name="SDL_SCANCODE_INTERNATIONAL4"]
    INTERNATIONAL4
    @[link_name="SDL_SCANCODE_INTERNATIONAL5"]
    INTERNATIONAL5
    @[link_name="SDL_SCANCODE_INTERNATIONAL6"]
    INTERNATIONAL6
    @[link_name="SDL_SCANCODE_INTERNATIONAL7"]
    INTERNATIONAL7
    @[link_name="SDL_SCANCODE_INTERNATIONAL8"]
    INTERNATIONAL8
    @[link_name="SDL_SCANCODE_INTERNATIONAL9"]
    INTERNATIONAL9
    @[link_name="SDL_SCANCODE_LANG1"]
    LANG1
    @[link_name="SDL_SCANCODE_LANG2"]
    LANG2
    @[link_name="SDL_SCANCODE_LANG3"]
    LANG3
    @[link_name="SDL_SCANCODE_LANG4"]
    LANG4
    @[link_name="SDL_SCANCODE_LANG5"]
    LANG5
    @[link_name="SDL_SCANCODE_LANG6"]
    LANG6
    @[link_name="SDL_SCANCODE_LANG7"]
    LANG7
    @[link_name="SDL_SCANCODE_LANG8"]
    LANG8
    @[link_name="SDL_SCANCODE_LANG9"]
    LANG9
    @[link_name="SDL_SCANCODE_ALTERASE"]
    ALTERASE
    @[link_name="SDL_SCANCODE_SYSREQ"]
    SYSREQ
    @[link_name="SDL_SCANCODE_CANCEL"]
    CANCEL
    @[link_name="SDL_SCANCODE_CLEAR"]
    CLEAR
    @[link_name="SDL_SCANCODE_PRIOR"]
    PRIOR
    @[link_name="SDL_SCANCODE_RETURN2"]
    RETURN2
    @[link_name="SDL_SCANCODE_SEPARATOR"]
    SEPARATOR
    @[link_name="SDL_SCANCODE_OUT"]
    OUT
    @[link_name="SDL_SCANCODE_OPER"]
    OPER
    @[link_name="SDL_SCANCODE_CLEARAGAIN"]
    CLEARAGAIN
    @[link_name="SDL_SCANCODE_CRSEL"]
    CRSEL
    @[link_name="SDL_SCANCODE_EXSEL"]
    EXSEL
    @[link_name="SDL_SCANCODE_KP_00"]
    KP_00
    @[link_name="SDL_SCANCODE_KP_000"]
    KP_000
    @[link_name="SDL_SCANCODE_THOUSANDSSEPARATOR"]
    THOUSANDSSEPARATOR
    @[link_name="SDL_SCANCODE_DECIMALSEPARATOR"]
    DECIMALSEPARATOR
    @[link_name="SDL_SCANCODE_CURRENCYUNIT"]
    CURRENCYUNIT
    @[link_name="SDL_SCANCODE_CURRENCYSUBUNIT"]
    CURRENCYSUBUNIT
    @[link_name="SDL_SCANCODE_KP_LEFTPAREN"]
    KP_LEFTPAREN
    @[link_name="SDL_SCANCODE_KP_RIGHTPAREN"]
    KP_RIGHTPAREN
    @[link_name="SDL_SCANCODE_KP_LEFTBRACE"]
    KP_LEFTBRACE
    @[link_name="SDL_SCANCODE_KP_RIGHTBRACE"]
    KP_RIGHTBRACE
    @[link_name="SDL_SCANCODE_KP_TAB"]
    KP_TAB
    @[link_name="SDL_SCANCODE_KP_BACKSPACE"]
    KP_BACKSPACE
    @[link_name="SDL_SCANCODE_KP_A"]
    KP_A
    @[link_name="SDL_SCANCODE_KP_B"]
    KP_B
    @[link_name="SDL_SCANCODE_KP_C"]
    KP_C
    @[link_name="SDL_SCANCODE_KP_D"]
    KP_D
    @[link_name="SDL_SCANCODE_KP_E"]
    KP_E
    @[link_name="SDL_SCANCODE_KP_F"]
    KP_F
    @[link_name="SDL_SCANCODE_KP_XOR"]
    KP_XOR
    @[link_name="SDL_SCANCODE_KP_POWER"]
    KP_POWER
    @[link_name="SDL_SCANCODE_KP_PERCENT"]
    KP_PERCENT
    @[link_name="SDL_SCANCODE_KP_LESS"]
    KP_LESS
    @[link_name="SDL_SCANCODE_KP_GREATER"]
    KP_GREATER
    @[link_name="SDL_SCANCODE_KP_AMPERSAND"]
    KP_AMPERSAND
    @[link_name="SDL_SCANCODE_KP_DBLAMPERSAND"]
    KP_DBLAMPERSAND
    @[link_name="SDL_SCANCODE_KP_VERTICALBAR"]
    KP_VERTICALBAR
    @[link_name="SDL_SCANCODE_KP_DBLVERTICALBAR"]
    KP_DBLVERTICALBAR
    @[link_name="SDL_SCANCODE_KP_COLON"]
    KP_COLON
    @[link_name="SDL_SCANCODE_KP_HASH"]
    KP_HASH
    @[link_name="SDL_SCANCODE_KP_SPACE"]
    KP_SPACE
    @[link_name="SDL_SCANCODE_KP_AT"]
    KP_AT
    @[link_name="SDL_SCANCODE_KP_EXCLAM"]
    KP_EXCLAM
    @[link_name="SDL_SCANCODE_KP_MEMSTORE"]
    KP_MEMSTORE
    @[link_name="SDL_SCANCODE_KP_MEMRECALL"]
    KP_MEMRECALL
    @[link_name="SDL_SCANCODE_KP_MEMCLEAR"]
    KP_MEMCLEAR
    @[link_name="SDL_SCANCODE_KP_MEMADD"]
    KP_MEMADD
    @[link_name="SDL_SCANCODE_KP_MEMSUBTRACT"]
    KP_MEMSUBTRACT
    @[link_name="SDL_SCANCODE_KP_MEMMULTIPLY"]
    KP_MEMMULTIPLY
    @[link_name="SDL_SCANCODE_KP_MEMDIVIDE"]
    KP_MEMDIVIDE
    @[link_name="SDL_SCANCODE_KP_PLUSMINUS"]
    KP_PLUSMINUS
    @[link_name="SDL_SCANCODE_KP_CLEAR"]
    KP_CLEAR
    @[link_name="SDL_SCANCODE_KP_CLEARENTRY"]
    KP_CLEARENTRY
    @[link_name="SDL_SCANCODE_KP_BINARY"]
    KP_BINARY
    @[link_name="SDL_SCANCODE_KP_OCTAL"]
    KP_OCTAL
    @[link_name="SDL_SCANCODE_KP_DECIMAL"]
    KP_DECIMAL
    @[link_name="SDL_SCANCODE_KP_HEXADECIMAL"]
    KP_HEXADECIMAL
    @[link_name="SDL_SCANCODE_LCTRL"]
    LCTRL
    @[link_name="SDL_SCANCODE_LSHIFT"]
    LSHIFT
    @[link_name="SDL_SCANCODE_LALT"]
    LALT
    @[link_name="SDL_SCANCODE_LGUI"]
    LGUI
    @[link_name="SDL_SCANCODE_RCTRL"]
    RCTRL
    @[link_name="SDL_SCANCODE_RSHIFT"]
    RSHIFT
    @[link_name="SDL_SCANCODE_RALT"]
    RALT
    @[link_name="SDL_SCANCODE_RGUI"]
    RGUI
    @[link_name="SDL_SCANCODE_MODE"]
    MODE
    @[link_name="SDL_SCANCODE_SLEEP"]
    SLEEP
    @[link_name="SDL_SCANCODE_WAKE"]
    WAKE
    @[link_name="SDL_SCANCODE_CHANNEL_INCREMENT"]
    CHANNEL_INCREMENT
    @[link_name="SDL_SCANCODE_CHANNEL_DECREMENT"]
    CHANNEL_DECREMENT
    @[link_name="SDL_SCANCODE_MEDIA_PLAY"]
    MEDIA_PLAY
    @[link_name="SDL_SCANCODE_MEDIA_PAUSE"]
    MEDIA_PAUSE
    @[link_name="SDL_SCANCODE_MEDIA_RECORD"]
    MEDIA_RECORD
    @[link_name="SDL_SCANCODE_MEDIA_FAST_FORWARD"]
    MEDIA_FAST_FORWARD
    @[link_name="SDL_SCANCODE_MEDIA_REWIND"]
    MEDIA_REWIND
    @[link_name="SDL_SCANCODE_MEDIA_NEXT_TRACK"]
    MEDIA_NEXT_TRACK
    @[link_name="SDL_SCANCODE_MEDIA_PREVIOUS_TRACK"]
    MEDIA_PREVIOUS_TRACK
    @[link_name="SDL_SCANCODE_MEDIA_STOP"]
    MEDIA_STOP
    @[link_name="SDL_SCANCODE_MEDIA_EJECT"]
    MEDIA_EJECT
    @[link_name="SDL_SCANCODE_MEDIA_PLAY_PAUSE"]
    MEDIA_PLAY_PAUSE
    @[link_name="SDL_SCANCODE_MEDIA_SELECT"]
    MEDIA_SELECT
    @[link_name="SDL_SCANCODE_AC_NEW"]
    AC_NEW
    @[link_name="SDL_SCANCODE_AC_OPEN"]
    AC_OPEN
    @[link_name="SDL_SCANCODE_AC_CLOSE"]
    AC_CLOSE
    @[link_name="SDL_SCANCODE_AC_EXIT"]
    AC_EXIT
    @[link_name="SDL_SCANCODE_AC_SAVE"]
    AC_SAVE
    @[link_name="SDL_SCANCODE_AC_PRINT"]
    AC_PRINT
    @[link_name="SDL_SCANCODE_AC_PROPERTIES"]
    AC_PROPERTIES
    @[link_name="SDL_SCANCODE_AC_SEARCH"]
    AC_SEARCH
    @[link_name="SDL_SCANCODE_AC_HOME"]
    AC_HOME
    @[link_name="SDL_SCANCODE_AC_BACK"]
    AC_BACK
    @[link_name="SDL_SCANCODE_AC_FORWARD"]
    AC_FORWARD
    @[link_name="SDL_SCANCODE_AC_STOP"]
    AC_STOP
    @[link_name="SDL_SCANCODE_AC_REFRESH"]
    AC_REFRESH
    @[link_name="SDL_SCANCODE_AC_BOOKMARKS"]
    AC_BOOKMARKS
    @[link_name="SDL_SCANCODE_SOFTLEFT"]
    SOFTLEFT
    @[link_name="SDL_SCANCODE_SOFTRIGHT"]
    SOFTRIGHT
    @[link_name="SDL_SCANCODE_CALL"]
    CALL
    @[link_name="SDL_SCANCODE_ENDCALL"]
    ENDCALL
    @[link_name="SDL_SCANCODE_RESERVED"]
    RESERVED
    @[link_name="SDL_SCANCODE_COUNT"]
    COUNT
  }




  @[link_name="SDL_Cursor"]
  opaque Cursor

  @[link_name="SDL_SystemCursor"]
  enum SystemCursor {
    @[link_name="SDL_SYSTEM_CURSOR_DEFAULT"]
    DEFAULT
    @[link_name="SDL_SYSTEM_CURSOR_TEXT"]
    TEXT
    @[link_name="SDL_SYSTEM_CURSOR_WAIT"]
    WAIT
    @[link_name="SDL_SYSTEM_CURSOR_CROSSHAIR"]
    CROSSHAIR
    @[link_name="SDL_SYSTEM_CURSOR_PROGRESS"]
    PROGRESS
    @[link_name="SDL_SYSTEM_CURSOR_NWSE_RESIZE"]
    NWSE_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_NESW_RESIZE"]
    NESW_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_EW_RESIZE"]
    EW_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_NS_RESIZE"]
    NS_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_MOVE"]
    MOVE
    @[link_name="SDL_SYSTEM_CURSOR_NOT_ALLOWED"]
    NOT_ALLOWED
    @[link_name="SDL_SYSTEM_CURSOR_POINTER"]
    POINTER
    @[link_name="SDL_SYSTEM_CURSOR_NW_RESIZE"]
    NW_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_N_RESIZE"]
    N_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_NE_RESIZE"]
    NE_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_E_RESIZE"]
    E_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_SE_RESIZE"]
    SE_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_S_RESIZE"]
    S_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_SW_RESIZE"]
    SW_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_W_RESIZE"]
    W_RESIZE
    @[link_name="SDL_SYSTEM_CURSOR_COUNT"]
    COUNT
  }


  @[link_name="SDL_TouchDeviceType"]
  enum TouchDeviceType {
    @[link_name="SDL_TOUCH_DEVICE_INVALID"]
    INVALID
    @[link_name="SDL_TOUCH_DEVICE_DIRECT"]
    DIRECT
    @[link_name="SDL_TOUCH_DEVICE_INDIRECT_ABSOLUTE"]
    INDIRECT_ABSOLUTE
    @[link_name="SDL_TOUCH_DEVICE_INDIRECT_RELATIVE"]
    INDIRECT_RELATIVE
  }


  @[link_name="SDL_Finger"]
  struct Finger {
    id FingerID
    x libc::c_float
    y libc::c_float
    pressure libc::c_float
  }

  @[link_name="SDL_CommonEvent"]
  struct CommonEvent {
    type u32
    reserved u32
    timestamp u64
  }

  @[link_name="SDL_EventType"]
  enum EventType {
    @[link_name="SDL_EVENT_FIRST"]
    FIRST
    @[link_name="SDL_EVENT_QUIT"]
    QUIT
    @[link_name="SDL_EVENT_TERMINATING"]
    TERMINATING
    @[link_name="SDL_EVENT_LOW_MEMORY"]
    LOW_MEMORY
    @[link_name="SDL_EVENT_WILL_ENTER_BACKGROUND"]
    WILL_ENTER_BACKGROUND
    @[link_name="SDL_EVENT_DID_ENTER_BACKGROUND"]
    DID_ENTER_BACKGROUND
    @[link_name="SDL_EVENT_WILL_ENTER_FOREGROUND"]
    WILL_ENTER_FOREGROUND
    @[link_name="SDL_EVENT_DID_ENTER_FOREGROUND"]
    DID_ENTER_FOREGROUND
    @[link_name="SDL_EVENT_LOCALE_CHANGED"]
    LOCALE_CHANGED
    @[link_name="SDL_EVENT_SYSTEM_THEME_CHANGED"]
    SYSTEM_THEME_CHANGED
    @[link_name="SDL_EVENT_DISPLAY_ORIENTATION"]
    DISPLAY_ORIENTATION
    @[link_name="SDL_EVENT_DISPLAY_ADDED"]
    DISPLAY_ADDED
    @[link_name="SDL_EVENT_DISPLAY_REMOVED"]
    DISPLAY_REMOVED
    @[link_name="SDL_EVENT_DISPLAY_MOVED"]
    DISPLAY_MOVED
    @[link_name="SDL_EVENT_DISPLAY_DESKTOP_MODE_CHANGED"]
    DISPLAY_DESKTOP_MODE_CHANGED
    @[link_name="SDL_EVENT_DISPLAY_CURRENT_MODE_CHANGED"]
    DISPLAY_CURRENT_MODE_CHANGED
    @[link_name="SDL_EVENT_DISPLAY_CONTENT_SCALE_CHANGED"]
    DISPLAY_CONTENT_SCALE_CHANGED
    @[link_name="SDL_EVENT_DISPLAY_FIRST"]
    DISPLAY_FIRST
    @[link_name="SDL_EVENT_DISPLAY_LAST"]
    DISPLAY_LAST
    @[link_name="SDL_EVENT_WINDOW_SHOWN"]
    WINDOW_SHOWN
    @[link_name="SDL_EVENT_WINDOW_HIDDEN"]
    WINDOW_HIDDEN
    @[link_name="SDL_EVENT_WINDOW_EXPOSED"]
    WINDOW_EXPOSED
    @[link_name="SDL_EVENT_WINDOW_MOVED"]
    WINDOW_MOVED
    @[link_name="SDL_EVENT_WINDOW_RESIZED"]
    WINDOW_RESIZED
    @[link_name="SDL_EVENT_WINDOW_PIXEL_SIZE_CHANGED"]
    WINDOW_PIXEL_SIZE_CHANGED
    @[link_name="SDL_EVENT_WINDOW_METAL_VIEW_RESIZED"]
    WINDOW_METAL_VIEW_RESIZED
    @[link_name="SDL_EVENT_WINDOW_MINIMIZED"]
    WINDOW_MINIMIZED
    @[link_name="SDL_EVENT_WINDOW_MAXIMIZED"]
    WINDOW_MAXIMIZED
    @[link_name="SDL_EVENT_WINDOW_RESTORED"]
    WINDOW_RESTORED
    @[link_name="SDL_EVENT_WINDOW_MOUSE_ENTER"]
    WINDOW_MOUSE_ENTER
    @[link_name="SDL_EVENT_WINDOW_MOUSE_LEAVE"]
    WINDOW_MOUSE_LEAVE
    @[link_name="SDL_EVENT_WINDOW_FOCUS_GAINED"]
    WINDOW_FOCUS_GAINED
    @[link_name="SDL_EVENT_WINDOW_FOCUS_LOST"]
    WINDOW_FOCUS_LOST
    @[link_name="SDL_EVENT_WINDOW_CLOSE_REQUESTED"]
    WINDOW_CLOSE_REQUESTED
    @[link_name="SDL_EVENT_WINDOW_HIT_TEST"]
    WINDOW_HIT_TEST
    @[link_name="SDL_EVENT_WINDOW_ICCPROF_CHANGED"]
    WINDOW_ICCPROF_CHANGED
    @[link_name="SDL_EVENT_WINDOW_DISPLAY_CHANGED"]
    WINDOW_DISPLAY_CHANGED
    @[link_name="SDL_EVENT_WINDOW_DISPLAY_SCALE_CHANGED"]
    WINDOW_DISPLAY_SCALE_CHANGED
    @[link_name="SDL_EVENT_WINDOW_SAFE_AREA_CHANGED"]
    WINDOW_SAFE_AREA_CHANGED
    @[link_name="SDL_EVENT_WINDOW_OCCLUDED"]
    WINDOW_OCCLUDED
    @[link_name="SDL_EVENT_WINDOW_ENTER_FULLSCREEN"]
    WINDOW_ENTER_FULLSCREEN
    @[link_name="SDL_EVENT_WINDOW_LEAVE_FULLSCREEN"]
    WINDOW_LEAVE_FULLSCREEN
    @[link_name="SDL_EVENT_WINDOW_DESTROYED"]
    WINDOW_DESTROYED
    @[link_name="SDL_EVENT_WINDOW_HDR_STATE_CHANGED"]
    WINDOW_HDR_STATE_CHANGED
    @[link_name="SDL_EVENT_WINDOW_FIRST"]
    WINDOW_FIRST
    @[link_name="SDL_EVENT_WINDOW_LAST"]
    WINDOW_LAST
    @[link_name="SDL_EVENT_KEY_DOWN"]
    KEY_DOWN
    @[link_name="SDL_EVENT_KEY_UP"]
    KEY_UP
    @[link_name="SDL_EVENT_TEXT_EDITING"]
    TEXT_EDITING
    @[link_name="SDL_EVENT_TEXT_INPUT"]
    TEXT_INPUT
    @[link_name="SDL_EVENT_KEYMAP_CHANGED"]
    KEYMAP_CHANGED
    @[link_name="SDL_EVENT_KEYBOARD_ADDED"]
    KEYBOARD_ADDED
    @[link_name="SDL_EVENT_KEYBOARD_REMOVED"]
    KEYBOARD_REMOVED
    @[link_name="SDL_EVENT_TEXT_EDITING_CANDIDATES"]
    TEXT_EDITING_CANDIDATES
    @[link_name="SDL_EVENT_MOUSE_MOTION"]
    MOUSE_MOTION
    @[link_name="SDL_EVENT_MOUSE_BUTTON_DOWN"]
    MOUSE_BUTTON_DOWN
    @[link_name="SDL_EVENT_MOUSE_BUTTON_UP"]
    MOUSE_BUTTON_UP
    @[link_name="SDL_EVENT_MOUSE_WHEEL"]
    MOUSE_WHEEL
    @[link_name="SDL_EVENT_MOUSE_ADDED"]
    MOUSE_ADDED
    @[link_name="SDL_EVENT_MOUSE_REMOVED"]
    MOUSE_REMOVED
    @[link_name="SDL_EVENT_JOYSTICK_AXIS_MOTION"]
    JOYSTICK_AXIS_MOTION
    @[link_name="SDL_EVENT_JOYSTICK_BALL_MOTION"]
    JOYSTICK_BALL_MOTION
    @[link_name="SDL_EVENT_JOYSTICK_HAT_MOTION"]
    JOYSTICK_HAT_MOTION
    @[link_name="SDL_EVENT_JOYSTICK_BUTTON_DOWN"]
    JOYSTICK_BUTTON_DOWN
    @[link_name="SDL_EVENT_JOYSTICK_BUTTON_UP"]
    JOYSTICK_BUTTON_UP
    @[link_name="SDL_EVENT_JOYSTICK_ADDED"]
    JOYSTICK_ADDED
    @[link_name="SDL_EVENT_JOYSTICK_REMOVED"]
    JOYSTICK_REMOVED
    @[link_name="SDL_EVENT_JOYSTICK_BATTERY_UPDATED"]
    JOYSTICK_BATTERY_UPDATED
    @[link_name="SDL_EVENT_JOYSTICK_UPDATE_COMPLETE"]
    JOYSTICK_UPDATE_COMPLETE
    @[link_name="SDL_EVENT_GAMEPAD_AXIS_MOTION"]
    GAMEPAD_AXIS_MOTION
    @[link_name="SDL_EVENT_GAMEPAD_BUTTON_DOWN"]
    GAMEPAD_BUTTON_DOWN
    @[link_name="SDL_EVENT_GAMEPAD_BUTTON_UP"]
    GAMEPAD_BUTTON_UP
    @[link_name="SDL_EVENT_GAMEPAD_ADDED"]
    GAMEPAD_ADDED
    @[link_name="SDL_EVENT_GAMEPAD_REMOVED"]
    GAMEPAD_REMOVED
    @[link_name="SDL_EVENT_GAMEPAD_REMAPPED"]
    GAMEPAD_REMAPPED
    @[link_name="SDL_EVENT_GAMEPAD_TOUCHPAD_DOWN"]
    GAMEPAD_TOUCHPAD_DOWN
    @[link_name="SDL_EVENT_GAMEPAD_TOUCHPAD_MOTION"]
    GAMEPAD_TOUCHPAD_MOTION
    @[link_name="SDL_EVENT_GAMEPAD_TOUCHPAD_UP"]
    GAMEPAD_TOUCHPAD_UP
    @[link_name="SDL_EVENT_GAMEPAD_SENSOR_UPDATE"]
    GAMEPAD_SENSOR_UPDATE
    @[link_name="SDL_EVENT_GAMEPAD_UPDATE_COMPLETE"]
    GAMEPAD_UPDATE_COMPLETE
    @[link_name="SDL_EVENT_GAMEPAD_STEAM_HANDLE_UPDATED"]
    GAMEPAD_STEAM_HANDLE_UPDATED
    @[link_name="SDL_EVENT_FINGER_DOWN"]
    FINGER_DOWN
    @[link_name="SDL_EVENT_FINGER_UP"]
    FINGER_UP
    @[link_name="SDL_EVENT_FINGER_MOTION"]
    FINGER_MOTION
    @[link_name="SDL_EVENT_FINGER_CANCELED"]
    FINGER_CANCELED
    @[link_name="SDL_EVENT_CLIPBOARD_UPDATE"]
    CLIPBOARD_UPDATE
    @[link_name="SDL_EVENT_DROP_FILE"]
    DROP_FILE
    @[link_name="SDL_EVENT_DROP_TEXT"]
    DROP_TEXT
    @[link_name="SDL_EVENT_DROP_BEGIN"]
    DROP_BEGIN
    @[link_name="SDL_EVENT_DROP_COMPLETE"]
    DROP_COMPLETE
    @[link_name="SDL_EVENT_DROP_POSITION"]
    DROP_POSITION
    @[link_name="SDL_EVENT_AUDIO_DEVICE_ADDED"]
    AUDIO_DEVICE_ADDED
    @[link_name="SDL_EVENT_AUDIO_DEVICE_REMOVED"]
    AUDIO_DEVICE_REMOVED
    @[link_name="SDL_EVENT_AUDIO_DEVICE_FORMAT_CHANGED"]
    AUDIO_DEVICE_FORMAT_CHANGED
    @[link_name="SDL_EVENT_SENSOR_UPDATE"]
    SENSOR_UPDATE
    @[link_name="SDL_EVENT_PEN_PROXIMITY_IN"]
    PEN_PROXIMITY_IN
    @[link_name="SDL_EVENT_PEN_PROXIMITY_OUT"]
    PEN_PROXIMITY_OUT
    @[link_name="SDL_EVENT_PEN_DOWN"]
    PEN_DOWN
    @[link_name="SDL_EVENT_PEN_UP"]
    PEN_UP
    @[link_name="SDL_EVENT_PEN_BUTTON_DOWN"]
    PEN_BUTTON_DOWN
    @[link_name="SDL_EVENT_PEN_BUTTON_UP"]
    PEN_BUTTON_UP
    @[link_name="SDL_EVENT_PEN_MOTION"]
    PEN_MOTION
    @[link_name="SDL_EVENT_PEN_AXIS"]
    PEN_AXIS
    @[link_name="SDL_EVENT_CAMERA_DEVICE_ADDED"]
    CAMERA_DEVICE_ADDED
    @[link_name="SDL_EVENT_CAMERA_DEVICE_REMOVED"]
    CAMERA_DEVICE_REMOVED
    @[link_name="SDL_EVENT_CAMERA_DEVICE_APPROVED"]
    CAMERA_DEVICE_APPROVED
    @[link_name="SDL_EVENT_CAMERA_DEVICE_DENIED"]
    CAMERA_DEVICE_DENIED
    @[link_name="SDL_EVENT_RENDER_TARGETS_RESET"]
    RENDER_TARGETS_RESET
    @[link_name="SDL_EVENT_RENDER_DEVICE_RESET"]
    RENDER_DEVICE_RESET
    @[link_name="SDL_EVENT_RENDER_DEVICE_LOST"]
    RENDER_DEVICE_LOST
    @[link_name="SDL_EVENT_PRIVATE0"]
    PRIVATE0
    @[link_name="SDL_EVENT_PRIVATE1"]
    PRIVATE1
    @[link_name="SDL_EVENT_PRIVATE2"]
    PRIVATE2
    @[link_name="SDL_EVENT_PRIVATE3"]
    PRIVATE3
    @[link_name="SDL_EVENT_POLL_SENTINEL"]
    POLL_SENTINEL
    @[link_name="SDL_EVENT_USER"]
    USER
    @[link_name="SDL_EVENT_LAST"]
    LAST
    @[link_name="SDL_EVENT_ENUM_PADDING"]
    ENUM_PADDING
  }

  @[link_name="SDL_DisplayEvent"]
  struct DisplayEvent {
    type EventType
    reserved u32
    timestamp u64
    displayID DisplayID
    data1 i32
    data2 i32
  }

  @[link_name="SDL_WindowEvent"]
  struct WindowEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    data1 i32
    data2 i32
  }

  @[link_name="SDL_KeyboardDeviceEvent"]
  struct KeyboardDeviceEvent {
    type EventType
    reserved u32
    timestamp u64
    which KeyboardID
  }

  @[link_name="SDL_KeyboardEvent"]
  struct KeyboardEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    which KeyboardID
    scancode Scancode
    key Keycode
    mod Keymod
    raw u16
    down bool
    repeat bool
  }

  @[link_name="SDL_TextEditingEvent"]
  struct TextEditingEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    text const cstring
    start i32
    length i32
  }

  @[link_name="SDL_TextEditingCandidatesEvent"]
  struct TextEditingCandidatesEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    candidates const *const cstring
    num_candidates i32
    selected_candidate i32
    horizontal bool
    padding1 u8
    padding2 u8
    padding3 u8
  }

  @[link_name="SDL_TextInputEvent"]
  struct TextInputEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    text const cstring
  }

  @[link_name="SDL_MouseDeviceEvent"]
  struct MouseDeviceEvent {
    type EventType
    reserved u32
    timestamp u64
    which MouseID
  }

  @[link_name="SDL_MouseMotionEvent"]
  struct MouseMotionEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    which MouseID
    state MouseButtonFlags
    x libc::c_float
    y libc::c_float
    xrel libc::c_float
    yrel libc::c_float
  }

  @[link_name="SDL_MouseButtonEvent"]
  struct MouseButtonEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    which MouseID
    button u8
    down bool
    clicks u8
    padding u8
    x libc::c_float
    y libc::c_float
  }

  @[link_name="SDL_MouseWheelDirection"]
  enum MouseWheelDirection {
    @[link_name="SDL_MOUSEWHEEL_NORMAL"]
    NORMAL
    @[link_name="SDL_MOUSEWHEEL_FLIPPED"]
    FLIPPED
  }

  @[link_name="SDL_MouseWheelEvent"]
  struct MouseWheelEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    which MouseID
    x libc::c_float
    y libc::c_float
    direction MouseWheelDirection
    mouse_x libc::c_float
    mouse_y libc::c_float
    integer_x i32
    integer_y i32
  }

  @[link_name="SDL_JoyDeviceEvent"]
  struct JoyDeviceEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
  }

  @[link_name="SDL_JoyAxisEvent"]
  struct JoyAxisEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
    axis u8
    padding1 u8
    padding2 u8
    padding3 u8
    value i16
    padding4 u16
  }

  @[link_name="SDL_JoyBallEvent"]
  struct JoyBallEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
    ball u8
    padding1 u8
    padding2 u8
    padding3 u8
    xrel i16
    yrel i16
  }

  @[link_name="SDL_JoyHatEvent"]
  struct JoyHatEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
    hat u8
    value u8
    padding1 u8
    padding2 u8
  }

  @[link_name="SDL_JoyButtonEvent"]
  struct JoyButtonEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
    button u8
    down bool
    padding1 u8
    padding2 u8
  }

  @[link_name="SDL_JoyBatteryEvent"]
  struct JoyBatteryEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
    state PowerState
    percent libc::c_int
  }

  @[link_name="SDL_GamepadDeviceEvent"]
  struct GamepadDeviceEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
  }

  @[link_name="SDL_GamepadAxisEvent"]
  struct GamepadAxisEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
    axis u8
    padding1 u8
    padding2 u8
    padding3 u8
    value i16
    padding4 u16
  }

  @[link_name="SDL_GamepadButtonEvent"]
  struct GamepadButtonEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
    button u8
    down bool
    padding1 u8
    padding2 u8
  }

  @[link_name="SDL_GamepadTouchpadEvent"]
  struct GamepadTouchpadEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
    touchpad i32
    finger i32
    x libc::c_float
    y libc::c_float
    pressure libc::c_float
  }

  @[link_name="SDL_GamepadSensorEvent"]
  struct GamepadSensorEvent {
    type EventType
    reserved u32
    timestamp u64
    which JoystickID
    sensor i32
    data libc::c_float[3]
    sensor_timestamp u64
  }

  @[link_name="SDL_AudioDeviceEvent"]
  struct AudioDeviceEvent {
    type EventType
    reserved u32
    timestamp u64
    which AudioDeviceID
    recording bool
    padding1 u8
    padding2 u8
    padding3 u8
  }

  @[link_name="SDL_CameraDeviceEvent"]
  struct CameraDeviceEvent {
    type EventType
    reserved u32
    timestamp u64
    which CameraID
  }

  @[link_name="SDL_SensorEvent"]
  struct SensorEvent {
    type EventType
    reserved u32
    timestamp u64
    which SensorID
    data libc::c_float[6]
    sensor_timestamp u64
  }

  @[link_name="SDL_QuitEvent"]
  struct QuitEvent {
    type EventType
    reserved u32
    timestamp u64
  }

  @[link_name="SDL_UserEvent"]
  struct UserEvent {
    type u32
    reserved u32
    timestamp u64
    windowID WindowID
    code i32
    data1 rawptr
    data2 rawptr
  }

  @[link_name="SDL_TouchFingerEvent"]
  struct TouchFingerEvent {
    type EventType
    reserved u32
    timestamp u64
    touchID TouchID
    fingerID FingerID
    x libc::c_float
    y libc::c_float
    dx libc::c_float
    dy libc::c_float
    pressure libc::c_float
    windowID WindowID
  }


  @[link_name="SDL_PenProximityEvent"]
  struct PenProximityEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    which PenID
  }


  @[link_name="SDL_PenTouchEvent"]
  struct PenTouchEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    which PenID
    pen_state PenInputFlags
    x libc::c_float
    y libc::c_float
    eraser bool
    down bool
  }

  @[link_name="SDL_PenMotionEvent"]
  struct PenMotionEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    which PenID
    pen_state PenInputFlags
    x libc::c_float
    y libc::c_float
  }

  @[link_name="SDL_PenButtonEvent"]
  struct PenButtonEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    which PenID
    pen_state PenInputFlags
    x libc::c_float
    y libc::c_float
    button u8
    down bool
  }

  @[link_name="SDL_PenAxis"]
  enum PenAxis {
    @[link_name="SDL_PEN_AXIS_PRESSURE"]
    PRESSURE
    @[link_name="SDL_PEN_AXIS_XTILT"]
    XTILT
    @[link_name="SDL_PEN_AXIS_YTILT"]
    YTILT
    @[link_name="SDL_PEN_AXIS_DISTANCE"]
    DISTANCE
    @[link_name="SDL_PEN_AXIS_ROTATION"]
    ROTATION
    @[link_name="SDL_PEN_AXIS_SLIDER"]
    SLIDER
    @[link_name="SDL_PEN_AXIS_TANGENTIAL_PRESSURE"]
    TANGENTIAL_PRESSURE
    @[link_name="SDL_PEN_AXIS_COUNT"]
    COUNT
  }

  @[link_name="SDL_PenAxisEvent"]
  struct PenAxisEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    which PenID
    pen_state PenInputFlags
    x libc::c_float
    y libc::c_float
    axis PenAxis
    value libc::c_float
  }

  @[link_name="SDL_RenderEvent"]
  struct RenderEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
  }

  @[link_name="SDL_DropEvent"]
  struct DropEvent {
    type EventType
    reserved u32
    timestamp u64
    windowID WindowID
    x libc::c_float
    y libc::c_float
    source const cstring
    data const cstring
  }

  @[link_name="SDL_ClipboardEvent"]
  struct ClipboardEvent {
    type EventType
    reserved u32
    timestamp u64
    owner bool
    num_mime_types i32
    mime_types *const cstring
  }

  @[link_name="SDL_Event"]
  union Event {
    type u32
    common CommonEvent
    display DisplayEvent
    window WindowEvent
    kdevice KeyboardDeviceEvent
    key KeyboardEvent
    edit TextEditingEvent
    edit_candidates TextEditingCandidatesEvent
    text TextInputEvent
    mdevice MouseDeviceEvent
    motion MouseMotionEvent
    button MouseButtonEvent
    wheel MouseWheelEvent
    jdevice JoyDeviceEvent
    jaxis JoyAxisEvent
    jball JoyBallEvent
    jhat JoyHatEvent
    jbutton JoyButtonEvent
    jbattery JoyBatteryEvent
    gdevice GamepadDeviceEvent
    gaxis GamepadAxisEvent
    gbutton GamepadButtonEvent
    gtouchpad GamepadTouchpadEvent
    gsensor GamepadSensorEvent
    adevice AudioDeviceEvent
    cdevice CameraDeviceEvent
    sensor SensorEvent
    quit QuitEvent
    user UserEvent
    tfinger TouchFingerEvent
    pproximity PenProximityEvent
    ptouch PenTouchEvent
    pmotion PenMotionEvent
    pbutton PenButtonEvent
    paxis PenAxisEvent
    render RenderEvent
    drop DropEvent
    clipboard ClipboardEvent
    padding u8[128]
  }

  @[link_name="SDL_EventAction"]
  enum EventAction {
    @[link_name="SDL_ADDEVENT"]
    ADDEVENT
    @[link_name="SDL_PEEKEVENT"]
    PEEKEVENT
    @[link_name="SDL_GETEVENT"]
    GETEVENT
  }


  @[link_name="SDL_Folder"]
  enum Folder {
    @[link_name="SDL_FOLDER_HOME"]
    HOME
    @[link_name="SDL_FOLDER_DESKTOP"]
    DESKTOP
    @[link_name="SDL_FOLDER_DOCUMENTS"]
    DOCUMENTS
    @[link_name="SDL_FOLDER_DOWNLOADS"]
    DOWNLOADS
    @[link_name="SDL_FOLDER_MUSIC"]
    MUSIC
    @[link_name="SDL_FOLDER_PICTURES"]
    PICTURES
    @[link_name="SDL_FOLDER_PUBLICSHARE"]
    PUBLICSHARE
    @[link_name="SDL_FOLDER_SAVEDGAMES"]
    SAVEDGAMES
    @[link_name="SDL_FOLDER_SCREENSHOTS"]
    SCREENSHOTS
    @[link_name="SDL_FOLDER_TEMPLATES"]
    TEMPLATES
    @[link_name="SDL_FOLDER_VIDEOS"]
    VIDEOS
    @[link_name="SDL_FOLDER_COUNT"]
    COUNT
  }

  @[link_name="SDL_EnumerationResult"]
  enum EnumerationResult {
    @[link_name="SDL_ENUM_CONTINUE"]
    CONTINUE
    @[link_name="SDL_ENUM_SUCCESS"]
    SUCCESS
    @[link_name="SDL_ENUM_FAILURE"]
    FAILURE
  }


  @[link_name="SDL_PathType"]
  enum PathType {
    @[link_name="SDL_PATHTYPE_NONE"]
    NONE
    @[link_name="SDL_PATHTYPE_FILE"]
    FILE
    @[link_name="SDL_PATHTYPE_DIRECTORY"]
    DIRECTORY
    @[link_name="SDL_PATHTYPE_OTHER"]
    OTHER
  }


  @[link_name="SDL_PathInfo"]
  struct PathInfo {
    type PathType
    size u64
    create_time Time
    modify_time Time
    access_time Time
  }



  @[link_name="SDL_GPUDevice"]
  opaque GPUDevice

  @[link_name="SDL_GPUComputePipelineCreateInfo"]
  struct GPUComputePipelineCreateInfo {
    code_size libc::c_size_t
    code const *u8
    entrypoint const cstring
    format GPUShaderFormat
    num_samplers u32
    num_readonly_storage_textures u32
    num_readonly_storage_buffers u32
    num_readwrite_storage_textures u32
    num_readwrite_storage_buffers u32
    num_uniform_buffers u32
    threadcount_x u32
    threadcount_y u32
    threadcount_z u32
    props PropertiesID
  }

  @[link_name="SDL_GPUComputePipeline"]
  opaque GPUComputePipeline

  @[link_name="SDL_GPUShader"]
  opaque GPUShader

  @[link_name="SDL_GPUVertexInputRate"]
  enum GPUVertexInputRate {
    @[link_name="SDL_GPU_VERTEXINPUTRATE_VERTEX"]
    VERTEX
    @[link_name="SDL_GPU_VERTEXINPUTRATE_INSTANCE"]
    INSTANCE
  }

  @[link_name="SDL_GPUVertexBufferDescription"]
  struct GPUVertexBufferDescription {
    slot u32
    pitch u32
    input_rate GPUVertexInputRate
    instance_step_rate u32
  }

  @[link_name="SDL_GPUVertexElementFormat"]
  enum GPUVertexElementFormat {
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INVALID"]
    INVALID
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INT"]
    INT
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INT2"]
    INT2
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INT3"]
    INT3
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_INT4"]
    INT4
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UINT"]
    UINT
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UINT2"]
    UINT2
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UINT3"]
    UINT3
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UINT4"]
    UINT4
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_FLOAT"]
    FLOAT
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_FLOAT2"]
    FLOAT2
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_FLOAT3"]
    FLOAT3
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_FLOAT4"]
    FLOAT4
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_BYTE2"]
    BYTE2
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_BYTE4"]
    BYTE4
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UBYTE2"]
    UBYTE2
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UBYTE4"]
    UBYTE4
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_BYTE2_NORM"]
    BYTE2_NORM
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_BYTE4_NORM"]
    BYTE4_NORM
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UBYTE2_NORM"]
    UBYTE2_NORM
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_UBYTE4_NORM"]
    UBYTE4_NORM
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_SHORT2"]
    SHORT2
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_SHORT4"]
    SHORT4
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_USHORT2"]
    USHORT2
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_USHORT4"]
    USHORT4
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_SHORT2_NORM"]
    SHORT2_NORM
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_SHORT4_NORM"]
    SHORT4_NORM
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_USHORT2_NORM"]
    USHORT2_NORM
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_USHORT4_NORM"]
    USHORT4_NORM
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_HALF2"]
    HALF2
    @[link_name="SDL_GPU_VERTEXELEMENTFORMAT_HALF4"]
    HALF4
  }

  @[link_name="SDL_GPUVertexAttribute"]
  struct GPUVertexAttribute {
    location u32
    buffer_slot u32
    format GPUVertexElementFormat
    offset u32
  }

  @[link_name="SDL_GPUVertexInputState"]
  struct GPUVertexInputState {
    vertex_buffer_descriptions const *GPUVertexBufferDescription
    num_vertex_buffers u32
    vertex_attributes const *GPUVertexAttribute
    num_vertex_attributes u32
  }

  @[link_name="SDL_GPUPrimitiveType"]
  enum GPUPrimitiveType {
    @[link_name="SDL_GPU_PRIMITIVETYPE_TRIANGLELIST"]
    TRIANGLELIST
    @[link_name="SDL_GPU_PRIMITIVETYPE_TRIANGLESTRIP"]
    TRIANGLESTRIP
    @[link_name="SDL_GPU_PRIMITIVETYPE_LINELIST"]
    LINELIST
    @[link_name="SDL_GPU_PRIMITIVETYPE_LINESTRIP"]
    LINESTRIP
    @[link_name="SDL_GPU_PRIMITIVETYPE_POINTLIST"]
    POINTLIST
  }

  @[link_name="SDL_GPUFillMode"]
  enum GPUFillMode {
    @[link_name="SDL_GPU_FILLMODE_FILL"]
    FILL
    @[link_name="SDL_GPU_FILLMODE_LINE"]
    LINE
  }

  @[link_name="SDL_GPUCullMode"]
  enum GPUCullMode {
    @[link_name="SDL_GPU_CULLMODE_NONE"]
    NONE
    @[link_name="SDL_GPU_CULLMODE_FRONT"]
    FRONT
    @[link_name="SDL_GPU_CULLMODE_BACK"]
    BACK
  }

  @[link_name="SDL_GPUFrontFace"]
  enum GPUFrontFace {
    @[link_name="SDL_GPU_FRONTFACE_COUNTER_CLOCKWISE"]
    COUNTER_CLOCKWISE
    @[link_name="SDL_GPU_FRONTFACE_CLOCKWISE"]
    CLOCKWISE
  }

  @[link_name="SDL_GPURasterizerState"]
  struct GPURasterizerState {
    fill_mode GPUFillMode
    cull_mode GPUCullMode
    front_face GPUFrontFace
    depth_bias_constant_factor libc::c_float
    depth_bias_clamp libc::c_float
    depth_bias_slope_factor libc::c_float
    enable_depth_bias bool
    enable_depth_clip bool
    padding1 u8
    padding2 u8
  }

  @[link_name="SDL_GPUSampleCount"]
  enum GPUSampleCount {
    @[link_name="SDL_GPU_SAMPLECOUNT_1"]
    GPU_SAMPLECOUNT_1
    @[link_name="SDL_GPU_SAMPLECOUNT_2"]
    GPU_SAMPLECOUNT_2
    @[link_name="SDL_GPU_SAMPLECOUNT_4"]
    GPU_SAMPLECOUNT_4
    @[link_name="SDL_GPU_SAMPLECOUNT_8"]
    GPU_SAMPLECOUNT_8
  }

  @[link_name="SDL_GPUMultisampleState"]
  struct GPUMultisampleState {
    sample_count GPUSampleCount
    sample_mask u32
    enable_mask bool
    padding1 u8
    padding2 u8
    padding3 u8
  }

  @[link_name="SDL_GPUCompareOp"]
  enum GPUCompareOp {
    @[link_name="SDL_GPU_COMPAREOP_INVALID"]
    INVALID
    @[link_name="SDL_GPU_COMPAREOP_NEVER"]
    NEVER
    @[link_name="SDL_GPU_COMPAREOP_LESS"]
    LESS
    @[link_name="SDL_GPU_COMPAREOP_EQUAL"]
    EQUAL
    @[link_name="SDL_GPU_COMPAREOP_LESS_OR_EQUAL"]
    LESS_OR_EQUAL
    @[link_name="SDL_GPU_COMPAREOP_GREATER"]
    GREATER
    @[link_name="SDL_GPU_COMPAREOP_NOT_EQUAL"]
    NOT_EQUAL
    @[link_name="SDL_GPU_COMPAREOP_GREATER_OR_EQUAL"]
    GREATER_OR_EQUAL
    @[link_name="SDL_GPU_COMPAREOP_ALWAYS"]
    ALWAYS
  }

  @[link_name="SDL_GPUStencilOp"]
  enum GPUStencilOp {
    @[link_name="SDL_GPU_STENCILOP_INVALID"]
    INVALID
    @[link_name="SDL_GPU_STENCILOP_KEEP"]
    KEEP
    @[link_name="SDL_GPU_STENCILOP_ZERO"]
    ZERO
    @[link_name="SDL_GPU_STENCILOP_REPLACE"]
    REPLACE
    @[link_name="SDL_GPU_STENCILOP_INCREMENT_AND_CLAMP"]
    INCREMENT_AND_CLAMP
    @[link_name="SDL_GPU_STENCILOP_DECREMENT_AND_CLAMP"]
    DECREMENT_AND_CLAMP
    @[link_name="SDL_GPU_STENCILOP_INVERT"]
    INVERT
    @[link_name="SDL_GPU_STENCILOP_INCREMENT_AND_WRAP"]
    INCREMENT_AND_WRAP
    @[link_name="SDL_GPU_STENCILOP_DECREMENT_AND_WRAP"]
    DECREMENT_AND_WRAP
  }

  @[link_name="SDL_GPUStencilOpState"]
  struct GPUStencilOpState {
    fail_op GPUStencilOp
    pass_op GPUStencilOp
    depth_fail_op GPUStencilOp
    compare_op GPUCompareOp
  }

  @[link_name="SDL_GPUDepthStencilState"]
  struct GPUDepthStencilState {
    compare_op GPUCompareOp
    back_stencil_state GPUStencilOpState
    front_stencil_state GPUStencilOpState
    compare_mask u8
    write_mask u8
    enable_depth_test bool
    enable_depth_write bool
    enable_stencil_test bool
    padding1 u8
    padding2 u8
    padding3 u8
  }

  @[link_name="SDL_GPUTextureFormat"]
  enum GPUTextureFormat {
    @[link_name="SDL_GPU_TEXTUREFORMAT_INVALID"]
    INVALID
    @[link_name="SDL_GPU_TEXTUREFORMAT_A8_UNORM"]
    A8_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8_UNORM"]
    R8_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8_UNORM"]
    R8G8_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_UNORM"]
    R8G8B8A8_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16_UNORM"]
    R16_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_UNORM"]
    R16G16_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_UNORM"]
    R16G16B16A16_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R10G10B10A2_UNORM"]
    R10G10B10A2_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_B5G6R5_UNORM"]
    B5G6R5_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_B5G5R5A1_UNORM"]
    B5G5R5A1_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_B4G4R4A4_UNORM"]
    B4G4R4A4_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_B8G8R8A8_UNORM"]
    B8G8R8A8_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC1_RGBA_UNORM"]
    BC1_RGBA_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC2_RGBA_UNORM"]
    BC2_RGBA_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC3_RGBA_UNORM"]
    BC3_RGBA_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC4_R_UNORM"]
    BC4_R_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC5_RG_UNORM"]
    BC5_RG_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC7_RGBA_UNORM"]
    BC7_RGBA_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC6H_RGB_FLOAT"]
    BC6H_RGB_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC6H_RGB_UFLOAT"]
    BC6H_RGB_UFLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8_SNORM"]
    R8_SNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8_SNORM"]
    R8G8_SNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_SNORM"]
    R8G8B8A8_SNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16_SNORM"]
    R16_SNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_SNORM"]
    R16G16_SNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_SNORM"]
    R16G16B16A16_SNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16_FLOAT"]
    R16_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_FLOAT"]
    R16G16_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_FLOAT"]
    R16G16B16A16_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R32_FLOAT"]
    R32_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32_FLOAT"]
    R32G32_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32B32A32_FLOAT"]
    R32G32B32A32_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R11G11B10_UFLOAT"]
    R11G11B10_UFLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8_UINT"]
    R8_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8_UINT"]
    R8G8_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_UINT"]
    R8G8B8A8_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16_UINT"]
    R16_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_UINT"]
    R16G16_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_UINT"]
    R16G16B16A16_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R32_UINT"]
    R32_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32_UINT"]
    R32G32_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32B32A32_UINT"]
    R32G32B32A32_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8_INT"]
    R8_INT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8_INT"]
    R8G8_INT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_INT"]
    R8G8B8A8_INT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16_INT"]
    R16_INT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16_INT"]
    R16G16_INT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R16G16B16A16_INT"]
    R16G16B16A16_INT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R32_INT"]
    R32_INT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32_INT"]
    R32G32_INT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R32G32B32A32_INT"]
    R32G32B32A32_INT
    @[link_name="SDL_GPU_TEXTUREFORMAT_R8G8B8A8_UNORM_SRGB"]
    R8G8B8A8_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_B8G8R8A8_UNORM_SRGB"]
    B8G8R8A8_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC1_RGBA_UNORM_SRGB"]
    BC1_RGBA_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC2_RGBA_UNORM_SRGB"]
    BC2_RGBA_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC3_RGBA_UNORM_SRGB"]
    BC3_RGBA_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_BC7_RGBA_UNORM_SRGB"]
    BC7_RGBA_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_D16_UNORM"]
    D16_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_D24_UNORM"]
    D24_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_D32_FLOAT"]
    D32_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_D24_UNORM_S8_UINT"]
    D24_UNORM_S8_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_D32_FLOAT_S8_UINT"]
    D32_FLOAT_S8_UINT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_4x4_UNORM"]
    ASTC_4x4_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x4_UNORM"]
    ASTC_5x4_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x5_UNORM"]
    ASTC_5x5_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x5_UNORM"]
    ASTC_6x5_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x6_UNORM"]
    ASTC_6x6_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x5_UNORM"]
    ASTC_8x5_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x6_UNORM"]
    ASTC_8x6_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x8_UNORM"]
    ASTC_8x8_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x5_UNORM"]
    ASTC_10x5_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x6_UNORM"]
    ASTC_10x6_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x8_UNORM"]
    ASTC_10x8_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x10_UNORM"]
    ASTC_10x10_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x10_UNORM"]
    ASTC_12x10_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x12_UNORM"]
    ASTC_12x12_UNORM
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_4x4_UNORM_SRGB"]
    ASTC_4x4_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x4_UNORM_SRGB"]
    ASTC_5x4_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x5_UNORM_SRGB"]
    ASTC_5x5_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x5_UNORM_SRGB"]
    ASTC_6x5_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x6_UNORM_SRGB"]
    ASTC_6x6_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x5_UNORM_SRGB"]
    ASTC_8x5_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x6_UNORM_SRGB"]
    ASTC_8x6_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x8_UNORM_SRGB"]
    ASTC_8x8_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x5_UNORM_SRGB"]
    ASTC_10x5_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x6_UNORM_SRGB"]
    ASTC_10x6_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x8_UNORM_SRGB"]
    ASTC_10x8_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x10_UNORM_SRGB"]
    ASTC_10x10_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x10_UNORM_SRGB"]
    ASTC_12x10_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x12_UNORM_SRGB"]
    ASTC_12x12_UNORM_SRGB
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_4x4_FLOAT"]
    ASTC_4x4_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x4_FLOAT"]
    ASTC_5x4_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_5x5_FLOAT"]
    ASTC_5x5_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x5_FLOAT"]
    ASTC_6x5_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_6x6_FLOAT"]
    ASTC_6x6_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x5_FLOAT"]
    ASTC_8x5_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x6_FLOAT"]
    ASTC_8x6_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_8x8_FLOAT"]
    ASTC_8x8_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x5_FLOAT"]
    ASTC_10x5_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x6_FLOAT"]
    ASTC_10x6_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x8_FLOAT"]
    ASTC_10x8_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_10x10_FLOAT"]
    ASTC_10x10_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x10_FLOAT"]
    ASTC_12x10_FLOAT
    @[link_name="SDL_GPU_TEXTUREFORMAT_ASTC_12x12_FLOAT"]
    ASTC_12x12_FLOAT
  }

  @[link_name="SDL_GPUBlendFactor"]
  enum GPUBlendFactor {
    @[link_name="SDL_GPU_BLENDFACTOR_INVALID"]
    INVALID
    @[link_name="SDL_GPU_BLENDFACTOR_ZERO"]
    ZERO
    @[link_name="SDL_GPU_BLENDFACTOR_ONE"]
    ONE
    @[link_name="SDL_GPU_BLENDFACTOR_SRC_COLOR"]
    SRC_COLOR
    @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_SRC_COLOR"]
    ONE_MINUS_SRC_COLOR
    @[link_name="SDL_GPU_BLENDFACTOR_DST_COLOR"]
    DST_COLOR
    @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_DST_COLOR"]
    ONE_MINUS_DST_COLOR
    @[link_name="SDL_GPU_BLENDFACTOR_SRC_ALPHA"]
    SRC_ALPHA
    @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_SRC_ALPHA"]
    ONE_MINUS_SRC_ALPHA
    @[link_name="SDL_GPU_BLENDFACTOR_DST_ALPHA"]
    DST_ALPHA
    @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_DST_ALPHA"]
    ONE_MINUS_DST_ALPHA
    @[link_name="SDL_GPU_BLENDFACTOR_CONSTANT_COLOR"]
    CONSTANT_COLOR
    @[link_name="SDL_GPU_BLENDFACTOR_ONE_MINUS_CONSTANT_COLOR"]
    ONE_MINUS_CONSTANT_COLOR
    @[link_name="SDL_GPU_BLENDFACTOR_SRC_ALPHA_SATURATE"]
    SRC_ALPHA_SATURATE
  }

  @[link_name="SDL_GPUBlendOp"]
  enum GPUBlendOp {
    @[link_name="SDL_GPU_BLENDOP_INVALID"]
    INVALID
    @[link_name="SDL_GPU_BLENDOP_ADD"]
    ADD
    @[link_name="SDL_GPU_BLENDOP_SUBTRACT"]
    SUBTRACT
    @[link_name="SDL_GPU_BLENDOP_REVERSE_SUBTRACT"]
    REVERSE_SUBTRACT
    @[link_name="SDL_GPU_BLENDOP_MIN"]
    MIN
    @[link_name="SDL_GPU_BLENDOP_MAX"]
    MAX
  }


  @[link_name="SDL_GPUColorTargetBlendState"]
  struct GPUColorTargetBlendState {
    src_color_blendfactor GPUBlendFactor
    dst_color_blendfactor GPUBlendFactor
    color_blend_op GPUBlendOp
    src_alpha_blendfactor GPUBlendFactor
    dst_alpha_blendfactor GPUBlendFactor
    alpha_blend_op GPUBlendOp
    color_write_mask GPUColorComponentFlags
    enable_blend bool
    enable_color_write_mask bool
    padding1 u8
    padding2 u8
  }

  @[link_name="SDL_GPUColorTargetDescription"]
  struct GPUColorTargetDescription {
    format GPUTextureFormat
    blend_state GPUColorTargetBlendState
  }

  @[link_name="SDL_GPUGraphicsPipelineTargetInfo"]
  struct GPUGraphicsPipelineTargetInfo {
    color_target_descriptions const *GPUColorTargetDescription
    num_color_targets u32
    depth_stencil_format GPUTextureFormat
    has_depth_stencil_target bool
    padding1 u8
    padding2 u8
    padding3 u8
  }

  @[link_name="SDL_GPUGraphicsPipelineCreateInfo"]
  struct GPUGraphicsPipelineCreateInfo {
    vertex_shader *GPUShader
    fragment_shader *GPUShader
    vertex_input_state GPUVertexInputState
    primitive_type GPUPrimitiveType
    rasterizer_state GPURasterizerState
    multisample_state GPUMultisampleState
    depth_stencil_state GPUDepthStencilState
    target_info GPUGraphicsPipelineTargetInfo
    props PropertiesID
  }

  @[link_name="SDL_GPUGraphicsPipeline"]
  opaque GPUGraphicsPipeline

  @[link_name="SDL_GPUFilter"]
  enum GPUFilter {
    @[link_name="SDL_GPU_FILTER_NEAREST"]
    NEAREST
    @[link_name="SDL_GPU_FILTER_LINEAR"]
    LINEAR
  }

  @[link_name="SDL_GPUSamplerMipmapMode"]
  enum GPUSamplerMipmapMode {
    @[link_name="SDL_GPU_SAMPLERMIPMAPMODE_NEAREST"]
    NEAREST
    @[link_name="SDL_GPU_SAMPLERMIPMAPMODE_LINEAR"]
    LINEAR
  }

  @[link_name="SDL_GPUSamplerAddressMode"]
  enum GPUSamplerAddressMode {
    @[link_name="SDL_GPU_SAMPLERADDRESSMODE_REPEAT"]
    REPEAT
    @[link_name="SDL_GPU_SAMPLERADDRESSMODE_MIRRORED_REPEAT"]
    MIRRORED_REPEAT
    @[link_name="SDL_GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE"]
    CLAMP_TO_EDGE
  }

  @[link_name="SDL_GPUSamplerCreateInfo"]
  struct GPUSamplerCreateInfo {
    min_filter GPUFilter
    mag_filter GPUFilter
    mipmap_mode GPUSamplerMipmapMode
    address_mode_u GPUSamplerAddressMode
    address_mode_v GPUSamplerAddressMode
    address_mode_w GPUSamplerAddressMode
    mip_lod_bias libc::c_float
    max_anisotropy libc::c_float
    compare_op GPUCompareOp
    min_lod libc::c_float
    max_lod libc::c_float
    enable_anisotropy bool
    enable_compare bool
    padding1 u8
    padding2 u8
    props PropertiesID
  }

  @[link_name="SDL_GPUSampler"]
  opaque GPUSampler

  @[link_name="SDL_GPUShaderStage"]
  enum GPUShaderStage {
    @[link_name="SDL_GPU_SHADERSTAGE_VERTEX"]
    VERTEX
    @[link_name="SDL_GPU_SHADERSTAGE_FRAGMENT"]
    FRAGMENT
  }

  @[link_name="SDL_GPUShaderCreateInfo"]
  struct GPUShaderCreateInfo {
    code_size libc::c_size_t
    code const *u8
    entrypoint const cstring
    format GPUShaderFormat
    stage GPUShaderStage
    num_samplers u32
    num_storage_textures u32
    num_storage_buffers u32
    num_uniform_buffers u32
    props PropertiesID
  }

  @[link_name="SDL_GPUTextureType"]
  enum GPUTextureType {
    @[link_name="SDL_GPU_TEXTURETYPE_2D"]
    GPU_TEXTURETYPE_2D
    @[link_name="SDL_GPU_TEXTURETYPE_2D_ARRAY"]
    GPU_TEXTURETYPE_2D_ARRAY
    @[link_name="SDL_GPU_TEXTURETYPE_3D"]
    GPU_TEXTURETYPE_3D
    @[link_name="SDL_GPU_TEXTURETYPE_CUBE"]
    CUBE
    @[link_name="SDL_GPU_TEXTURETYPE_CUBE_ARRAY"]
    CUBE_ARRAY
  }


  @[link_name="SDL_GPUTextureCreateInfo"]
  struct GPUTextureCreateInfo {
    type GPUTextureType
    format GPUTextureFormat
    usage GPUTextureUsageFlags
    width u32
    height u32
    layer_count_or_depth u32
    num_levels u32
    sample_count GPUSampleCount
    props PropertiesID
  }

  @[link_name="SDL_GPUTexture"]
  opaque GPUTexture


  @[link_name="SDL_GPUBufferCreateInfo"]
  struct GPUBufferCreateInfo {
    usage GPUBufferUsageFlags
    size u32
    props PropertiesID
  }

  @[link_name="SDL_GPUBuffer"]
  opaque GPUBuffer

  @[link_name="SDL_GPUTransferBufferUsage"]
  enum GPUTransferBufferUsage {
    @[link_name="SDL_GPU_TRANSFERBUFFERUSAGE_UPLOAD"]
    UPLOAD
    @[link_name="SDL_GPU_TRANSFERBUFFERUSAGE_DOWNLOAD"]
    DOWNLOAD
  }

  @[link_name="SDL_GPUTransferBufferCreateInfo"]
  struct GPUTransferBufferCreateInfo {
    usage GPUTransferBufferUsage
    size u32
    props PropertiesID
  }

  @[link_name="SDL_GPUTransferBuffer"]
  opaque GPUTransferBuffer

  @[link_name="SDL_GPUCommandBuffer"]
  opaque GPUCommandBuffer

  @[link_name="SDL_FColor"]
  struct FColor {
    r libc::c_float
    g libc::c_float
    b libc::c_float
    a libc::c_float
  }

  @[link_name="SDL_GPULoadOp"]
  enum GPULoadOp {
    @[link_name="SDL_GPU_LOADOP_LOAD"]
    LOAD
    @[link_name="SDL_GPU_LOADOP_CLEAR"]
    CLEAR
    @[link_name="SDL_GPU_LOADOP_DONT_CARE"]
    DONT_CARE
  }

  @[link_name="SDL_GPUStoreOp"]
  enum GPUStoreOp {
    @[link_name="SDL_GPU_STOREOP_STORE"]
    STORE
    @[link_name="SDL_GPU_STOREOP_DONT_CARE"]
    DONT_CARE
    @[link_name="SDL_GPU_STOREOP_RESOLVE"]
    RESOLVE
    @[link_name="SDL_GPU_STOREOP_RESOLVE_AND_STORE"]
    RESOLVE_AND_STORE
  }

  @[link_name="SDL_GPUColorTargetInfo"]
  struct GPUColorTargetInfo {
    texture *GPUTexture
    mip_level u32
    layer_or_depth_plane u32
    clear_color FColor
    load_op GPULoadOp
    store_op GPUStoreOp
    resolve_texture *GPUTexture
    resolve_mip_level u32
    resolve_layer u32
    cycle bool
    cycle_resolve_texture bool
    padding1 u8
    padding2 u8
  }

  @[link_name="SDL_GPUDepthStencilTargetInfo"]
  struct GPUDepthStencilTargetInfo {
    texture *GPUTexture
    clear_depth libc::c_float
    load_op GPULoadOp
    store_op GPUStoreOp
    stencil_load_op GPULoadOp
    stencil_store_op GPUStoreOp
    cycle bool
    clear_stencil u8
    padding1 u8
    padding2 u8
  }

  @[link_name="SDL_GPURenderPass"]
  opaque GPURenderPass

  @[link_name="SDL_GPUViewport"]
  struct GPUViewport {
    x libc::c_float
    y libc::c_float
    w libc::c_float
    h libc::c_float
    min_depth libc::c_float
    max_depth libc::c_float
  }

  @[link_name="SDL_GPUBufferBinding"]
  struct GPUBufferBinding {
    buffer *GPUBuffer
    offset u32
  }

  @[link_name="SDL_GPUIndexElementSize"]
  enum GPUIndexElementSize {
    @[link_name="SDL_GPU_INDEXELEMENTSIZE_16BIT"]
    GPU_INDEXELEMENTSIZE_16BIT
    @[link_name="SDL_GPU_INDEXELEMENTSIZE_32BIT"]
    GPU_INDEXELEMENTSIZE_32BIT
  }

  @[link_name="SDL_GPUTextureSamplerBinding"]
  struct GPUTextureSamplerBinding {
    texture *GPUTexture
    sampler *GPUSampler
  }

  @[link_name="SDL_GPUStorageTextureReadWriteBinding"]
  struct GPUStorageTextureReadWriteBinding {
    texture *GPUTexture
    mip_level u32
    layer u32
    cycle bool
    padding1 u8
    padding2 u8
    padding3 u8
  }

  @[link_name="SDL_GPUStorageBufferReadWriteBinding"]
  struct GPUStorageBufferReadWriteBinding {
    buffer *GPUBuffer
    cycle bool
    padding1 u8
    padding2 u8
    padding3 u8
  }

  @[link_name="SDL_GPUComputePass"]
  opaque GPUComputePass

  @[link_name="SDL_GPUCopyPass"]
  opaque GPUCopyPass

  @[link_name="SDL_GPUTextureTransferInfo"]
  struct GPUTextureTransferInfo {
    transfer_buffer *GPUTransferBuffer
    offset u32
    pixels_per_row u32
    rows_per_layer u32
  }

  @[link_name="SDL_GPUTextureRegion"]
  struct GPUTextureRegion {
    texture *GPUTexture
    mip_level u32
    layer u32
    x u32
    y u32
    z u32
    w u32
    h u32
    d u32
  }

  @[link_name="SDL_GPUTransferBufferLocation"]
  struct GPUTransferBufferLocation {
    transfer_buffer *GPUTransferBuffer
    offset u32
  }

  @[link_name="SDL_GPUBufferRegion"]
  struct GPUBufferRegion {
    buffer *GPUBuffer
    offset u32
    size u32
  }

  @[link_name="SDL_GPUTextureLocation"]
  struct GPUTextureLocation {
    texture *GPUTexture
    mip_level u32
    layer u32
    x u32
    y u32
    z u32
  }

  @[link_name="SDL_GPUBufferLocation"]
  struct GPUBufferLocation {
    buffer *GPUBuffer
    offset u32
  }

  @[link_name="SDL_GPUBlitRegion"]
  struct GPUBlitRegion {
    texture *GPUTexture
    mip_level u32
    layer_or_depth_plane u32
    x u32
    y u32
    w u32
    h u32
  }

  @[link_name="SDL_GPUBlitInfo"]
  struct GPUBlitInfo {
    source GPUBlitRegion
    destination GPUBlitRegion
    load_op GPULoadOp
    clear_color FColor
    flip_mode FlipMode
    filter GPUFilter
    cycle bool
    padding1 u8
    padding2 u8
    padding3 u8
  }

  @[link_name="SDL_GPUSwapchainComposition"]
  enum GPUSwapchainComposition {
    @[link_name="SDL_GPU_SWAPCHAINCOMPOSITION_SDR"]
    SDR
    @[link_name="SDL_GPU_SWAPCHAINCOMPOSITION_SDR_LINEAR"]
    SDR_LINEAR
    @[link_name="SDL_GPU_SWAPCHAINCOMPOSITION_HDR_EXTENDED_LINEAR"]
    HDR_EXTENDED_LINEAR
    @[link_name="SDL_GPU_SWAPCHAINCOMPOSITION_HDR10_ST2084"]
    HDR10_ST2084
  }

  @[link_name="SDL_GPUPresentMode"]
  enum GPUPresentMode {
    @[link_name="SDL_GPU_PRESENTMODE_VSYNC"]
    VSYNC
    @[link_name="SDL_GPU_PRESENTMODE_IMMEDIATE"]
    IMMEDIATE
    @[link_name="SDL_GPU_PRESENTMODE_MAILBOX"]
    MAILBOX
  }

  @[link_name="SDL_GPUFence"]
  opaque GPUFence


  @[link_name="SDL_Haptic"]
  opaque Haptic

  @[link_name="SDL_HapticDirection"]
  struct HapticDirection {
    type u8
    dir i32[3]
  }

  @[link_name="SDL_HapticConstant"]
  struct HapticConstant {
    type u16
    direction HapticDirection
    length u32
    delay u16
    button u16
    interval u16
    level i16
    attack_length u16
    attack_level u16
    fade_length u16
    fade_level u16
  }

  @[link_name="SDL_HapticPeriodic"]
  struct HapticPeriodic {
    type u16
    direction HapticDirection
    length u32
    delay u16
    button u16
    interval u16
    period u16
    magnitude i16
    offset i16
    phase u16
    attack_length u16
    attack_level u16
    fade_length u16
    fade_level u16
  }

  @[link_name="SDL_HapticCondition"]
  struct HapticCondition {
    type u16
    direction HapticDirection
    length u32
    delay u16
    button u16
    interval u16
    right_sat u16[3]
    left_sat u16[3]
    right_coeff i16[3]
    left_coeff i16[3]
    deadband u16[3]
    center i16[3]
  }

  @[link_name="SDL_HapticRamp"]
  struct HapticRamp {
    type u16
    direction HapticDirection
    length u32
    delay u16
    button u16
    interval u16
    start i16
    end i16
    attack_length u16
    attack_level u16
    fade_length u16
    fade_level u16
  }

  @[link_name="SDL_HapticLeftRight"]
  struct HapticLeftRight {
    type u16
    length u32
    large_magnitude u16
    small_magnitude u16
  }

  @[link_name="SDL_HapticCustom"]
  struct HapticCustom {
    type u16
    direction HapticDirection
    length u32
    delay u16
    button u16
    interval u16
    channels u8
    period u16
    samples u16
    data *u16
    attack_length u16
    attack_level u16
    fade_length u16
    fade_level u16
  }

  @[link_name="SDL_HapticEffect"]
  union HapticEffect {
    type u16
    constant HapticConstant
    periodic HapticPeriodic
    condition HapticCondition
    ramp HapticRamp
    leftright HapticLeftRight
    custom HapticCustom
  }

  @[link_name="SDL_hid_bus_type"]
  enum hid_bus_type {
    @[link_name="SDL_HID_API_BUS_UNKNOWN"]
    UNKNOWN
    @[link_name="SDL_HID_API_BUS_USB"]
    USB
    @[link_name="SDL_HID_API_BUS_BLUETOOTH"]
    BLUETOOTH
    @[link_name="SDL_HID_API_BUS_I2C"]
    I2C
    @[link_name="SDL_HID_API_BUS_SPI"]
    SPI
  }

  @[link_name="SDL_hid_device_info"]
  struct hid_device_info {
    path cstring
    vendor_id libc::c_ushort
    product_id libc::c_ushort
    serial_number *libc::c_wchar_t
    release_number libc::c_ushort
    manufacturer_string *libc::c_wchar_t
    product_string *libc::c_wchar_t
    usage_page libc::c_ushort
    usage libc::c_ushort
    interface_number libc::c_int
    interface_class libc::c_int
    interface_subclass libc::c_int
    interface_protocol libc::c_int
    bus_type hid_bus_type
    next *hid_device_info
  }

  @[link_name="SDL_hid_device"]
  opaque hid_device

  @[link_name="SDL_HintPriority"]
  enum HintPriority {
    @[link_name="SDL_HINT_DEFAULT"]
    DEFAULT
    @[link_name="SDL_HINT_NORMAL"]
    NORMAL
    @[link_name="SDL_HINT_OVERRIDE"]
    OVERRIDE
  }




  @[link_name="SDL_SharedObject"]
  opaque SharedObject

  @[link_name="SDL_Locale"]
  struct Locale {
    language const cstring
    country const cstring
  }

  @[link_name="SDL_LogPriority"]
  enum LogPriority {
    @[link_name="SDL_LOG_PRIORITY_INVALID"]
    INVALID
    @[link_name="SDL_LOG_PRIORITY_TRACE"]
    TRACE
    @[link_name="SDL_LOG_PRIORITY_VERBOSE"]
    VERBOSE
    @[link_name="SDL_LOG_PRIORITY_DEBUG"]
    DEBUG
    @[link_name="SDL_LOG_PRIORITY_INFO"]
    INFO
    @[link_name="SDL_LOG_PRIORITY_WARN"]
    WARN
    @[link_name="SDL_LOG_PRIORITY_ERROR"]
    ERROR
    @[link_name="SDL_LOG_PRIORITY_CRITICAL"]
    CRITICAL
    @[link_name="SDL_LOG_PRIORITY_COUNT"]
    COUNT
  }




  @[link_name="SDL_MessageBoxButtonData"]
  struct MessageBoxButtonData {
    flags MessageBoxButtonFlags
    buttonID libc::c_int
    text const cstring
  }

  @[link_name="SDL_MessageBoxColor"]
  struct MessageBoxColor {
    r u8
    g u8
    b u8
  }

  @[link_name="SDL_MessageBoxColorScheme"]
  struct MessageBoxColorScheme {
    colors MessageBoxColor[5]
  }

  @[link_name="SDL_MessageBoxData"]
  struct MessageBoxData {
    flags MessageBoxFlags
    window *Window
    title const cstring
    message const cstring
    numbuttons libc::c_int
    buttons const *MessageBoxButtonData
    colorScheme const *MessageBoxColorScheme
  }


  @[link_name="SDL_Process"]
  opaque Process

  @[link_name="SDL_Renderer"]
  opaque Renderer

  @[link_name="SDL_TextureAccess"]
  enum TextureAccess {
    @[link_name="SDL_TEXTUREACCESS_STATIC"]
    STATIC
    @[link_name="SDL_TEXTUREACCESS_STREAMING"]
    STREAMING
    @[link_name="SDL_TEXTUREACCESS_TARGET"]
    TARGET
  }

  @[link_name="SDL_Texture"]
  struct Texture {
    format PixelFormat
    w libc::c_int
    h libc::c_int
    refcount libc::c_int
  }

  @[link_name="SDL_RendererLogicalPresentation"]
  enum RendererLogicalPresentation {
    @[link_name="SDL_LOGICAL_PRESENTATION_DISABLED"]
    DISABLED
    @[link_name="SDL_LOGICAL_PRESENTATION_STRETCH"]
    STRETCH
    @[link_name="SDL_LOGICAL_PRESENTATION_LETTERBOX"]
    LETTERBOX
    @[link_name="SDL_LOGICAL_PRESENTATION_OVERSCAN"]
    OVERSCAN
    @[link_name="SDL_LOGICAL_PRESENTATION_INTEGER_SCALE"]
    INTEGER_SCALE
  }

  @[link_name="SDL_Vertex"]
  struct Vertex {
    position FPoint
    color FColor
    tex_coord FPoint
  }

  @[link_name="SDL_Storage"]
  opaque Storage

  @[link_name="SDL_StorageInterface"]
  struct StorageInterface {
    version u32
    close StorageInterface_closeCallback
    ready StorageInterface_readyCallback
    enumerate StorageInterface_enumerateCallback
    info StorageInterface_infoCallback
    read_file StorageInterface_read_fileCallback
    write_file StorageInterface_write_fileCallback
    mkdir StorageInterface_mkdirCallback
    remove StorageInterface_removeCallback
    rename StorageInterface_renameCallback
    copy StorageInterface_copyCallback
    space_remaining StorageInterface_space_remainingCallback
  }

  @[link_name="XEvent"]
  opaque XEvent


  @[link_name="SDL_Sandbox"]
  enum Sandbox {
    @[link_name="SDL_SANDBOX_NONE"]
    NONE
    @[link_name="SDL_SANDBOX_UNKNOWN_CONTAINER"]
    UNKNOWN_CONTAINER
    @[link_name="SDL_SANDBOX_FLATPAK"]
    FLATPAK
    @[link_name="SDL_SANDBOX_SNAP"]
    SNAP
    @[link_name="SDL_SANDBOX_MACOS"]
    MACOS
  }

  @[link_name="SDL_DateFormat"]
  enum DateFormat {
    @[link_name="SDL_DATE_FORMAT_YYYYMMDD"]
    YYYYMMDD
    @[link_name="SDL_DATE_FORMAT_DDMMYYYY"]
    DDMMYYYY
    @[link_name="SDL_DATE_FORMAT_MMDDYYYY"]
    MMDDYYYY
  }

  @[link_name="SDL_TimeFormat"]
  enum TimeFormat {
    @[link_name="SDL_TIME_FORMAT_24HR"]
    TIME_FORMAT_24HR
    @[link_name="SDL_TIME_FORMAT_12HR"]
    TIME_FORMAT_12HR
  }

  @[link_name="SDL_DateTime"]
  struct DateTime {
    year libc::c_int
    month libc::c_int
    day libc::c_int
    hour libc::c_int
    minute libc::c_int
    second libc::c_int
    nanosecond libc::c_int
    day_of_week libc::c_int
    utc_offset libc::c_int
  }




  @[link_name="SDL_Tray"]
  opaque Tray

  @[link_name="SDL_TrayMenu"]
  opaque TrayMenu

  @[link_name="SDL_TrayEntry"]
  opaque TrayEntry




  @[link_name="SDL_AppResult"]
  enum AppResult {
    @[link_name="SDL_APP_CONTINUE"]
    CONTINUE
    @[link_name="SDL_APP_SUCCESS"]
    SUCCESS
    @[link_name="SDL_APP_FAILURE"]
    FAILURE
  }





  @[link_name="struct VkInstance_T"]
  opaque VkInstance_T


  @[link_name="struct VkAllocationCallbacks"]
  opaque VkAllocationCallbacks

  @[link_name="struct VkSurfaceKHR_T"]
  opaque VkSurfaceKHR_T


  @[link_name="struct VkPhysicalDevice_T"]
  opaque VkPhysicalDevice_T


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

  @[link_name="SDL_HINT_ALLOW_ALT_TAB_WHILE_GRABBED"]
  const HINT_ALLOW_ALT_TAB_WHILE_GRABBED const cstring

  @[link_name="SDL_HINT_ANDROID_ALLOW_RECREATE_ACTIVITY"]
  const HINT_ANDROID_ALLOW_RECREATE_ACTIVITY const cstring

  @[link_name="SDL_HINT_ANDROID_BLOCK_ON_PAUSE"]
  const HINT_ANDROID_BLOCK_ON_PAUSE const cstring

  @[link_name="SDL_HINT_ANDROID_LOW_LATENCY_AUDIO"]
  const HINT_ANDROID_LOW_LATENCY_AUDIO const cstring

  @[link_name="SDL_HINT_ANDROID_TRAP_BACK_BUTTON"]
  const HINT_ANDROID_TRAP_BACK_BUTTON const cstring

  @[link_name="SDL_HINT_APPLE_TV_CONTROLLER_UI_EVENTS"]
  const HINT_APPLE_TV_CONTROLLER_UI_EVENTS const cstring

  @[link_name="SDL_HINT_APPLE_TV_REMOTE_ALLOW_ROTATION"]
  const HINT_APPLE_TV_REMOTE_ALLOW_ROTATION const cstring

  @[link_name="SDL_HINT_APP_ID"]
  const HINT_APP_ID const cstring

  @[link_name="SDL_HINT_APP_NAME"]
  const HINT_APP_NAME const cstring

  @[link_name="SDL_HINT_ASSERT"]
  const HINT_ASSERT const cstring

  @[link_name="SDL_HINT_AUDIO_ALSA_DEFAULT_DEVICE"]
  const HINT_AUDIO_ALSA_DEFAULT_DEVICE const cstring

  @[link_name="SDL_HINT_AUDIO_ALSA_DEFAULT_PLAYBACK_DEVICE"]
  const HINT_AUDIO_ALSA_DEFAULT_PLAYBACK_DEVICE const cstring

  @[link_name="SDL_HINT_AUDIO_ALSA_DEFAULT_RECORDING_DEVICE"]
  const HINT_AUDIO_ALSA_DEFAULT_RECORDING_DEVICE const cstring

  @[link_name="SDL_HINT_AUDIO_CATEGORY"]
  const HINT_AUDIO_CATEGORY const cstring

  @[link_name="SDL_HINT_AUDIO_CHANNELS"]
  const HINT_AUDIO_CHANNELS const cstring

  @[link_name="SDL_HINT_AUDIO_DEVICE_APP_ICON_NAME"]
  const HINT_AUDIO_DEVICE_APP_ICON_NAME const cstring

  @[link_name="SDL_HINT_AUDIO_DEVICE_SAMPLE_FRAMES"]
  const HINT_AUDIO_DEVICE_SAMPLE_FRAMES const cstring

  @[link_name="SDL_HINT_AUDIO_DEVICE_STREAM_NAME"]
  const HINT_AUDIO_DEVICE_STREAM_NAME const cstring

  @[link_name="SDL_HINT_AUDIO_DEVICE_STREAM_ROLE"]
  const HINT_AUDIO_DEVICE_STREAM_ROLE const cstring

  @[link_name="SDL_HINT_AUDIO_DISK_INPUT_FILE"]
  const HINT_AUDIO_DISK_INPUT_FILE const cstring

  @[link_name="SDL_HINT_AUDIO_DISK_OUTPUT_FILE"]
  const HINT_AUDIO_DISK_OUTPUT_FILE const cstring

  @[link_name="SDL_HINT_AUDIO_DISK_TIMESCALE"]
  const HINT_AUDIO_DISK_TIMESCALE const cstring

  @[link_name="SDL_HINT_AUDIO_DRIVER"]
  const HINT_AUDIO_DRIVER const cstring

  @[link_name="SDL_HINT_AUDIO_DUMMY_TIMESCALE"]
  const HINT_AUDIO_DUMMY_TIMESCALE const cstring

  @[link_name="SDL_HINT_AUDIO_FORMAT"]
  const HINT_AUDIO_FORMAT const cstring

  @[link_name="SDL_HINT_AUDIO_FREQUENCY"]
  const HINT_AUDIO_FREQUENCY const cstring

  @[link_name="SDL_HINT_AUDIO_INCLUDE_MONITORS"]
  const HINT_AUDIO_INCLUDE_MONITORS const cstring

  @[link_name="SDL_HINT_AUTO_UPDATE_JOYSTICKS"]
  const HINT_AUTO_UPDATE_JOYSTICKS const cstring

  @[link_name="SDL_HINT_AUTO_UPDATE_SENSORS"]
  const HINT_AUTO_UPDATE_SENSORS const cstring

  @[link_name="SDL_HINT_BMP_SAVE_LEGACY_FORMAT"]
  const HINT_BMP_SAVE_LEGACY_FORMAT const cstring

  @[link_name="SDL_HINT_CAMERA_DRIVER"]
  const HINT_CAMERA_DRIVER const cstring

  @[link_name="SDL_HINT_CPU_FEATURE_MASK"]
  const HINT_CPU_FEATURE_MASK const cstring

  @[link_name="SDL_HINT_DISPLAY_USABLE_BOUNDS"]
  const HINT_DISPLAY_USABLE_BOUNDS const cstring

  @[link_name="SDL_HINT_EGL_LIBRARY"]
  const HINT_EGL_LIBRARY const cstring

  @[link_name="SDL_HINT_EMSCRIPTEN_ASYNCIFY"]
  const HINT_EMSCRIPTEN_ASYNCIFY const cstring

  @[link_name="SDL_HINT_EMSCRIPTEN_CANVAS_SELECTOR"]
  const HINT_EMSCRIPTEN_CANVAS_SELECTOR const cstring

  @[link_name="SDL_HINT_EMSCRIPTEN_KEYBOARD_ELEMENT"]
  const HINT_EMSCRIPTEN_KEYBOARD_ELEMENT const cstring

  @[link_name="SDL_HINT_ENABLE_SCREEN_KEYBOARD"]
  const HINT_ENABLE_SCREEN_KEYBOARD const cstring

  @[link_name="SDL_HINT_EVDEV_DEVICES"]
  const HINT_EVDEV_DEVICES const cstring

  @[link_name="SDL_HINT_EVENT_LOGGING"]
  const HINT_EVENT_LOGGING const cstring

  @[link_name="SDL_HINT_FILE_DIALOG_DRIVER"]
  const HINT_FILE_DIALOG_DRIVER const cstring

  @[link_name="SDL_HINT_FORCE_RAISEWINDOW"]
  const HINT_FORCE_RAISEWINDOW const cstring

  @[link_name="SDL_HINT_FRAMEBUFFER_ACCELERATION"]
  const HINT_FRAMEBUFFER_ACCELERATION const cstring

  @[link_name="SDL_HINT_GAMECONTROLLERCONFIG"]
  const HINT_GAMECONTROLLERCONFIG const cstring

  @[link_name="SDL_HINT_GAMECONTROLLERCONFIG_FILE"]
  const HINT_GAMECONTROLLERCONFIG_FILE const cstring

  @[link_name="SDL_HINT_GAMECONTROLLERTYPE"]
  const HINT_GAMECONTROLLERTYPE const cstring

  @[link_name="SDL_HINT_GAMECONTROLLER_IGNORE_DEVICES"]
  const HINT_GAMECONTROLLER_IGNORE_DEVICES const cstring

  @[link_name="SDL_HINT_GAMECONTROLLER_IGNORE_DEVICES_EXCEPT"]
  const HINT_GAMECONTROLLER_IGNORE_DEVICES_EXCEPT const cstring

  @[link_name="SDL_HINT_GAMECONTROLLER_SENSOR_FUSION"]
  const HINT_GAMECONTROLLER_SENSOR_FUSION const cstring

  @[link_name="SDL_HINT_GDK_TEXTINPUT_DEFAULT_TEXT"]
  const HINT_GDK_TEXTINPUT_DEFAULT_TEXT const cstring

  @[link_name="SDL_HINT_GDK_TEXTINPUT_DESCRIPTION"]
  const HINT_GDK_TEXTINPUT_DESCRIPTION const cstring

  @[link_name="SDL_HINT_GDK_TEXTINPUT_MAX_LENGTH"]
  const HINT_GDK_TEXTINPUT_MAX_LENGTH const cstring

  @[link_name="SDL_HINT_GDK_TEXTINPUT_SCOPE"]
  const HINT_GDK_TEXTINPUT_SCOPE const cstring

  @[link_name="SDL_HINT_GDK_TEXTINPUT_TITLE"]
  const HINT_GDK_TEXTINPUT_TITLE const cstring

  @[link_name="SDL_HINT_GPU_DRIVER"]
  const HINT_GPU_DRIVER const cstring

  @[link_name="SDL_HINT_HIDAPI_ENUMERATE_ONLY_CONTROLLERS"]
  const HINT_HIDAPI_ENUMERATE_ONLY_CONTROLLERS const cstring

  @[link_name="SDL_HINT_HIDAPI_IGNORE_DEVICES"]
  const HINT_HIDAPI_IGNORE_DEVICES const cstring

  @[link_name="SDL_HINT_HIDAPI_LIBUSB"]
  const HINT_HIDAPI_LIBUSB const cstring

  @[link_name="SDL_HINT_HIDAPI_LIBUSB_WHITELIST"]
  const HINT_HIDAPI_LIBUSB_WHITELIST const cstring

  @[link_name="SDL_HINT_HIDAPI_UDEV"]
  const HINT_HIDAPI_UDEV const cstring

  @[link_name="SDL_HINT_IME_IMPLEMENTED_UI"]
  const HINT_IME_IMPLEMENTED_UI const cstring

  @[link_name="SDL_HINT_IOS_HIDE_HOME_INDICATOR"]
  const HINT_IOS_HIDE_HOME_INDICATOR const cstring

  @[link_name="SDL_HINT_JOYSTICK_ALLOW_BACKGROUND_EVENTS"]
  const HINT_JOYSTICK_ALLOW_BACKGROUND_EVENTS const cstring

  @[link_name="SDL_HINT_JOYSTICK_ARCADESTICK_DEVICES"]
  const HINT_JOYSTICK_ARCADESTICK_DEVICES const cstring

  @[link_name="SDL_HINT_JOYSTICK_ARCADESTICK_DEVICES_EXCLUDED"]
  const HINT_JOYSTICK_ARCADESTICK_DEVICES_EXCLUDED const cstring

  @[link_name="SDL_HINT_JOYSTICK_BLACKLIST_DEVICES"]
  const HINT_JOYSTICK_BLACKLIST_DEVICES const cstring

  @[link_name="SDL_HINT_JOYSTICK_BLACKLIST_DEVICES_EXCLUDED"]
  const HINT_JOYSTICK_BLACKLIST_DEVICES_EXCLUDED const cstring

  @[link_name="SDL_HINT_JOYSTICK_DEVICE"]
  const HINT_JOYSTICK_DEVICE const cstring

  @[link_name="SDL_HINT_JOYSTICK_DIRECTINPUT"]
  const HINT_JOYSTICK_DIRECTINPUT const cstring

  @[link_name="SDL_HINT_JOYSTICK_ENHANCED_REPORTS"]
  const HINT_JOYSTICK_ENHANCED_REPORTS const cstring

  @[link_name="SDL_HINT_JOYSTICK_FLIGHTSTICK_DEVICES"]
  const HINT_JOYSTICK_FLIGHTSTICK_DEVICES const cstring

  @[link_name="SDL_HINT_JOYSTICK_FLIGHTSTICK_DEVICES_EXCLUDED"]
  const HINT_JOYSTICK_FLIGHTSTICK_DEVICES_EXCLUDED const cstring

  @[link_name="SDL_HINT_JOYSTICK_GAMECUBE_DEVICES"]
  const HINT_JOYSTICK_GAMECUBE_DEVICES const cstring

  @[link_name="SDL_HINT_JOYSTICK_GAMECUBE_DEVICES_EXCLUDED"]
  const HINT_JOYSTICK_GAMECUBE_DEVICES_EXCLUDED const cstring

  @[link_name="SDL_HINT_JOYSTICK_GAMEINPUT"]
  const HINT_JOYSTICK_GAMEINPUT const cstring

  @[link_name="SDL_HINT_JOYSTICK_HAPTIC_AXES"]
  const HINT_JOYSTICK_HAPTIC_AXES const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI"]
  const HINT_JOYSTICK_HIDAPI const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_COMBINE_JOY_CONS"]
  const HINT_JOYSTICK_HIDAPI_COMBINE_JOY_CONS const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_GAMECUBE"]
  const HINT_JOYSTICK_HIDAPI_GAMECUBE const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_GAMECUBE_RUMBLE_BRAKE"]
  const HINT_JOYSTICK_HIDAPI_GAMECUBE_RUMBLE_BRAKE const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_JOYCON_HOME_LED"]
  const HINT_JOYSTICK_HIDAPI_JOYCON_HOME_LED const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_JOY_CONS"]
  const HINT_JOYSTICK_HIDAPI_JOY_CONS const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_LUNA"]
  const HINT_JOYSTICK_HIDAPI_LUNA const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_NINTENDO_CLASSIC"]
  const HINT_JOYSTICK_HIDAPI_NINTENDO_CLASSIC const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_PS3"]
  const HINT_JOYSTICK_HIDAPI_PS3 const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_PS3_SIXAXIS_DRIVER"]
  const HINT_JOYSTICK_HIDAPI_PS3_SIXAXIS_DRIVER const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_PS4"]
  const HINT_JOYSTICK_HIDAPI_PS4 const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_PS4_REPORT_INTERVAL"]
  const HINT_JOYSTICK_HIDAPI_PS4_REPORT_INTERVAL const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_PS5"]
  const HINT_JOYSTICK_HIDAPI_PS5 const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_PS5_PLAYER_LED"]
  const HINT_JOYSTICK_HIDAPI_PS5_PLAYER_LED const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_SHIELD"]
  const HINT_JOYSTICK_HIDAPI_SHIELD const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_STADIA"]
  const HINT_JOYSTICK_HIDAPI_STADIA const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_STEAM"]
  const HINT_JOYSTICK_HIDAPI_STEAM const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_STEAMDECK"]
  const HINT_JOYSTICK_HIDAPI_STEAMDECK const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_STEAM_HOME_LED"]
  const HINT_JOYSTICK_HIDAPI_STEAM_HOME_LED const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_STEAM_HORI"]
  const HINT_JOYSTICK_HIDAPI_STEAM_HORI const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_SWITCH"]
  const HINT_JOYSTICK_HIDAPI_SWITCH const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_SWITCH_HOME_LED"]
  const HINT_JOYSTICK_HIDAPI_SWITCH_HOME_LED const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_SWITCH_PLAYER_LED"]
  const HINT_JOYSTICK_HIDAPI_SWITCH_PLAYER_LED const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_VERTICAL_JOY_CONS"]
  const HINT_JOYSTICK_HIDAPI_VERTICAL_JOY_CONS const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_WII"]
  const HINT_JOYSTICK_HIDAPI_WII const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_WII_PLAYER_LED"]
  const HINT_JOYSTICK_HIDAPI_WII_PLAYER_LED const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_XBOX"]
  const HINT_JOYSTICK_HIDAPI_XBOX const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_XBOX_360"]
  const HINT_JOYSTICK_HIDAPI_XBOX_360 const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_XBOX_360_PLAYER_LED"]
  const HINT_JOYSTICK_HIDAPI_XBOX_360_PLAYER_LED const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_XBOX_360_WIRELESS"]
  const HINT_JOYSTICK_HIDAPI_XBOX_360_WIRELESS const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_XBOX_ONE"]
  const HINT_JOYSTICK_HIDAPI_XBOX_ONE const cstring

  @[link_name="SDL_HINT_JOYSTICK_HIDAPI_XBOX_ONE_HOME_LED"]
  const HINT_JOYSTICK_HIDAPI_XBOX_ONE_HOME_LED const cstring

  @[link_name="SDL_HINT_JOYSTICK_IOKIT"]
  const HINT_JOYSTICK_IOKIT const cstring

  @[link_name="SDL_HINT_JOYSTICK_LINUX_CLASSIC"]
  const HINT_JOYSTICK_LINUX_CLASSIC const cstring

  @[link_name="SDL_HINT_JOYSTICK_LINUX_DEADZONES"]
  const HINT_JOYSTICK_LINUX_DEADZONES const cstring

  @[link_name="SDL_HINT_JOYSTICK_LINUX_DIGITAL_HATS"]
  const HINT_JOYSTICK_LINUX_DIGITAL_HATS const cstring

  @[link_name="SDL_HINT_JOYSTICK_LINUX_HAT_DEADZONES"]
  const HINT_JOYSTICK_LINUX_HAT_DEADZONES const cstring

  @[link_name="SDL_HINT_JOYSTICK_MFI"]
  const HINT_JOYSTICK_MFI const cstring

  @[link_name="SDL_HINT_JOYSTICK_RAWINPUT"]
  const HINT_JOYSTICK_RAWINPUT const cstring

  @[link_name="SDL_HINT_JOYSTICK_RAWINPUT_CORRELATE_XINPUT"]
  const HINT_JOYSTICK_RAWINPUT_CORRELATE_XINPUT const cstring

  @[link_name="SDL_HINT_JOYSTICK_ROG_CHAKRAM"]
  const HINT_JOYSTICK_ROG_CHAKRAM const cstring

  @[link_name="SDL_HINT_JOYSTICK_THREAD"]
  const HINT_JOYSTICK_THREAD const cstring

  @[link_name="SDL_HINT_JOYSTICK_THROTTLE_DEVICES"]
  const HINT_JOYSTICK_THROTTLE_DEVICES const cstring

  @[link_name="SDL_HINT_JOYSTICK_THROTTLE_DEVICES_EXCLUDED"]
  const HINT_JOYSTICK_THROTTLE_DEVICES_EXCLUDED const cstring

  @[link_name="SDL_HINT_JOYSTICK_WGI"]
  const HINT_JOYSTICK_WGI const cstring

  @[link_name="SDL_HINT_JOYSTICK_WHEEL_DEVICES"]
  const HINT_JOYSTICK_WHEEL_DEVICES const cstring

  @[link_name="SDL_HINT_JOYSTICK_WHEEL_DEVICES_EXCLUDED"]
  const HINT_JOYSTICK_WHEEL_DEVICES_EXCLUDED const cstring

  @[link_name="SDL_HINT_JOYSTICK_ZERO_CENTERED_DEVICES"]
  const HINT_JOYSTICK_ZERO_CENTERED_DEVICES const cstring

  @[link_name="SDL_HINT_KEYCODE_OPTIONS"]
  const HINT_KEYCODE_OPTIONS const cstring

  @[link_name="SDL_HINT_KMSDRM_DEVICE_INDEX"]
  const HINT_KMSDRM_DEVICE_INDEX const cstring

  @[link_name="SDL_HINT_KMSDRM_REQUIRE_DRM_MASTER"]
  const HINT_KMSDRM_REQUIRE_DRM_MASTER const cstring

  @[link_name="SDL_HINT_LOGGING"]
  const HINT_LOGGING const cstring

  @[link_name="SDL_HINT_MAC_BACKGROUND_APP"]
  const HINT_MAC_BACKGROUND_APP const cstring

  @[link_name="SDL_HINT_MAC_CTRL_CLICK_EMULATE_RIGHT_CLICK"]
  const HINT_MAC_CTRL_CLICK_EMULATE_RIGHT_CLICK const cstring

  @[link_name="SDL_HINT_MAC_OPENGL_ASYNC_DISPATCH"]
  const HINT_MAC_OPENGL_ASYNC_DISPATCH const cstring

  @[link_name="SDL_HINT_MAC_OPTION_AS_ALT"]
  const HINT_MAC_OPTION_AS_ALT const cstring

  @[link_name="SDL_HINT_MAC_SCROLL_MOMENTUM"]
  const HINT_MAC_SCROLL_MOMENTUM const cstring

  @[link_name="SDL_HINT_MAIN_CALLBACK_RATE"]
  const HINT_MAIN_CALLBACK_RATE const cstring

  @[link_name="SDL_HINT_MOUSE_AUTO_CAPTURE"]
  const HINT_MOUSE_AUTO_CAPTURE const cstring

  @[link_name="SDL_HINT_MOUSE_DEFAULT_SYSTEM_CURSOR"]
  const HINT_MOUSE_DEFAULT_SYSTEM_CURSOR const cstring

  @[link_name="SDL_HINT_MOUSE_DOUBLE_CLICK_RADIUS"]
  const HINT_MOUSE_DOUBLE_CLICK_RADIUS const cstring

  @[link_name="SDL_HINT_MOUSE_DOUBLE_CLICK_TIME"]
  const HINT_MOUSE_DOUBLE_CLICK_TIME const cstring

  @[link_name="SDL_HINT_MOUSE_EMULATE_WARP_WITH_RELATIVE"]
  const HINT_MOUSE_EMULATE_WARP_WITH_RELATIVE const cstring

  @[link_name="SDL_HINT_MOUSE_FOCUS_CLICKTHROUGH"]
  const HINT_MOUSE_FOCUS_CLICKTHROUGH const cstring

  @[link_name="SDL_HINT_MOUSE_NORMAL_SPEED_SCALE"]
  const HINT_MOUSE_NORMAL_SPEED_SCALE const cstring

  @[link_name="SDL_HINT_MOUSE_RELATIVE_CURSOR_VISIBLE"]
  const HINT_MOUSE_RELATIVE_CURSOR_VISIBLE const cstring

  @[link_name="SDL_HINT_MOUSE_RELATIVE_MODE_CENTER"]
  const HINT_MOUSE_RELATIVE_MODE_CENTER const cstring

  @[link_name="SDL_HINT_MOUSE_RELATIVE_SPEED_SCALE"]
  const HINT_MOUSE_RELATIVE_SPEED_SCALE const cstring

  @[link_name="SDL_HINT_MOUSE_RELATIVE_SYSTEM_SCALE"]
  const HINT_MOUSE_RELATIVE_SYSTEM_SCALE const cstring

  @[link_name="SDL_HINT_MOUSE_RELATIVE_WARP_MOTION"]
  const HINT_MOUSE_RELATIVE_WARP_MOTION const cstring

  @[link_name="SDL_HINT_MOUSE_TOUCH_EVENTS"]
  const HINT_MOUSE_TOUCH_EVENTS const cstring

  @[link_name="SDL_HINT_MUTE_CONSOLE_KEYBOARD"]
  const HINT_MUTE_CONSOLE_KEYBOARD const cstring

  @[link_name="SDL_HINT_NO_SIGNAL_HANDLERS"]
  const HINT_NO_SIGNAL_HANDLERS const cstring

  @[link_name="SDL_HINT_OPENGL_ES_DRIVER"]
  const HINT_OPENGL_ES_DRIVER const cstring

  @[link_name="SDL_HINT_OPENGL_LIBRARY"]
  const HINT_OPENGL_LIBRARY const cstring

  @[link_name="SDL_HINT_OPENVR_LIBRARY"]
  const HINT_OPENVR_LIBRARY const cstring

  @[link_name="SDL_HINT_ORIENTATIONS"]
  const HINT_ORIENTATIONS const cstring

  @[link_name="SDL_HINT_PEN_MOUSE_EVENTS"]
  const HINT_PEN_MOUSE_EVENTS const cstring

  @[link_name="SDL_HINT_PEN_TOUCH_EVENTS"]
  const HINT_PEN_TOUCH_EVENTS const cstring

  @[link_name="SDL_HINT_POLL_SENTINEL"]
  const HINT_POLL_SENTINEL const cstring

  @[link_name="SDL_HINT_PREFERRED_LOCALES"]
  const HINT_PREFERRED_LOCALES const cstring

  @[link_name="SDL_HINT_QUIT_ON_LAST_WINDOW_CLOSE"]
  const HINT_QUIT_ON_LAST_WINDOW_CLOSE const cstring

  @[link_name="SDL_HINT_RENDER_DIRECT3D11_DEBUG"]
  const HINT_RENDER_DIRECT3D11_DEBUG const cstring

  @[link_name="SDL_HINT_RENDER_DIRECT3D_THREADSAFE"]
  const HINT_RENDER_DIRECT3D_THREADSAFE const cstring

  @[link_name="SDL_HINT_RENDER_DRIVER"]
  const HINT_RENDER_DRIVER const cstring

  @[link_name="SDL_HINT_RENDER_GPU_DEBUG"]
  const HINT_RENDER_GPU_DEBUG const cstring

  @[link_name="SDL_HINT_RENDER_GPU_LOW_POWER"]
  const HINT_RENDER_GPU_LOW_POWER const cstring

  @[link_name="SDL_HINT_RENDER_LINE_METHOD"]
  const HINT_RENDER_LINE_METHOD const cstring

  @[link_name="SDL_HINT_RENDER_METAL_PREFER_LOW_POWER_DEVICE"]
  const HINT_RENDER_METAL_PREFER_LOW_POWER_DEVICE const cstring

  @[link_name="SDL_HINT_RENDER_VSYNC"]
  const HINT_RENDER_VSYNC const cstring

  @[link_name="SDL_HINT_RENDER_VULKAN_DEBUG"]
  const HINT_RENDER_VULKAN_DEBUG const cstring

  @[link_name="SDL_HINT_RETURN_KEY_HIDES_IME"]
  const HINT_RETURN_KEY_HIDES_IME const cstring

  @[link_name="SDL_HINT_ROG_GAMEPAD_MICE"]
  const HINT_ROG_GAMEPAD_MICE const cstring

  @[link_name="SDL_HINT_ROG_GAMEPAD_MICE_EXCLUDED"]
  const HINT_ROG_GAMEPAD_MICE_EXCLUDED const cstring

  @[link_name="SDL_HINT_RPI_VIDEO_LAYER"]
  const HINT_RPI_VIDEO_LAYER const cstring

  @[link_name="SDL_HINT_SCREENSAVER_INHIBIT_ACTIVITY_NAME"]
  const HINT_SCREENSAVER_INHIBIT_ACTIVITY_NAME const cstring

  @[link_name="SDL_HINT_SHUTDOWN_DBUS_ON_QUIT"]
  const HINT_SHUTDOWN_DBUS_ON_QUIT const cstring

  @[link_name="SDL_HINT_STORAGE_TITLE_DRIVER"]
  const HINT_STORAGE_TITLE_DRIVER const cstring

  @[link_name="SDL_HINT_STORAGE_USER_DRIVER"]
  const HINT_STORAGE_USER_DRIVER const cstring

  @[link_name="SDL_HINT_THREAD_FORCE_REALTIME_TIME_CRITICAL"]
  const HINT_THREAD_FORCE_REALTIME_TIME_CRITICAL const cstring

  @[link_name="SDL_HINT_THREAD_PRIORITY_POLICY"]
  const HINT_THREAD_PRIORITY_POLICY const cstring

  @[link_name="SDL_HINT_TIMER_RESOLUTION"]
  const HINT_TIMER_RESOLUTION const cstring

  @[link_name="SDL_HINT_TOUCH_MOUSE_EVENTS"]
  const HINT_TOUCH_MOUSE_EVENTS const cstring

  @[link_name="SDL_HINT_TRACKPAD_IS_TOUCH_ONLY"]
  const HINT_TRACKPAD_IS_TOUCH_ONLY const cstring

  @[link_name="SDL_HINT_TV_REMOTE_AS_JOYSTICK"]
  const HINT_TV_REMOTE_AS_JOYSTICK const cstring

  @[link_name="SDL_HINT_VIDEO_ALLOW_SCREENSAVER"]
  const HINT_VIDEO_ALLOW_SCREENSAVER const cstring

  @[link_name="SDL_HINT_VIDEO_DISPLAY_PRIORITY"]
  const HINT_VIDEO_DISPLAY_PRIORITY const cstring

  @[link_name="SDL_HINT_VIDEO_DOUBLE_BUFFER"]
  const HINT_VIDEO_DOUBLE_BUFFER const cstring

  @[link_name="SDL_HINT_VIDEO_DRIVER"]
  const HINT_VIDEO_DRIVER const cstring

  @[link_name="SDL_HINT_VIDEO_DUMMY_SAVE_FRAMES"]
  const HINT_VIDEO_DUMMY_SAVE_FRAMES const cstring

  @[link_name="SDL_HINT_VIDEO_EGL_ALLOW_GETDISPLAY_FALLBACK"]
  const HINT_VIDEO_EGL_ALLOW_GETDISPLAY_FALLBACK const cstring

  @[link_name="SDL_HINT_VIDEO_FORCE_EGL"]
  const HINT_VIDEO_FORCE_EGL const cstring

  @[link_name="SDL_HINT_VIDEO_MAC_FULLSCREEN_MENU_VISIBILITY"]
  const HINT_VIDEO_MAC_FULLSCREEN_MENU_VISIBILITY const cstring

  @[link_name="SDL_HINT_VIDEO_MAC_FULLSCREEN_SPACES"]
  const HINT_VIDEO_MAC_FULLSCREEN_SPACES const cstring

  @[link_name="SDL_HINT_VIDEO_MINIMIZE_ON_FOCUS_LOSS"]
  const HINT_VIDEO_MINIMIZE_ON_FOCUS_LOSS const cstring

  @[link_name="SDL_HINT_VIDEO_OFFSCREEN_SAVE_FRAMES"]
  const HINT_VIDEO_OFFSCREEN_SAVE_FRAMES const cstring

  @[link_name="SDL_HINT_VIDEO_SYNC_WINDOW_OPERATIONS"]
  const HINT_VIDEO_SYNC_WINDOW_OPERATIONS const cstring

  @[link_name="SDL_HINT_VIDEO_WAYLAND_ALLOW_LIBDECOR"]
  const HINT_VIDEO_WAYLAND_ALLOW_LIBDECOR const cstring

  @[link_name="SDL_HINT_VIDEO_WAYLAND_MODE_EMULATION"]
  const HINT_VIDEO_WAYLAND_MODE_EMULATION const cstring

  @[link_name="SDL_HINT_VIDEO_WAYLAND_MODE_SCALING"]
  const HINT_VIDEO_WAYLAND_MODE_SCALING const cstring

  @[link_name="SDL_HINT_VIDEO_WAYLAND_PREFER_LIBDECOR"]
  const HINT_VIDEO_WAYLAND_PREFER_LIBDECOR const cstring

  @[link_name="SDL_HINT_VIDEO_WAYLAND_SCALE_TO_DISPLAY"]
  const HINT_VIDEO_WAYLAND_SCALE_TO_DISPLAY const cstring

  @[link_name="SDL_HINT_VIDEO_WIN_D3DCOMPILER"]
  const HINT_VIDEO_WIN_D3DCOMPILER const cstring

  @[link_name="SDL_HINT_VIDEO_X11_EXTERNAL_WINDOW_INPUT"]
  const HINT_VIDEO_X11_EXTERNAL_WINDOW_INPUT const cstring

  @[link_name="SDL_HINT_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR"]
  const HINT_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR const cstring

  @[link_name="SDL_HINT_VIDEO_X11_NET_WM_PING"]
  const HINT_VIDEO_X11_NET_WM_PING const cstring

  @[link_name="SDL_HINT_VIDEO_X11_NODIRECTCOLOR"]
  const HINT_VIDEO_X11_NODIRECTCOLOR const cstring

  @[link_name="SDL_HINT_VIDEO_X11_SCALING_FACTOR"]
  const HINT_VIDEO_X11_SCALING_FACTOR const cstring

  @[link_name="SDL_HINT_VIDEO_X11_VISUALID"]
  const HINT_VIDEO_X11_VISUALID const cstring

  @[link_name="SDL_HINT_VIDEO_X11_WINDOW_VISUALID"]
  const HINT_VIDEO_X11_WINDOW_VISUALID const cstring

  @[link_name="SDL_HINT_VIDEO_X11_XRANDR"]
  const HINT_VIDEO_X11_XRANDR const cstring

  @[link_name="SDL_HINT_VITA_ENABLE_BACK_TOUCH"]
  const HINT_VITA_ENABLE_BACK_TOUCH const cstring

  @[link_name="SDL_HINT_VITA_ENABLE_FRONT_TOUCH"]
  const HINT_VITA_ENABLE_FRONT_TOUCH const cstring

  @[link_name="SDL_HINT_VITA_MODULE_PATH"]
  const HINT_VITA_MODULE_PATH const cstring

  @[link_name="SDL_HINT_VITA_PVR_INIT"]
  const HINT_VITA_PVR_INIT const cstring

  @[link_name="SDL_HINT_VITA_PVR_OPENGL"]
  const HINT_VITA_PVR_OPENGL const cstring

  @[link_name="SDL_HINT_VITA_RESOLUTION"]
  const HINT_VITA_RESOLUTION const cstring

  @[link_name="SDL_HINT_VITA_TOUCH_MOUSE_DEVICE"]
  const HINT_VITA_TOUCH_MOUSE_DEVICE const cstring

  @[link_name="SDL_HINT_VULKAN_DISPLAY"]
  const HINT_VULKAN_DISPLAY const cstring

  @[link_name="SDL_HINT_VULKAN_LIBRARY"]
  const HINT_VULKAN_LIBRARY const cstring

  @[link_name="SDL_HINT_WAVE_CHUNK_LIMIT"]
  const HINT_WAVE_CHUNK_LIMIT const cstring

  @[link_name="SDL_HINT_WAVE_FACT_CHUNK"]
  const HINT_WAVE_FACT_CHUNK const cstring

  @[link_name="SDL_HINT_WAVE_RIFF_CHUNK_SIZE"]
  const HINT_WAVE_RIFF_CHUNK_SIZE const cstring

  @[link_name="SDL_HINT_WAVE_TRUNCATION"]
  const HINT_WAVE_TRUNCATION const cstring

  @[link_name="SDL_HINT_WINDOWS_CLOSE_ON_ALT_F4"]
  const HINT_WINDOWS_CLOSE_ON_ALT_F4 const cstring

  @[link_name="SDL_HINT_WINDOWS_ENABLE_MENU_MNEMONICS"]
  const HINT_WINDOWS_ENABLE_MENU_MNEMONICS const cstring

  @[link_name="SDL_HINT_WINDOWS_ENABLE_MESSAGELOOP"]
  const HINT_WINDOWS_ENABLE_MESSAGELOOP const cstring

  @[link_name="SDL_HINT_WINDOWS_ERASE_BACKGROUND_MODE"]
  const HINT_WINDOWS_ERASE_BACKGROUND_MODE const cstring

  @[link_name="SDL_HINT_WINDOWS_FORCE_SEMAPHORE_KERNEL"]
  const HINT_WINDOWS_FORCE_SEMAPHORE_KERNEL const cstring

  @[link_name="SDL_HINT_WINDOWS_GAMEINPUT"]
  const HINT_WINDOWS_GAMEINPUT const cstring

  @[link_name="SDL_HINT_WINDOWS_INTRESOURCE_ICON"]
  const HINT_WINDOWS_INTRESOURCE_ICON const cstring

  @[link_name="SDL_HINT_WINDOWS_INTRESOURCE_ICON_SMALL"]
  const HINT_WINDOWS_INTRESOURCE_ICON_SMALL const cstring

  @[link_name="SDL_HINT_WINDOWS_RAW_KEYBOARD"]
  const HINT_WINDOWS_RAW_KEYBOARD const cstring

  @[link_name="SDL_HINT_WINDOWS_USE_D3D9EX"]
  const HINT_WINDOWS_USE_D3D9EX const cstring

  @[link_name="SDL_HINT_WINDOW_ACTIVATE_WHEN_RAISED"]
  const HINT_WINDOW_ACTIVATE_WHEN_RAISED const cstring

  @[link_name="SDL_HINT_WINDOW_ACTIVATE_WHEN_SHOWN"]
  const HINT_WINDOW_ACTIVATE_WHEN_SHOWN const cstring

  @[link_name="SDL_HINT_WINDOW_ALLOW_TOPMOST"]
  const HINT_WINDOW_ALLOW_TOPMOST const cstring

  @[link_name="SDL_HINT_WINDOW_FRAME_USABLE_WHILE_CURSOR_HIDDEN"]
  const HINT_WINDOW_FRAME_USABLE_WHILE_CURSOR_HIDDEN const cstring

  @[link_name="SDL_HINT_X11_FORCE_OVERRIDE_REDIRECT"]
  const HINT_X11_FORCE_OVERRIDE_REDIRECT const cstring

  @[link_name="SDL_HINT_X11_WINDOW_TYPE"]
  const HINT_X11_WINDOW_TYPE const cstring

  @[link_name="SDL_HINT_X11_XCB_LIBRARY"]
  const HINT_X11_XCB_LIBRARY const cstring

  @[link_name="SDL_HINT_XINPUT_ENABLED"]
  const HINT_XINPUT_ENABLED const cstring

  @[link_name="SDL_ICONV_E2BIG"]
  const ICONV_E2BIG libc::c_ulong

  @[link_name="SDL_ICONV_EILSEQ"]
  const ICONV_EILSEQ libc::c_ulong

  @[link_name="SDL_ICONV_EINVAL"]
  const ICONV_EINVAL libc::c_ulong

  @[link_name="SDL_ICONV_ERROR"]
  const ICONV_ERROR libc::c_ulong

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

  @[link_name="SDL_PROP_APP_METADATA_COPYRIGHT_STRING"]
  const PROP_APP_METADATA_COPYRIGHT_STRING const cstring

  @[link_name="SDL_PROP_APP_METADATA_CREATOR_STRING"]
  const PROP_APP_METADATA_CREATOR_STRING const cstring

  @[link_name="SDL_PROP_APP_METADATA_IDENTIFIER_STRING"]
  const PROP_APP_METADATA_IDENTIFIER_STRING const cstring

  @[link_name="SDL_PROP_APP_METADATA_NAME_STRING"]
  const PROP_APP_METADATA_NAME_STRING const cstring

  @[link_name="SDL_PROP_APP_METADATA_TYPE_STRING"]
  const PROP_APP_METADATA_TYPE_STRING const cstring

  @[link_name="SDL_PROP_APP_METADATA_URL_STRING"]
  const PROP_APP_METADATA_URL_STRING const cstring

  @[link_name="SDL_PROP_APP_METADATA_VERSION_STRING"]
  const PROP_APP_METADATA_VERSION_STRING const cstring

  @[link_name="SDL_PROP_DISPLAY_HDR_ENABLED_BOOLEAN"]
  const PROP_DISPLAY_HDR_ENABLED_BOOLEAN const cstring

  @[link_name="SDL_PROP_DISPLAY_KMSDRM_PANEL_ORIENTATION_NUMBER"]
  const PROP_DISPLAY_KMSDRM_PANEL_ORIENTATION_NUMBER const cstring

  @[link_name="SDL_PROP_FILE_DIALOG_ACCEPT_STRING"]
  const PROP_FILE_DIALOG_ACCEPT_STRING const cstring

  @[link_name="SDL_PROP_FILE_DIALOG_CANCEL_STRING"]
  const PROP_FILE_DIALOG_CANCEL_STRING const cstring

  @[link_name="SDL_PROP_FILE_DIALOG_FILTERS_POINTER"]
  const PROP_FILE_DIALOG_FILTERS_POINTER const cstring

  @[link_name="SDL_PROP_FILE_DIALOG_LOCATION_STRING"]
  const PROP_FILE_DIALOG_LOCATION_STRING const cstring

  @[link_name="SDL_PROP_FILE_DIALOG_MANY_BOOLEAN"]
  const PROP_FILE_DIALOG_MANY_BOOLEAN const cstring

  @[link_name="SDL_PROP_FILE_DIALOG_NFILTERS_NUMBER"]
  const PROP_FILE_DIALOG_NFILTERS_NUMBER const cstring

  @[link_name="SDL_PROP_FILE_DIALOG_TITLE_STRING"]
  const PROP_FILE_DIALOG_TITLE_STRING const cstring

  @[link_name="SDL_PROP_FILE_DIALOG_WINDOW_POINTER"]
  const PROP_FILE_DIALOG_WINDOW_POINTER const cstring

  @[link_name="SDL_PROP_GAMEPAD_CAP_MONO_LED_BOOLEAN"]
  const PROP_GAMEPAD_CAP_MONO_LED_BOOLEAN const cstring

  @[link_name="SDL_PROP_GAMEPAD_CAP_PLAYER_LED_BOOLEAN"]
  const PROP_GAMEPAD_CAP_PLAYER_LED_BOOLEAN const cstring

  @[link_name="SDL_PROP_GAMEPAD_CAP_RGB_LED_BOOLEAN"]
  const PROP_GAMEPAD_CAP_RGB_LED_BOOLEAN const cstring

  @[link_name="SDL_PROP_GAMEPAD_CAP_RUMBLE_BOOLEAN"]
  const PROP_GAMEPAD_CAP_RUMBLE_BOOLEAN const cstring

  @[link_name="SDL_PROP_GAMEPAD_CAP_TRIGGER_RUMBLE_BOOLEAN"]
  const PROP_GAMEPAD_CAP_TRIGGER_RUMBLE_BOOLEAN const cstring

  @[link_name="SDL_PROP_GLOBAL_VIDEO_WAYLAND_WL_DISPLAY_POINTER"]
  const PROP_GLOBAL_VIDEO_WAYLAND_WL_DISPLAY_POINTER const cstring

  @[link_name="SDL_PROP_GPU_BUFFER_CREATE_NAME_STRING"]
  const PROP_GPU_BUFFER_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_GPU_COMPUTEPIPELINE_CREATE_NAME_STRING"]
  const PROP_GPU_COMPUTEPIPELINE_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_D3D12_SEMANTIC_NAME_STRING"]
  const PROP_GPU_DEVICE_CREATE_D3D12_SEMANTIC_NAME_STRING const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_DEBUGMODE_BOOLEAN"]
  const PROP_GPU_DEVICE_CREATE_DEBUGMODE_BOOLEAN const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_NAME_STRING"]
  const PROP_GPU_DEVICE_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_PREFERLOWPOWER_BOOLEAN"]
  const PROP_GPU_DEVICE_CREATE_PREFERLOWPOWER_BOOLEAN const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_SHADERS_DXBC_BOOLEAN"]
  const PROP_GPU_DEVICE_CREATE_SHADERS_DXBC_BOOLEAN const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_SHADERS_DXIL_BOOLEAN"]
  const PROP_GPU_DEVICE_CREATE_SHADERS_DXIL_BOOLEAN const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_SHADERS_METALLIB_BOOLEAN"]
  const PROP_GPU_DEVICE_CREATE_SHADERS_METALLIB_BOOLEAN const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_SHADERS_MSL_BOOLEAN"]
  const PROP_GPU_DEVICE_CREATE_SHADERS_MSL_BOOLEAN const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_SHADERS_PRIVATE_BOOLEAN"]
  const PROP_GPU_DEVICE_CREATE_SHADERS_PRIVATE_BOOLEAN const cstring

  @[link_name="SDL_PROP_GPU_DEVICE_CREATE_SHADERS_SPIRV_BOOLEAN"]
  const PROP_GPU_DEVICE_CREATE_SHADERS_SPIRV_BOOLEAN const cstring

  @[link_name="SDL_PROP_GPU_GRAPHICSPIPELINE_CREATE_NAME_STRING"]
  const PROP_GPU_GRAPHICSPIPELINE_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_GPU_SAMPLER_CREATE_NAME_STRING"]
  const PROP_GPU_SAMPLER_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_GPU_SHADER_CREATE_NAME_STRING"]
  const PROP_GPU_SHADER_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_A_FLOAT"]
  const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_A_FLOAT const cstring

  @[link_name="SDL_PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_B_FLOAT"]
  const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_B_FLOAT const cstring

  @[link_name="SDL_PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_DEPTH_FLOAT"]
  const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_DEPTH_FLOAT const cstring

  @[link_name="SDL_PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_G_FLOAT"]
  const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_G_FLOAT const cstring

  @[link_name="SDL_PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_R_FLOAT"]
  const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_R_FLOAT const cstring

  @[link_name="SDL_PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_STENCIL_NUMBER"]
  const PROP_GPU_TEXTURE_CREATE_D3D12_CLEAR_STENCIL_NUMBER const cstring

  @[link_name="SDL_PROP_GPU_TEXTURE_CREATE_NAME_STRING"]
  const PROP_GPU_TEXTURE_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_GPU_TRANSFERBUFFER_CREATE_NAME_STRING"]
  const PROP_GPU_TRANSFERBUFFER_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_IOSTREAM_ANDROID_AASSET_POINTER"]
  const PROP_IOSTREAM_ANDROID_AASSET_POINTER const cstring

  @[link_name="SDL_PROP_IOSTREAM_DYNAMIC_CHUNKSIZE_NUMBER"]
  const PROP_IOSTREAM_DYNAMIC_CHUNKSIZE_NUMBER const cstring

  @[link_name="SDL_PROP_IOSTREAM_DYNAMIC_MEMORY_POINTER"]
  const PROP_IOSTREAM_DYNAMIC_MEMORY_POINTER const cstring

  @[link_name="SDL_PROP_IOSTREAM_FILE_DESCRIPTOR_NUMBER"]
  const PROP_IOSTREAM_FILE_DESCRIPTOR_NUMBER const cstring

  @[link_name="SDL_PROP_IOSTREAM_MEMORY_POINTER"]
  const PROP_IOSTREAM_MEMORY_POINTER const cstring

  @[link_name="SDL_PROP_IOSTREAM_MEMORY_SIZE_NUMBER"]
  const PROP_IOSTREAM_MEMORY_SIZE_NUMBER const cstring

  @[link_name="SDL_PROP_IOSTREAM_STDIO_FILE_POINTER"]
  const PROP_IOSTREAM_STDIO_FILE_POINTER const cstring

  @[link_name="SDL_PROP_IOSTREAM_WINDOWS_HANDLE_POINTER"]
  const PROP_IOSTREAM_WINDOWS_HANDLE_POINTER const cstring

  @[link_name="SDL_PROP_JOYSTICK_CAP_MONO_LED_BOOLEAN"]
  const PROP_JOYSTICK_CAP_MONO_LED_BOOLEAN const cstring

  @[link_name="SDL_PROP_JOYSTICK_CAP_PLAYER_LED_BOOLEAN"]
  const PROP_JOYSTICK_CAP_PLAYER_LED_BOOLEAN const cstring

  @[link_name="SDL_PROP_JOYSTICK_CAP_RGB_LED_BOOLEAN"]
  const PROP_JOYSTICK_CAP_RGB_LED_BOOLEAN const cstring

  @[link_name="SDL_PROP_JOYSTICK_CAP_RUMBLE_BOOLEAN"]
  const PROP_JOYSTICK_CAP_RUMBLE_BOOLEAN const cstring

  @[link_name="SDL_PROP_JOYSTICK_CAP_TRIGGER_RUMBLE_BOOLEAN"]
  const PROP_JOYSTICK_CAP_TRIGGER_RUMBLE_BOOLEAN const cstring

  @[link_name="SDL_PROP_PROCESS_BACKGROUND_BOOLEAN"]
  const PROP_PROCESS_BACKGROUND_BOOLEAN const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_ARGS_POINTER"]
  const PROP_PROCESS_CREATE_ARGS_POINTER const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_BACKGROUND_BOOLEAN"]
  const PROP_PROCESS_CREATE_BACKGROUND_BOOLEAN const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_ENVIRONMENT_POINTER"]
  const PROP_PROCESS_CREATE_ENVIRONMENT_POINTER const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_STDERR_NUMBER"]
  const PROP_PROCESS_CREATE_STDERR_NUMBER const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_STDERR_POINTER"]
  const PROP_PROCESS_CREATE_STDERR_POINTER const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_STDERR_TO_STDOUT_BOOLEAN"]
  const PROP_PROCESS_CREATE_STDERR_TO_STDOUT_BOOLEAN const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_STDIN_NUMBER"]
  const PROP_PROCESS_CREATE_STDIN_NUMBER const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_STDIN_POINTER"]
  const PROP_PROCESS_CREATE_STDIN_POINTER const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_STDOUT_NUMBER"]
  const PROP_PROCESS_CREATE_STDOUT_NUMBER const cstring

  @[link_name="SDL_PROP_PROCESS_CREATE_STDOUT_POINTER"]
  const PROP_PROCESS_CREATE_STDOUT_POINTER const cstring

  @[link_name="SDL_PROP_PROCESS_PID_NUMBER"]
  const PROP_PROCESS_PID_NUMBER const cstring

  @[link_name="SDL_PROP_PROCESS_STDERR_POINTER"]
  const PROP_PROCESS_STDERR_POINTER const cstring

  @[link_name="SDL_PROP_PROCESS_STDIN_POINTER"]
  const PROP_PROCESS_STDIN_POINTER const cstring

  @[link_name="SDL_PROP_PROCESS_STDOUT_POINTER"]
  const PROP_PROCESS_STDOUT_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_NAME_STRING"]
  const PROP_RENDERER_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_OUTPUT_COLORSPACE_NUMBER"]
  const PROP_RENDERER_CREATE_OUTPUT_COLORSPACE_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_PRESENT_VSYNC_NUMBER"]
  const PROP_RENDERER_CREATE_PRESENT_VSYNC_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_SURFACE_POINTER"]
  const PROP_RENDERER_CREATE_SURFACE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_VULKAN_DEVICE_POINTER"]
  const PROP_RENDERER_CREATE_VULKAN_DEVICE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_VULKAN_GRAPHICS_QUEUE_FAMILY_INDEX_NUMBER"]
  const PROP_RENDERER_CREATE_VULKAN_GRAPHICS_QUEUE_FAMILY_INDEX_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_VULKAN_INSTANCE_POINTER"]
  const PROP_RENDERER_CREATE_VULKAN_INSTANCE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_VULKAN_PHYSICAL_DEVICE_POINTER"]
  const PROP_RENDERER_CREATE_VULKAN_PHYSICAL_DEVICE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_VULKAN_PRESENT_QUEUE_FAMILY_INDEX_NUMBER"]
  const PROP_RENDERER_CREATE_VULKAN_PRESENT_QUEUE_FAMILY_INDEX_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_VULKAN_SURFACE_NUMBER"]
  const PROP_RENDERER_CREATE_VULKAN_SURFACE_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_CREATE_WINDOW_POINTER"]
  const PROP_RENDERER_CREATE_WINDOW_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_D3D11_DEVICE_POINTER"]
  const PROP_RENDERER_D3D11_DEVICE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_D3D11_SWAPCHAIN_POINTER"]
  const PROP_RENDERER_D3D11_SWAPCHAIN_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_D3D12_COMMAND_QUEUE_POINTER"]
  const PROP_RENDERER_D3D12_COMMAND_QUEUE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_D3D12_DEVICE_POINTER"]
  const PROP_RENDERER_D3D12_DEVICE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_D3D12_SWAPCHAIN_POINTER"]
  const PROP_RENDERER_D3D12_SWAPCHAIN_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_D3D9_DEVICE_POINTER"]
  const PROP_RENDERER_D3D9_DEVICE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_GPU_DEVICE_POINTER"]
  const PROP_RENDERER_GPU_DEVICE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_HDR_ENABLED_BOOLEAN"]
  const PROP_RENDERER_HDR_ENABLED_BOOLEAN const cstring

  @[link_name="SDL_PROP_RENDERER_HDR_HEADROOM_FLOAT"]
  const PROP_RENDERER_HDR_HEADROOM_FLOAT const cstring

  @[link_name="SDL_PROP_RENDERER_MAX_TEXTURE_SIZE_NUMBER"]
  const PROP_RENDERER_MAX_TEXTURE_SIZE_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_NAME_STRING"]
  const PROP_RENDERER_NAME_STRING const cstring

  @[link_name="SDL_PROP_RENDERER_OUTPUT_COLORSPACE_NUMBER"]
  const PROP_RENDERER_OUTPUT_COLORSPACE_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_SDR_WHITE_POINT_FLOAT"]
  const PROP_RENDERER_SDR_WHITE_POINT_FLOAT const cstring

  @[link_name="SDL_PROP_RENDERER_SURFACE_POINTER"]
  const PROP_RENDERER_SURFACE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_TEXTURE_FORMATS_POINTER"]
  const PROP_RENDERER_TEXTURE_FORMATS_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_VSYNC_NUMBER"]
  const PROP_RENDERER_VSYNC_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_VULKAN_DEVICE_POINTER"]
  const PROP_RENDERER_VULKAN_DEVICE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_VULKAN_GRAPHICS_QUEUE_FAMILY_INDEX_NUMBER"]
  const PROP_RENDERER_VULKAN_GRAPHICS_QUEUE_FAMILY_INDEX_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_VULKAN_INSTANCE_POINTER"]
  const PROP_RENDERER_VULKAN_INSTANCE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_VULKAN_PHYSICAL_DEVICE_POINTER"]
  const PROP_RENDERER_VULKAN_PHYSICAL_DEVICE_POINTER const cstring

  @[link_name="SDL_PROP_RENDERER_VULKAN_PRESENT_QUEUE_FAMILY_INDEX_NUMBER"]
  const PROP_RENDERER_VULKAN_PRESENT_QUEUE_FAMILY_INDEX_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_VULKAN_SURFACE_NUMBER"]
  const PROP_RENDERER_VULKAN_SURFACE_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_VULKAN_SWAPCHAIN_IMAGE_COUNT_NUMBER"]
  const PROP_RENDERER_VULKAN_SWAPCHAIN_IMAGE_COUNT_NUMBER const cstring

  @[link_name="SDL_PROP_RENDERER_WINDOW_POINTER"]
  const PROP_RENDERER_WINDOW_POINTER const cstring

  @[link_name="SDL_PROP_SURFACE_HDR_HEADROOM_FLOAT"]
  const PROP_SURFACE_HDR_HEADROOM_FLOAT const cstring

  @[link_name="SDL_PROP_SURFACE_HOTSPOT_X_NUMBER"]
  const PROP_SURFACE_HOTSPOT_X_NUMBER const cstring

  @[link_name="SDL_PROP_SURFACE_HOTSPOT_Y_NUMBER"]
  const PROP_SURFACE_HOTSPOT_Y_NUMBER const cstring

  @[link_name="SDL_PROP_SURFACE_SDR_WHITE_POINT_FLOAT"]
  const PROP_SURFACE_SDR_WHITE_POINT_FLOAT const cstring

  @[link_name="SDL_PROP_SURFACE_TONEMAP_OPERATOR_STRING"]
  const PROP_SURFACE_TONEMAP_OPERATOR_STRING const cstring

  @[link_name="SDL_PROP_TEXTINPUT_ANDROID_INPUTTYPE_NUMBER"]
  const PROP_TEXTINPUT_ANDROID_INPUTTYPE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTINPUT_AUTOCORRECT_BOOLEAN"]
  const PROP_TEXTINPUT_AUTOCORRECT_BOOLEAN const cstring

  @[link_name="SDL_PROP_TEXTINPUT_CAPITALIZATION_NUMBER"]
  const PROP_TEXTINPUT_CAPITALIZATION_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTINPUT_MULTILINE_BOOLEAN"]
  const PROP_TEXTINPUT_MULTILINE_BOOLEAN const cstring

  @[link_name="SDL_PROP_TEXTINPUT_TYPE_NUMBER"]
  const PROP_TEXTINPUT_TYPE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_ACCESS_NUMBER"]
  const PROP_TEXTURE_ACCESS_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_COLORSPACE_NUMBER"]
  const PROP_TEXTURE_COLORSPACE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_ACCESS_NUMBER"]
  const PROP_TEXTURE_CREATE_ACCESS_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_COLORSPACE_NUMBER"]
  const PROP_TEXTURE_CREATE_COLORSPACE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_D3D11_TEXTURE_POINTER"]
  const PROP_TEXTURE_CREATE_D3D11_TEXTURE_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_D3D11_TEXTURE_U_POINTER"]
  const PROP_TEXTURE_CREATE_D3D11_TEXTURE_U_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_D3D11_TEXTURE_V_POINTER"]
  const PROP_TEXTURE_CREATE_D3D11_TEXTURE_V_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_D3D12_TEXTURE_POINTER"]
  const PROP_TEXTURE_CREATE_D3D12_TEXTURE_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_D3D12_TEXTURE_U_POINTER"]
  const PROP_TEXTURE_CREATE_D3D12_TEXTURE_U_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_D3D12_TEXTURE_V_POINTER"]
  const PROP_TEXTURE_CREATE_D3D12_TEXTURE_V_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_FORMAT_NUMBER"]
  const PROP_TEXTURE_CREATE_FORMAT_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_HDR_HEADROOM_FLOAT"]
  const PROP_TEXTURE_CREATE_HDR_HEADROOM_FLOAT const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_HEIGHT_NUMBER"]
  const PROP_TEXTURE_CREATE_HEIGHT_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_METAL_PIXELBUFFER_POINTER"]
  const PROP_TEXTURE_CREATE_METAL_PIXELBUFFER_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_NUMBER"]
  const PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_UV_NUMBER"]
  const PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_UV_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_U_NUMBER"]
  const PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_U_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_V_NUMBER"]
  const PROP_TEXTURE_CREATE_OPENGLES2_TEXTURE_V_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_OPENGL_TEXTURE_NUMBER"]
  const PROP_TEXTURE_CREATE_OPENGL_TEXTURE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_OPENGL_TEXTURE_UV_NUMBER"]
  const PROP_TEXTURE_CREATE_OPENGL_TEXTURE_UV_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_OPENGL_TEXTURE_U_NUMBER"]
  const PROP_TEXTURE_CREATE_OPENGL_TEXTURE_U_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_OPENGL_TEXTURE_V_NUMBER"]
  const PROP_TEXTURE_CREATE_OPENGL_TEXTURE_V_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_SDR_WHITE_POINT_FLOAT"]
  const PROP_TEXTURE_CREATE_SDR_WHITE_POINT_FLOAT const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_VULKAN_TEXTURE_NUMBER"]
  const PROP_TEXTURE_CREATE_VULKAN_TEXTURE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_CREATE_WIDTH_NUMBER"]
  const PROP_TEXTURE_CREATE_WIDTH_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_D3D11_TEXTURE_POINTER"]
  const PROP_TEXTURE_D3D11_TEXTURE_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_D3D11_TEXTURE_U_POINTER"]
  const PROP_TEXTURE_D3D11_TEXTURE_U_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_D3D11_TEXTURE_V_POINTER"]
  const PROP_TEXTURE_D3D11_TEXTURE_V_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_D3D12_TEXTURE_POINTER"]
  const PROP_TEXTURE_D3D12_TEXTURE_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_D3D12_TEXTURE_U_POINTER"]
  const PROP_TEXTURE_D3D12_TEXTURE_U_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_D3D12_TEXTURE_V_POINTER"]
  const PROP_TEXTURE_D3D12_TEXTURE_V_POINTER const cstring

  @[link_name="SDL_PROP_TEXTURE_FORMAT_NUMBER"]
  const PROP_TEXTURE_FORMAT_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_HDR_HEADROOM_FLOAT"]
  const PROP_TEXTURE_HDR_HEADROOM_FLOAT const cstring

  @[link_name="SDL_PROP_TEXTURE_HEIGHT_NUMBER"]
  const PROP_TEXTURE_HEIGHT_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGLES2_TEXTURE_NUMBER"]
  const PROP_TEXTURE_OPENGLES2_TEXTURE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGLES2_TEXTURE_TARGET_NUMBER"]
  const PROP_TEXTURE_OPENGLES2_TEXTURE_TARGET_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGLES2_TEXTURE_UV_NUMBER"]
  const PROP_TEXTURE_OPENGLES2_TEXTURE_UV_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGLES2_TEXTURE_U_NUMBER"]
  const PROP_TEXTURE_OPENGLES2_TEXTURE_U_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGLES2_TEXTURE_V_NUMBER"]
  const PROP_TEXTURE_OPENGLES2_TEXTURE_V_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGL_TEXTURE_NUMBER"]
  const PROP_TEXTURE_OPENGL_TEXTURE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGL_TEXTURE_TARGET_NUMBER"]
  const PROP_TEXTURE_OPENGL_TEXTURE_TARGET_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGL_TEXTURE_UV_NUMBER"]
  const PROP_TEXTURE_OPENGL_TEXTURE_UV_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGL_TEXTURE_U_NUMBER"]
  const PROP_TEXTURE_OPENGL_TEXTURE_U_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGL_TEXTURE_V_NUMBER"]
  const PROP_TEXTURE_OPENGL_TEXTURE_V_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGL_TEX_H_FLOAT"]
  const PROP_TEXTURE_OPENGL_TEX_H_FLOAT const cstring

  @[link_name="SDL_PROP_TEXTURE_OPENGL_TEX_W_FLOAT"]
  const PROP_TEXTURE_OPENGL_TEX_W_FLOAT const cstring

  @[link_name="SDL_PROP_TEXTURE_SDR_WHITE_POINT_FLOAT"]
  const PROP_TEXTURE_SDR_WHITE_POINT_FLOAT const cstring

  @[link_name="SDL_PROP_TEXTURE_VULKAN_TEXTURE_NUMBER"]
  const PROP_TEXTURE_VULKAN_TEXTURE_NUMBER const cstring

  @[link_name="SDL_PROP_TEXTURE_WIDTH_NUMBER"]
  const PROP_TEXTURE_WIDTH_NUMBER const cstring

  @[link_name="SDL_PROP_THREAD_CREATE_ENTRY_FUNCTION_POINTER"]
  const PROP_THREAD_CREATE_ENTRY_FUNCTION_POINTER const cstring

  @[link_name="SDL_PROP_THREAD_CREATE_NAME_STRING"]
  const PROP_THREAD_CREATE_NAME_STRING const cstring

  @[link_name="SDL_PROP_THREAD_CREATE_STACKSIZE_NUMBER"]
  const PROP_THREAD_CREATE_STACKSIZE_NUMBER const cstring

  @[link_name="SDL_PROP_THREAD_CREATE_USERDATA_POINTER"]
  const PROP_THREAD_CREATE_USERDATA_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_ANDROID_SURFACE_POINTER"]
  const PROP_WINDOW_ANDROID_SURFACE_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_ANDROID_WINDOW_POINTER"]
  const PROP_WINDOW_ANDROID_WINDOW_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_COCOA_METAL_VIEW_TAG_NUMBER"]
  const PROP_WINDOW_COCOA_METAL_VIEW_TAG_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_COCOA_WINDOW_POINTER"]
  const PROP_WINDOW_COCOA_WINDOW_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_ALWAYS_ON_TOP_BOOLEAN"]
  const PROP_WINDOW_CREATE_ALWAYS_ON_TOP_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_BORDERLESS_BOOLEAN"]
  const PROP_WINDOW_CREATE_BORDERLESS_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_COCOA_VIEW_POINTER"]
  const PROP_WINDOW_CREATE_COCOA_VIEW_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_COCOA_WINDOW_POINTER"]
  const PROP_WINDOW_CREATE_COCOA_WINDOW_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_CONSTRAIN_POPUP_BOOLEAN"]
  const PROP_WINDOW_CREATE_CONSTRAIN_POPUP_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_EXTERNAL_GRAPHICS_CONTEXT_BOOLEAN"]
  const PROP_WINDOW_CREATE_EXTERNAL_GRAPHICS_CONTEXT_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_FLAGS_NUMBER"]
  const PROP_WINDOW_CREATE_FLAGS_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_FOCUSABLE_BOOLEAN"]
  const PROP_WINDOW_CREATE_FOCUSABLE_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_FULLSCREEN_BOOLEAN"]
  const PROP_WINDOW_CREATE_FULLSCREEN_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_HEIGHT_NUMBER"]
  const PROP_WINDOW_CREATE_HEIGHT_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_HIDDEN_BOOLEAN"]
  const PROP_WINDOW_CREATE_HIDDEN_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_HIGH_PIXEL_DENSITY_BOOLEAN"]
  const PROP_WINDOW_CREATE_HIGH_PIXEL_DENSITY_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_MAXIMIZED_BOOLEAN"]
  const PROP_WINDOW_CREATE_MAXIMIZED_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_MENU_BOOLEAN"]
  const PROP_WINDOW_CREATE_MENU_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_METAL_BOOLEAN"]
  const PROP_WINDOW_CREATE_METAL_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_MINIMIZED_BOOLEAN"]
  const PROP_WINDOW_CREATE_MINIMIZED_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_MODAL_BOOLEAN"]
  const PROP_WINDOW_CREATE_MODAL_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_MOUSE_GRABBED_BOOLEAN"]
  const PROP_WINDOW_CREATE_MOUSE_GRABBED_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_OPENGL_BOOLEAN"]
  const PROP_WINDOW_CREATE_OPENGL_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_PARENT_POINTER"]
  const PROP_WINDOW_CREATE_PARENT_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_RESIZABLE_BOOLEAN"]
  const PROP_WINDOW_CREATE_RESIZABLE_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_TITLE_STRING"]
  const PROP_WINDOW_CREATE_TITLE_STRING const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_TOOLTIP_BOOLEAN"]
  const PROP_WINDOW_CREATE_TOOLTIP_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_TRANSPARENT_BOOLEAN"]
  const PROP_WINDOW_CREATE_TRANSPARENT_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_UTILITY_BOOLEAN"]
  const PROP_WINDOW_CREATE_UTILITY_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_VULKAN_BOOLEAN"]
  const PROP_WINDOW_CREATE_VULKAN_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_WAYLAND_CREATE_EGL_WINDOW_BOOLEAN"]
  const PROP_WINDOW_CREATE_WAYLAND_CREATE_EGL_WINDOW_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_WAYLAND_SURFACE_ROLE_CUSTOM_BOOLEAN"]
  const PROP_WINDOW_CREATE_WAYLAND_SURFACE_ROLE_CUSTOM_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_WAYLAND_WL_SURFACE_POINTER"]
  const PROP_WINDOW_CREATE_WAYLAND_WL_SURFACE_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_WIDTH_NUMBER"]
  const PROP_WINDOW_CREATE_WIDTH_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_WIN32_HWND_POINTER"]
  const PROP_WINDOW_CREATE_WIN32_HWND_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_WIN32_PIXEL_FORMAT_HWND_POINTER"]
  const PROP_WINDOW_CREATE_WIN32_PIXEL_FORMAT_HWND_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_X11_WINDOW_NUMBER"]
  const PROP_WINDOW_CREATE_X11_WINDOW_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_X_NUMBER"]
  const PROP_WINDOW_CREATE_X_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_CREATE_Y_NUMBER"]
  const PROP_WINDOW_CREATE_Y_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_HDR_ENABLED_BOOLEAN"]
  const PROP_WINDOW_HDR_ENABLED_BOOLEAN const cstring

  @[link_name="SDL_PROP_WINDOW_HDR_HEADROOM_FLOAT"]
  const PROP_WINDOW_HDR_HEADROOM_FLOAT const cstring

  @[link_name="SDL_PROP_WINDOW_KMSDRM_DEVICE_INDEX_NUMBER"]
  const PROP_WINDOW_KMSDRM_DEVICE_INDEX_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_KMSDRM_DRM_FD_NUMBER"]
  const PROP_WINDOW_KMSDRM_DRM_FD_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_KMSDRM_GBM_DEVICE_POINTER"]
  const PROP_WINDOW_KMSDRM_GBM_DEVICE_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_OPENVR_OVERLAY_ID"]
  const PROP_WINDOW_OPENVR_OVERLAY_ID const cstring

  @[link_name="SDL_PROP_WINDOW_SDR_WHITE_LEVEL_FLOAT"]
  const PROP_WINDOW_SDR_WHITE_LEVEL_FLOAT const cstring

  @[link_name="SDL_PROP_WINDOW_SHAPE_POINTER"]
  const PROP_WINDOW_SHAPE_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_UIKIT_METAL_VIEW_TAG_NUMBER"]
  const PROP_WINDOW_UIKIT_METAL_VIEW_TAG_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_UIKIT_OPENGL_FRAMEBUFFER_NUMBER"]
  const PROP_WINDOW_UIKIT_OPENGL_FRAMEBUFFER_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_UIKIT_OPENGL_RENDERBUFFER_NUMBER"]
  const PROP_WINDOW_UIKIT_OPENGL_RENDERBUFFER_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_UIKIT_OPENGL_RESOLVE_FRAMEBUFFER_NUMBER"]
  const PROP_WINDOW_UIKIT_OPENGL_RESOLVE_FRAMEBUFFER_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_UIKIT_WINDOW_POINTER"]
  const PROP_WINDOW_UIKIT_WINDOW_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_VIVANTE_DISPLAY_POINTER"]
  const PROP_WINDOW_VIVANTE_DISPLAY_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_VIVANTE_SURFACE_POINTER"]
  const PROP_WINDOW_VIVANTE_SURFACE_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_VIVANTE_WINDOW_POINTER"]
  const PROP_WINDOW_VIVANTE_WINDOW_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WAYLAND_DISPLAY_POINTER"]
  const PROP_WINDOW_WAYLAND_DISPLAY_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WAYLAND_EGL_WINDOW_POINTER"]
  const PROP_WINDOW_WAYLAND_EGL_WINDOW_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WAYLAND_SURFACE_POINTER"]
  const PROP_WINDOW_WAYLAND_SURFACE_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WAYLAND_VIEWPORT_POINTER"]
  const PROP_WINDOW_WAYLAND_VIEWPORT_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WAYLAND_XDG_POPUP_POINTER"]
  const PROP_WINDOW_WAYLAND_XDG_POPUP_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WAYLAND_XDG_POSITIONER_POINTER"]
  const PROP_WINDOW_WAYLAND_XDG_POSITIONER_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WAYLAND_XDG_SURFACE_POINTER"]
  const PROP_WINDOW_WAYLAND_XDG_SURFACE_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WAYLAND_XDG_TOPLEVEL_EXPORT_HANDLE_STRING"]
  const PROP_WINDOW_WAYLAND_XDG_TOPLEVEL_EXPORT_HANDLE_STRING const cstring

  @[link_name="SDL_PROP_WINDOW_WAYLAND_XDG_TOPLEVEL_POINTER"]
  const PROP_WINDOW_WAYLAND_XDG_TOPLEVEL_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WIN32_HDC_POINTER"]
  const PROP_WINDOW_WIN32_HDC_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WIN32_HWND_POINTER"]
  const PROP_WINDOW_WIN32_HWND_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_WIN32_INSTANCE_POINTER"]
  const PROP_WINDOW_WIN32_INSTANCE_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_X11_DISPLAY_POINTER"]
  const PROP_WINDOW_X11_DISPLAY_POINTER const cstring

  @[link_name="SDL_PROP_WINDOW_X11_SCREEN_NUMBER"]
  const PROP_WINDOW_X11_SCREEN_NUMBER const cstring

  @[link_name="SDL_PROP_WINDOW_X11_WINDOW_NUMBER"]
  const PROP_WINDOW_X11_WINDOW_NUMBER const cstring

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

  @[link_name="SDL_malloc"]
  fn malloc(size libc::c_size_t) rawptr

  @[link_name="SDL_calloc"]
  fn calloc(nmemb libc::c_size_t, size libc::c_size_t) rawptr

  @[link_name="SDL_realloc"]
  fn realloc(mem rawptr, size libc::c_size_t) rawptr

  @[link_name="SDL_free"]
  fn free(mem rawptr)

  @[link_name="SDL_GetOriginalMemoryFunctions"]
  fn get_original_memory_functions(malloc_func *malloc_func, calloc_func *calloc_func, realloc_func *realloc_func, free_func *free_func)

  @[link_name="SDL_GetMemoryFunctions"]
  fn get_memory_functions(malloc_func *malloc_func, calloc_func *calloc_func, realloc_func *realloc_func, free_func *free_func)

  @[link_name="SDL_SetMemoryFunctions"]
  fn set_memory_functions(malloc_func malloc_func, calloc_func calloc_func, realloc_func realloc_func, free_func free_func) bool

  @[link_name="SDL_aligned_alloc"]
  fn aligned_alloc(alignment libc::c_size_t, size libc::c_size_t) rawptr

  @[link_name="SDL_aligned_free"]
  fn aligned_free(mem rawptr)

  @[link_name="SDL_GetNumAllocations"]
  fn get_num_allocations() libc::c_int

  @[link_name="SDL_GetEnvironment"]
  fn get_environment() *Environment

  @[link_name="SDL_CreateEnvironment"]
  fn create_environment(populated bool) *Environment

  @[link_name="SDL_GetEnvironmentVariable"]
  fn get_environment_variable(env *Environment, name const cstring) const cstring

  @[link_name="SDL_GetEnvironmentVariables"]
  fn get_environment_variables(env *Environment) *cstring

  @[link_name="SDL_SetEnvironmentVariable"]
  fn set_environment_variable(env *Environment, name const cstring, value const cstring, overwrite bool) bool

  @[link_name="SDL_UnsetEnvironmentVariable"]
  fn unset_environment_variable(env *Environment, name const cstring) bool

  @[link_name="SDL_DestroyEnvironment"]
  fn destroy_environment(env *Environment)

  @[link_name="SDL_getenv"]
  fn getenv(name const cstring) const cstring

  @[link_name="SDL_getenv_unsafe"]
  fn getenv_unsafe(name const cstring) const cstring

  @[link_name="SDL_setenv_unsafe"]
  fn setenv_unsafe(name const cstring, value const cstring, overwrite libc::c_int) libc::c_int

  @[link_name="SDL_unsetenv_unsafe"]
  fn unsetenv_unsafe(name const cstring) libc::c_int

  @[link_name="SDL_qsort"]
  fn qsort(base rawptr, nmemb libc::c_size_t, size libc::c_size_t, compare CompareCallback)

  @[link_name="SDL_bsearch"]
  fn bsearch(key const *Void, base const *Void, nmemb libc::c_size_t, size libc::c_size_t, compare CompareCallback) rawptr

  @[link_name="SDL_qsort_r"]
  fn qsort_r(base rawptr, nmemb libc::c_size_t, size libc::c_size_t, compare CompareCallback_r, userdata rawptr)

  @[link_name="SDL_bsearch_r"]
  fn bsearch_r(key const *Void, base const *Void, nmemb libc::c_size_t, size libc::c_size_t, compare CompareCallback_r, userdata rawptr) rawptr

  @[link_name="SDL_abs"]
  fn abs(x libc::c_int) libc::c_int

  @[link_name="SDL_isalpha"]
  fn isalpha(x libc::c_int) libc::c_int

  @[link_name="SDL_isalnum"]
  fn isalnum(x libc::c_int) libc::c_int

  @[link_name="SDL_isblank"]
  fn isblank(x libc::c_int) libc::c_int

  @[link_name="SDL_iscntrl"]
  fn iscntrl(x libc::c_int) libc::c_int

  @[link_name="SDL_isdigit"]
  fn isdigit(x libc::c_int) libc::c_int

  @[link_name="SDL_isxdigit"]
  fn isxdigit(x libc::c_int) libc::c_int

  @[link_name="SDL_ispunct"]
  fn ispunct(x libc::c_int) libc::c_int

  @[link_name="SDL_isspace"]
  fn isspace(x libc::c_int) libc::c_int

  @[link_name="SDL_isupper"]
  fn isupper(x libc::c_int) libc::c_int

  @[link_name="SDL_islower"]
  fn islower(x libc::c_int) libc::c_int

  @[link_name="SDL_isprint"]
  fn isprint(x libc::c_int) libc::c_int

  @[link_name="SDL_isgraph"]
  fn isgraph(x libc::c_int) libc::c_int

  @[link_name="SDL_toupper"]
  fn toupper(x libc::c_int) libc::c_int

  @[link_name="SDL_tolower"]
  fn tolower(x libc::c_int) libc::c_int

  @[link_name="SDL_crc16"]
  fn crc16(crc u16, data const *Void, len libc::c_size_t) u16

  @[link_name="SDL_crc32"]
  fn crc32(crc u32, data const *Void, len libc::c_size_t) u32

  @[link_name="SDL_murmur3_32"]
  fn murmur3_32(data const *Void, len libc::c_size_t, seed u32) u32

  @[link_name="SDL_memcpy"]
  fn memcpy(dst rawptr, src const *Void, len libc::c_size_t) rawptr

  @[link_name="SDL_memmove"]
  fn memmove(dst rawptr, src const *Void, len libc::c_size_t) rawptr

  @[link_name="SDL_memset"]
  fn memset(dst rawptr, c libc::c_int, len libc::c_size_t) rawptr

  @[link_name="SDL_memset4"]
  fn memset4(dst rawptr, val u32, dwords libc::c_size_t) rawptr

  @[link_name="SDL_memcmp"]
  fn memcmp(s1 const *Void, s2 const *Void, len libc::c_size_t) libc::c_int

  @[link_name="SDL_wcslen"]
  fn wcslen(wstr const *libc::c_wchar_t) libc::c_size_t

  @[link_name="SDL_wcsnlen"]
  fn wcsnlen(wstr const *libc::c_wchar_t, maxlen libc::c_size_t) libc::c_size_t

  @[link_name="SDL_wcslcpy"]
  fn wcslcpy(dst *libc::c_wchar_t, src const *libc::c_wchar_t, maxlen libc::c_size_t) libc::c_size_t

  @[link_name="SDL_wcslcat"]
  fn wcslcat(dst *libc::c_wchar_t, src const *libc::c_wchar_t, maxlen libc::c_size_t) libc::c_size_t

  @[link_name="SDL_wcsdup"]
  fn wcsdup(wstr const *libc::c_wchar_t) *libc::c_wchar_t

  @[link_name="SDL_wcsstr"]
  fn wcsstr(haystack const *libc::c_wchar_t, needle const *libc::c_wchar_t) *libc::c_wchar_t

  @[link_name="SDL_wcsnstr"]
  fn wcsnstr(haystack const *libc::c_wchar_t, needle const *libc::c_wchar_t, maxlen libc::c_size_t) *libc::c_wchar_t

  @[link_name="SDL_wcscmp"]
  fn wcscmp(str1 const *libc::c_wchar_t, str2 const *libc::c_wchar_t) libc::c_int

  @[link_name="SDL_wcsncmp"]
  fn wcsncmp(str1 const *libc::c_wchar_t, str2 const *libc::c_wchar_t, maxlen libc::c_size_t) libc::c_int

  @[link_name="SDL_wcscasecmp"]
  fn wcscasecmp(str1 const *libc::c_wchar_t, str2 const *libc::c_wchar_t) libc::c_int

  @[link_name="SDL_wcsncasecmp"]
  fn wcsncasecmp(str1 const *libc::c_wchar_t, str2 const *libc::c_wchar_t, maxlen libc::c_size_t) libc::c_int

  @[link_name="SDL_wcstol"]
  fn wcstol(str const *libc::c_wchar_t, endp **libc::c_wchar_t, base libc::c_int) libc::c_long

  @[link_name="SDL_strlen"]
  fn strlen(str const cstring) libc::c_size_t

  @[link_name="SDL_strnlen"]
  fn strnlen(str const cstring, maxlen libc::c_size_t) libc::c_size_t

  @[link_name="SDL_strlcpy"]
  fn strlcpy(dst cstring, src const cstring, maxlen libc::c_size_t) libc::c_size_t

  @[link_name="SDL_utf8strlcpy"]
  fn utf8strlcpy(dst cstring, src const cstring, dst_bytes libc::c_size_t) libc::c_size_t

  @[link_name="SDL_strlcat"]
  fn strlcat(dst cstring, src const cstring, maxlen libc::c_size_t) libc::c_size_t

  @[link_name="SDL_strdup"]
  fn strdup(str const cstring) cstring

  @[link_name="SDL_strndup"]
  fn strndup(str const cstring, maxlen libc::c_size_t) cstring

  @[link_name="SDL_strrev"]
  fn strrev(str cstring) cstring

  @[link_name="SDL_strupr"]
  fn strupr(str cstring) cstring

  @[link_name="SDL_strlwr"]
  fn strlwr(str cstring) cstring

  @[link_name="SDL_strchr"]
  fn strchr(str const cstring, c libc::c_int) cstring

  @[link_name="SDL_strrchr"]
  fn strrchr(str const cstring, c libc::c_int) cstring

  @[link_name="SDL_strstr"]
  fn strstr(haystack const cstring, needle const cstring) cstring

  @[link_name="SDL_strnstr"]
  fn strnstr(haystack const cstring, needle const cstring, maxlen libc::c_size_t) cstring

  @[link_name="SDL_strcasestr"]
  fn strcasestr(haystack const cstring, needle const cstring) cstring

  @[link_name="SDL_strtok_r"]
  fn strtok_r(str cstring, delim const cstring, saveptr *cstring) cstring

  @[link_name="SDL_utf8strlen"]
  fn utf8strlen(str const cstring) libc::c_size_t

  @[link_name="SDL_utf8strnlen"]
  fn utf8strnlen(str const cstring, bytes libc::c_size_t) libc::c_size_t

  @[link_name="SDL_itoa"]
  fn itoa(value libc::c_int, str cstring, radix libc::c_int) cstring

  @[link_name="SDL_uitoa"]
  fn uitoa(value libc::c_uint, str cstring, radix libc::c_int) cstring

  @[link_name="SDL_ltoa"]
  fn ltoa(value libc::c_long, str cstring, radix libc::c_int) cstring

  @[link_name="SDL_ultoa"]
  fn ultoa(value libc::c_ulong, str cstring, radix libc::c_int) cstring

  @[link_name="SDL_lltoa"]
  fn lltoa(value libc::c_longlong, str cstring, radix libc::c_int) cstring

  @[link_name="SDL_ulltoa"]
  fn ulltoa(value libc::c_ulonglong, str cstring, radix libc::c_int) cstring

  @[link_name="SDL_atoi"]
  fn atoi(str const cstring) libc::c_int

  @[link_name="SDL_atof"]
  fn atof(str const cstring) libc::c_double

  @[link_name="SDL_strtol"]
  fn strtol(str const cstring, endp *cstring, base libc::c_int) libc::c_long

  @[link_name="SDL_strtoul"]
  fn strtoul(str const cstring, endp *cstring, base libc::c_int) libc::c_ulong

  @[link_name="SDL_strtoll"]
  fn strtoll(str const cstring, endp *cstring, base libc::c_int) libc::c_longlong

  @[link_name="SDL_strtoull"]
  fn strtoull(str const cstring, endp *cstring, base libc::c_int) libc::c_ulonglong

  @[link_name="SDL_strtod"]
  fn strtod(str const cstring, endp *cstring) libc::c_double

  @[link_name="SDL_strcmp"]
  fn strcmp(str1 const cstring, str2 const cstring) libc::c_int

  @[link_name="SDL_strncmp"]
  fn strncmp(str1 const cstring, str2 const cstring, maxlen libc::c_size_t) libc::c_int

  @[link_name="SDL_strcasecmp"]
  fn strcasecmp(str1 const cstring, str2 const cstring) libc::c_int

  @[link_name="SDL_strncasecmp"]
  fn strncasecmp(str1 const cstring, str2 const cstring, maxlen libc::c_size_t) libc::c_int

  @[link_name="SDL_strpbrk"]
  fn strpbrk(str const cstring, breakset const cstring) cstring

  @[link_name="SDL_StepUTF8"]
  fn step_utf8(pstr *const cstring, pslen *libc::c_size_t) u32

  @[link_name="SDL_StepBackUTF8"]
  fn step_back_utf8(start const cstring, pstr *const cstring) u32

  @[link_name="SDL_UCS4ToUTF8"]
  fn ucs4_to_utf8(codepoint u32, dst cstring) cstring

  @[link_name="SDL_sscanf"]
  fn sscanf(text const cstring, fmt const cstring, args @c ...) libc::c_int

  @[link_name="SDL_vsscanf"]
  fn vsscanf(text const cstring, fmt const cstring, ap libc::c_va_list) libc::c_int

  @[link_name="SDL_snprintf"]
  fn snprintf(text cstring, maxlen libc::c_size_t, fmt const cstring, args @c ...) libc::c_int

  @[link_name="SDL_swprintf"]
  fn swprintf(text *libc::c_wchar_t, maxlen libc::c_size_t, fmt const *libc::c_wchar_t, args @c ...) libc::c_int

  @[link_name="SDL_vsnprintf"]
  fn vsnprintf(text cstring, maxlen libc::c_size_t, fmt const cstring, ap libc::c_va_list) libc::c_int

  @[link_name="SDL_vswprintf"]
  fn vswprintf(text *libc::c_wchar_t, maxlen libc::c_size_t, fmt const *libc::c_wchar_t, ap libc::c_va_list) libc::c_int

  @[link_name="SDL_asprintf"]
  fn asprintf(strp *cstring, fmt const cstring, args @c ...) libc::c_int

  @[link_name="SDL_vasprintf"]
  fn vasprintf(strp *cstring, fmt const cstring, ap libc::c_va_list) libc::c_int

  @[link_name="SDL_srand"]
  fn srand(seed u64)

  @[link_name="SDL_rand"]
  fn rand(n i32) i32

  @[link_name="SDL_randf"]
  fn randf() libc::c_float

  @[link_name="SDL_rand_bits"]
  fn rand_bits() u32

  @[link_name="SDL_rand_r"]
  fn rand_r(state *u64, n i32) i32

  @[link_name="SDL_randf_r"]
  fn randf_r(state *u64) libc::c_float

  @[link_name="SDL_rand_bits_r"]
  fn rand_bits_r(state *u64) u32

  @[link_name="SDL_acos"]
  fn acos(x libc::c_double) libc::c_double

  @[link_name="SDL_acosf"]
  fn acosf(x libc::c_float) libc::c_float

  @[link_name="SDL_asin"]
  fn asin(x libc::c_double) libc::c_double

  @[link_name="SDL_asinf"]
  fn asinf(x libc::c_float) libc::c_float

  @[link_name="SDL_atan"]
  fn atan(x libc::c_double) libc::c_double

  @[link_name="SDL_atanf"]
  fn atanf(x libc::c_float) libc::c_float

  @[link_name="SDL_atan2"]
  fn atan2(y libc::c_double, x libc::c_double) libc::c_double

  @[link_name="SDL_atan2f"]
  fn atan2f(y libc::c_float, x libc::c_float) libc::c_float

  @[link_name="SDL_ceil"]
  fn ceil(x libc::c_double) libc::c_double

  @[link_name="SDL_ceilf"]
  fn ceilf(x libc::c_float) libc::c_float

  @[link_name="SDL_copysign"]
  fn copysign(x libc::c_double, y libc::c_double) libc::c_double

  @[link_name="SDL_copysignf"]
  fn copysignf(x libc::c_float, y libc::c_float) libc::c_float

  @[link_name="SDL_cos"]
  fn cos(x libc::c_double) libc::c_double

  @[link_name="SDL_cosf"]
  fn cosf(x libc::c_float) libc::c_float

  @[link_name="SDL_exp"]
  fn exp(x libc::c_double) libc::c_double

  @[link_name="SDL_expf"]
  fn expf(x libc::c_float) libc::c_float

  @[link_name="SDL_fabs"]
  fn fabs(x libc::c_double) libc::c_double

  @[link_name="SDL_fabsf"]
  fn fabsf(x libc::c_float) libc::c_float

  @[link_name="SDL_floor"]
  fn floor(x libc::c_double) libc::c_double

  @[link_name="SDL_floorf"]
  fn floorf(x libc::c_float) libc::c_float

  @[link_name="SDL_trunc"]
  fn trunc(x libc::c_double) libc::c_double

  @[link_name="SDL_truncf"]
  fn truncf(x libc::c_float) libc::c_float

  @[link_name="SDL_fmod"]
  fn fmod(x libc::c_double, y libc::c_double) libc::c_double

  @[link_name="SDL_fmodf"]
  fn fmodf(x libc::c_float, y libc::c_float) libc::c_float

  @[link_name="SDL_isinf"]
  fn isinf(x libc::c_double) libc::c_int

  @[link_name="SDL_isinff"]
  fn isinff(x libc::c_float) libc::c_int

  @[link_name="SDL_isnan"]
  fn isnan(x libc::c_double) libc::c_int

  @[link_name="SDL_isnanf"]
  fn isnanf(x libc::c_float) libc::c_int

  @[link_name="SDL_log"]
  fn log_f64(x libc::c_double) libc::c_double

  @[link_name="SDL_logf"]
  fn logf(x libc::c_float) libc::c_float

  @[link_name="SDL_log10"]
  fn log10(x libc::c_double) libc::c_double

  @[link_name="SDL_log10f"]
  fn log10f(x libc::c_float) libc::c_float

  @[link_name="SDL_modf"]
  fn modf(x libc::c_double, y *libc::c_double) libc::c_double

  @[link_name="SDL_modff"]
  fn modff(x libc::c_float, y *libc::c_float) libc::c_float

  @[link_name="SDL_pow"]
  fn pow(x libc::c_double, y libc::c_double) libc::c_double

  @[link_name="SDL_powf"]
  fn powf(x libc::c_float, y libc::c_float) libc::c_float

  @[link_name="SDL_round"]
  fn round(x libc::c_double) libc::c_double

  @[link_name="SDL_roundf"]
  fn roundf(x libc::c_float) libc::c_float

  @[link_name="SDL_lround"]
  fn lround(x libc::c_double) libc::c_long

  @[link_name="SDL_lroundf"]
  fn lroundf(x libc::c_float) libc::c_long

  @[link_name="SDL_scalbn"]
  fn scalbn(x libc::c_double, n libc::c_int) libc::c_double

  @[link_name="SDL_scalbnf"]
  fn scalbnf(x libc::c_float, n libc::c_int) libc::c_float

  @[link_name="SDL_sin"]
  fn sin(x libc::c_double) libc::c_double

  @[link_name="SDL_sinf"]
  fn sinf(x libc::c_float) libc::c_float

  @[link_name="SDL_sqrt"]
  fn sqrt(x libc::c_double) libc::c_double

  @[link_name="SDL_sqrtf"]
  fn sqrtf(x libc::c_float) libc::c_float

  @[link_name="SDL_tan"]
  fn tan(x libc::c_double) libc::c_double

  @[link_name="SDL_tanf"]
  fn tanf(x libc::c_float) libc::c_float

  @[link_name="SDL_iconv_open"]
  fn iconv_open(tocode const cstring, fromcode const cstring) iconv_t

  @[link_name="SDL_iconv_close"]
  fn iconv_close(cd iconv_t) libc::c_int

  @[link_name="SDL_iconv"]
  fn iconv(cd iconv_t, inbuf *const cstring, inbytesleft *libc::c_size_t, outbuf *cstring, outbytesleft *libc::c_size_t) libc::c_size_t

  @[link_name="SDL_iconv_string"]
  fn iconv_string(tocode const cstring, fromcode const cstring, inbuf const cstring, inbytesleft libc::c_size_t) cstring

  @[link_name="SDL_size_mul_check_overflow"]
  fn size_mul_check_overflow(a libc::c_size_t, b libc::c_size_t, ret *libc::c_size_t) bool

  @[link_name="SDL_size_mul_check_overflow_builtin"]
  fn size_mul_check_overflow_builtin(a libc::c_size_t, b libc::c_size_t, ret *libc::c_size_t) bool

  @[link_name="SDL_size_add_check_overflow"]
  fn size_add_check_overflow(a libc::c_size_t, b libc::c_size_t, ret *libc::c_size_t) bool

  @[link_name="SDL_size_add_check_overflow_builtin"]
  fn size_add_check_overflow_builtin(a libc::c_size_t, b libc::c_size_t, ret *libc::c_size_t) bool

  @[link_name="SDL_ReportAssertion"]
  fn report_assertion(data *AssertData, func const cstring, file const cstring, line libc::c_int) AssertState

  @[link_name="SDL_SetAssertionHandler"]
  fn set_assertion_handler(handler AssertionHandler, userdata rawptr)

  @[link_name="SDL_GetDefaultAssertionHandler"]
  fn get_default_assertion_handler() AssertionHandler

  @[link_name="SDL_GetAssertionHandler"]
  fn get_assertion_handler(puserdata *rawptr) AssertionHandler

  @[link_name="SDL_GetAssertionReport"]
  fn get_assertion_report() const *AssertData

  @[link_name="SDL_ResetAssertionReport"]
  fn reset_assertion_report()

  @[link_name="SDL_AsyncIOFromFile"]
  fn async_io_from_file(file const cstring, mode const cstring) *AsyncIO

  @[link_name="SDL_GetAsyncIOSize"]
  fn get_async_io_size(asyncio *AsyncIO) i64

  @[link_name="SDL_ReadAsyncIO"]
  fn read_async_io(asyncio *AsyncIO, ptr rawptr, offset u64, size u64, queue *AsyncIOQueue, userdata rawptr) bool

  @[link_name="SDL_WriteAsyncIO"]
  fn write_async_io(asyncio *AsyncIO, ptr rawptr, offset u64, size u64, queue *AsyncIOQueue, userdata rawptr) bool

  @[link_name="SDL_CloseAsyncIO"]
  fn close_async_io(asyncio *AsyncIO, flush bool, queue *AsyncIOQueue, userdata rawptr) bool

  @[link_name="SDL_CreateAsyncIOQueue"]
  fn create_async_io_queue() *AsyncIOQueue

  @[link_name="SDL_DestroyAsyncIOQueue"]
  fn destroy_async_io_queue(queue *AsyncIOQueue)

  @[link_name="SDL_GetAsyncIOResult"]
  fn get_async_io_result(queue *AsyncIOQueue, outcome *AsyncIOOutcome) bool

  @[link_name="SDL_WaitAsyncIOResult"]
  fn wait_async_io_result(queue *AsyncIOQueue, outcome *AsyncIOOutcome, timeoutMS i32) bool

  @[link_name="SDL_SignalAsyncIOQueue"]
  fn signal_async_io_queue(queue *AsyncIOQueue)

  @[link_name="SDL_LoadFileAsync"]
  fn load_file_async(file const cstring, queue *AsyncIOQueue, userdata rawptr) bool

  @[link_name="SDL_TryLockSpinlock"]
  fn try_lock_spinlock(lock *SpinLock) bool

  @[link_name="SDL_LockSpinlock"]
  fn lock_spinlock(lock *SpinLock)

  @[link_name="SDL_UnlockSpinlock"]
  fn unlock_spinlock(lock *SpinLock)

  @[link_name="SDL_MemoryBarrierReleaseFunction"]
  fn memory_barrier_release_function()

  @[link_name="SDL_MemoryBarrierAcquireFunction"]
  fn memory_barrier_acquire_function()

  @[link_name="SDL_CompareAndSwapAtomicInt"]
  fn compare_and_swap_atomic_int(a *AtomicInt, oldval libc::c_int, newval libc::c_int) bool

  @[link_name="SDL_SetAtomicInt"]
  fn set_atomic_int(a *AtomicInt, v libc::c_int) libc::c_int

  @[link_name="SDL_GetAtomicInt"]
  fn get_atomic_int(a *AtomicInt) libc::c_int

  @[link_name="SDL_AddAtomicInt"]
  fn add_atomic_int(a *AtomicInt, v libc::c_int) libc::c_int

  @[link_name="SDL_CompareAndSwapAtomicU32"]
  fn compare_and_swap_atomic_u32(a *AtomicU32, oldval u32, newval u32) bool

  @[link_name="SDL_SetAtomicU32"]
  fn set_atomic_u32(a *AtomicU32, v u32) u32

  @[link_name="SDL_GetAtomicU32"]
  fn get_atomic_u32(a *AtomicU32) u32

  @[link_name="SDL_CompareAndSwapAtomicPointer"]
  fn compare_and_swap_atomic_pointer(a *rawptr, oldval rawptr, newval rawptr) bool

  @[link_name="SDL_SetAtomicPointer"]
  fn set_atomic_pointer(a *rawptr, v rawptr) rawptr

  @[link_name="SDL_GetAtomicPointer"]
  fn get_atomic_pointer(a *rawptr) rawptr

  @[link_name="SDL_SwapFloat"]
  fn swap_float(x libc::c_float) libc::c_float

  @[link_name="SDL_SetError"]
  fn set_error(fmt const cstring, args @c ...) bool

  @[link_name="SDL_SetErrorV"]
  fn set_error_v(fmt const cstring, ap libc::c_va_list) bool

  @[link_name="SDL_OutOfMemory"]
  fn out_of_memory() bool

  @[link_name="SDL_GetError"]
  fn get_error() const cstring

  @[link_name="SDL_ClearError"]
  fn clear_error() bool

  @[link_name="SDL_GetGlobalProperties"]
  fn get_global_properties() PropertiesID

  @[link_name="SDL_CreateProperties"]
  fn create_properties() PropertiesID

  @[link_name="SDL_CopyProperties"]
  fn copy_properties(src PropertiesID, dst PropertiesID) bool

  @[link_name="SDL_LockProperties"]
  fn lock_properties(props PropertiesID) bool

  @[link_name="SDL_UnlockProperties"]
  fn unlock_properties(props PropertiesID)

  @[link_name="SDL_SetPointerPropertyWithCleanup"]
  fn set_pointer_property_with_cleanup(props PropertiesID, name const cstring, value rawptr, cleanup CleanupPropertyCallback, userdata rawptr) bool

  @[link_name="SDL_SetPointerProperty"]
  fn set_pointer_property(props PropertiesID, name const cstring, value rawptr) bool

  @[link_name="SDL_SetStringProperty"]
  fn set_string_property(props PropertiesID, name const cstring, value const cstring) bool

  @[link_name="SDL_SetNumberProperty"]
  fn set_number_property(props PropertiesID, name const cstring, value i64) bool

  @[link_name="SDL_SetFloatProperty"]
  fn set_float_property(props PropertiesID, name const cstring, value libc::c_float) bool

  @[link_name="SDL_SetBooleanProperty"]
  fn set_boolean_property(props PropertiesID, name const cstring, value bool) bool

  @[link_name="SDL_HasProperty"]
  fn has_property(props PropertiesID, name const cstring) bool

  @[link_name="SDL_GetPropertyType"]
  fn get_property_type(props PropertiesID, name const cstring) PropertyType

  @[link_name="SDL_GetPointerProperty"]
  fn get_pointer_property(props PropertiesID, name const cstring, default_value rawptr) rawptr

  @[link_name="SDL_GetStringProperty"]
  fn get_string_property(props PropertiesID, name const cstring, default_value const cstring) const cstring

  @[link_name="SDL_GetNumberProperty"]
  fn get_number_property(props PropertiesID, name const cstring, default_value i64) i64

  @[link_name="SDL_GetFloatProperty"]
  fn get_float_property(props PropertiesID, name const cstring, default_value libc::c_float) libc::c_float

  @[link_name="SDL_GetBooleanProperty"]
  fn get_boolean_property(props PropertiesID, name const cstring, default_value bool) bool

  @[link_name="SDL_ClearProperty"]
  fn clear_property(props PropertiesID, name const cstring) bool

  @[link_name="SDL_EnumerateProperties"]
  fn enumerate_properties(props PropertiesID, callback EnumeratePropertiesCallback, userdata rawptr) bool

  @[link_name="SDL_DestroyProperties"]
  fn destroy_properties(props PropertiesID)

  @[link_name="SDL_CreateThreadRuntime"]
  fn create_thread_runtime(fn_ ThreadFunction, name const cstring, data rawptr, pfnBeginThread FunctionPointer, pfnEndThread FunctionPointer) *Thread

  @[link_name="SDL_CreateThreadWithPropertiesRuntime"]
  fn create_thread_with_properties_runtime(props PropertiesID, pfnBeginThread FunctionPointer, pfnEndThread FunctionPointer) *Thread

  @[link_name="SDL_GetThreadName"]
  fn get_thread_name(thread *Thread) const cstring

  @[link_name="SDL_GetCurrentThreadID"]
  fn get_current_thread_id() ThreadID

  @[link_name="SDL_GetThreadID"]
  fn get_thread_id(thread *Thread) ThreadID

  @[link_name="SDL_SetCurrentThreadPriority"]
  fn set_current_thread_priority(priority ThreadPriority) bool

  @[link_name="SDL_WaitThread"]
  fn wait_thread(thread *Thread, status *libc::c_int)

  @[link_name="SDL_GetThreadState"]
  fn get_thread_state(thread *Thread) ThreadState

  @[link_name="SDL_DetachThread"]
  fn detach_thread(thread *Thread)

  @[link_name="SDL_GetTLS"]
  fn get_tls(id *TLSID) rawptr

  @[link_name="SDL_SetTLS"]
  fn set_tls(id *TLSID, value const *Void, destructor TLSDestructorCallback) bool

  @[link_name="SDL_CleanupTLS"]
  fn cleanup_tls()

  @[link_name="SDL_CreateMutex"]
  fn create_mutex() *Mutex

  @[link_name="SDL_LockMutex"]
  fn lock_mutex(mutex *Mutex)

  @[link_name="SDL_TryLockMutex"]
  fn try_lock_mutex(mutex *Mutex) bool

  @[link_name="SDL_UnlockMutex"]
  fn unlock_mutex(mutex *Mutex)

  @[link_name="SDL_DestroyMutex"]
  fn destroy_mutex(mutex *Mutex)

  @[link_name="SDL_CreateRWLock"]
  fn create_rw_lock() *RWLock

  @[link_name="SDL_LockRWLockForReading"]
  fn lock_rw_lock_for_reading(rwlock *RWLock)

  @[link_name="SDL_LockRWLockForWriting"]
  fn lock_rw_lock_for_writing(rwlock *RWLock)

  @[link_name="SDL_TryLockRWLockForReading"]
  fn try_lock_rw_lock_for_reading(rwlock *RWLock) bool

  @[link_name="SDL_TryLockRWLockForWriting"]
  fn try_lock_rw_lock_for_writing(rwlock *RWLock) bool

  @[link_name="SDL_UnlockRWLock"]
  fn unlock_rw_lock(rwlock *RWLock)

  @[link_name="SDL_DestroyRWLock"]
  fn destroy_rw_lock(rwlock *RWLock)

  @[link_name="SDL_CreateSemaphore"]
  fn create_semaphore(initial_value u32) *Semaphore

  @[link_name="SDL_DestroySemaphore"]
  fn destroy_semaphore(sem *Semaphore)

  @[link_name="SDL_WaitSemaphore"]
  fn wait_semaphore(sem *Semaphore)

  @[link_name="SDL_TryWaitSemaphore"]
  fn try_wait_semaphore(sem *Semaphore) bool

  @[link_name="SDL_WaitSemaphoreTimeout"]
  fn wait_semaphore_timeout(sem *Semaphore, timeoutMS i32) bool

  @[link_name="SDL_SignalSemaphore"]
  fn signal_semaphore(sem *Semaphore)

  @[link_name="SDL_GetSemaphoreValue"]
  fn get_semaphore_value(sem *Semaphore) u32

  @[link_name="SDL_CreateCondition"]
  fn create_condition() *Condition

  @[link_name="SDL_DestroyCondition"]
  fn destroy_condition(cond *Condition)

  @[link_name="SDL_SignalCondition"]
  fn signal_condition(cond *Condition)

  @[link_name="SDL_BroadcastCondition"]
  fn broadcast_condition(cond *Condition)

  @[link_name="SDL_WaitCondition"]
  fn wait_condition(cond *Condition, mutex *Mutex)

  @[link_name="SDL_WaitConditionTimeout"]
  fn wait_condition_timeout(cond *Condition, mutex *Mutex, timeoutMS i32) bool

  @[link_name="SDL_ShouldInit"]
  fn should_init(state *InitState) bool

  @[link_name="SDL_ShouldQuit"]
  fn should_quit(state *InitState) bool

  @[link_name="SDL_SetInitialized"]
  fn set_initialized(state *InitState, initialized bool)

  @[link_name="SDL_IOFromFile"]
  fn io_from_file(file const cstring, mode const cstring) *IOStream

  @[link_name="SDL_IOFromMem"]
  fn io_from_mem(mem rawptr, size libc::c_size_t) *IOStream

  @[link_name="SDL_IOFromConstMem"]
  fn io_from_const_mem(mem const *Void, size libc::c_size_t) *IOStream

  @[link_name="SDL_IOFromDynamicMem"]
  fn io_from_dynamic_mem() *IOStream

  @[link_name="SDL_OpenIO"]
  fn open_io(iface const *IOStreamInterface, userdata rawptr) *IOStream

  @[link_name="SDL_CloseIO"]
  fn close_io(context *IOStream) bool

  @[link_name="SDL_GetIOProperties"]
  fn get_io_properties(context *IOStream) PropertiesID

  @[link_name="SDL_GetIOStatus"]
  fn get_io_status(context *IOStream) IOStatus

  @[link_name="SDL_GetIOSize"]
  fn get_io_size(context *IOStream) i64

  @[link_name="SDL_SeekIO"]
  fn seek_io(context *IOStream, offset i64, whence IOWhence) i64

  @[link_name="SDL_TellIO"]
  fn tell_io(context *IOStream) i64

  @[link_name="SDL_ReadIO"]
  fn read_io(context *IOStream, ptr rawptr, size libc::c_size_t) libc::c_size_t

  @[link_name="SDL_WriteIO"]
  fn write_io(context *IOStream, ptr const *Void, size libc::c_size_t) libc::c_size_t

  @[link_name="SDL_IOprintf"]
  fn i_oprintf(context *IOStream, fmt const cstring, args @c ...) libc::c_size_t

  @[link_name="SDL_IOvprintf"]
  fn i_ovprintf(context *IOStream, fmt const cstring, ap libc::c_va_list) libc::c_size_t

  @[link_name="SDL_FlushIO"]
  fn flush_io(context *IOStream) bool

  @[link_name="SDL_LoadFile_IO"]
  fn load_file_io(src *IOStream, datasize *libc::c_size_t, closeio bool) rawptr

  @[link_name="SDL_LoadFile"]
  fn load_file(file const cstring, datasize *libc::c_size_t) rawptr

  @[link_name="SDL_SaveFile_IO"]
  fn save_file_io(src *IOStream, data const *Void, datasize libc::c_size_t, closeio bool) bool

  @[link_name="SDL_SaveFile"]
  fn save_file(file const cstring, data const *Void, datasize libc::c_size_t) bool

  @[link_name="SDL_ReadU8"]
  fn read_u8(src *IOStream, value *u8) bool

  @[link_name="SDL_ReadS8"]
  fn read_s8(src *IOStream, value *i8) bool

  @[link_name="SDL_ReadU16LE"]
  fn read_u16_le(src *IOStream, value *u16) bool

  @[link_name="SDL_ReadS16LE"]
  fn read_s16_le(src *IOStream, value *i16) bool

  @[link_name="SDL_ReadU16BE"]
  fn read_u16_be(src *IOStream, value *u16) bool

  @[link_name="SDL_ReadS16BE"]
  fn read_s16_be(src *IOStream, value *i16) bool

  @[link_name="SDL_ReadU32LE"]
  fn read_u32_le(src *IOStream, value *u32) bool

  @[link_name="SDL_ReadS32LE"]
  fn read_s32_le(src *IOStream, value *i32) bool

  @[link_name="SDL_ReadU32BE"]
  fn read_u32_be(src *IOStream, value *u32) bool

  @[link_name="SDL_ReadS32BE"]
  fn read_s32_be(src *IOStream, value *i32) bool

  @[link_name="SDL_ReadU64LE"]
  fn read_u64_le(src *IOStream, value *u64) bool

  @[link_name="SDL_ReadS64LE"]
  fn read_s64_le(src *IOStream, value *i64) bool

  @[link_name="SDL_ReadU64BE"]
  fn read_u64_be(src *IOStream, value *u64) bool

  @[link_name="SDL_ReadS64BE"]
  fn read_s64_be(src *IOStream, value *i64) bool

  @[link_name="SDL_WriteU8"]
  fn write_u8(dst *IOStream, value u8) bool

  @[link_name="SDL_WriteS8"]
  fn write_s8(dst *IOStream, value i8) bool

  @[link_name="SDL_WriteU16LE"]
  fn write_u16_le(dst *IOStream, value u16) bool

  @[link_name="SDL_WriteS16LE"]
  fn write_s16_le(dst *IOStream, value i16) bool

  @[link_name="SDL_WriteU16BE"]
  fn write_u16_be(dst *IOStream, value u16) bool

  @[link_name="SDL_WriteS16BE"]
  fn write_s16_be(dst *IOStream, value i16) bool

  @[link_name="SDL_WriteU32LE"]
  fn write_u32_le(dst *IOStream, value u32) bool

  @[link_name="SDL_WriteS32LE"]
  fn write_s32_le(dst *IOStream, value i32) bool

  @[link_name="SDL_WriteU32BE"]
  fn write_u32_be(dst *IOStream, value u32) bool

  @[link_name="SDL_WriteS32BE"]
  fn write_s32_be(dst *IOStream, value i32) bool

  @[link_name="SDL_WriteU64LE"]
  fn write_u64_le(dst *IOStream, value u64) bool

  @[link_name="SDL_WriteS64LE"]
  fn write_s64_le(dst *IOStream, value i64) bool

  @[link_name="SDL_WriteU64BE"]
  fn write_u64_be(dst *IOStream, value u64) bool

  @[link_name="SDL_WriteS64BE"]
  fn write_s64_be(dst *IOStream, value i64) bool

  @[link_name="SDL_GetNumAudioDrivers"]
  fn get_num_audio_drivers() libc::c_int

  @[link_name="SDL_GetAudioDriver"]
  fn get_audio_driver(index libc::c_int) const cstring

  @[link_name="SDL_GetCurrentAudioDriver"]
  fn get_current_audio_driver() const cstring

  @[link_name="SDL_GetAudioPlaybackDevices"]
  fn get_audio_playback_devices(count *libc::c_int) *AudioDeviceID

  @[link_name="SDL_GetAudioRecordingDevices"]
  fn get_audio_recording_devices(count *libc::c_int) *AudioDeviceID

  @[link_name="SDL_GetAudioDeviceName"]
  fn get_audio_device_name(devid AudioDeviceID) const cstring

  @[link_name="SDL_GetAudioDeviceFormat"]
  fn get_audio_device_format(devid AudioDeviceID, spec *AudioSpec, sample_frames *libc::c_int) bool

  @[link_name="SDL_GetAudioDeviceChannelMap"]
  fn get_audio_device_channel_map(devid AudioDeviceID, count *libc::c_int) *libc::c_int

  @[link_name="SDL_OpenAudioDevice"]
  fn open_audio_device(devid AudioDeviceID, spec const *AudioSpec) AudioDeviceID

  @[link_name="SDL_IsAudioDevicePhysical"]
  fn is_audio_device_physical(devid AudioDeviceID) bool

  @[link_name="SDL_IsAudioDevicePlayback"]
  fn is_audio_device_playback(devid AudioDeviceID) bool

  @[link_name="SDL_PauseAudioDevice"]
  fn pause_audio_device(devid AudioDeviceID) bool

  @[link_name="SDL_ResumeAudioDevice"]
  fn resume_audio_device(devid AudioDeviceID) bool

  @[link_name="SDL_AudioDevicePaused"]
  fn audio_device_paused(devid AudioDeviceID) bool

  @[link_name="SDL_GetAudioDeviceGain"]
  fn get_audio_device_gain(devid AudioDeviceID) libc::c_float

  @[link_name="SDL_SetAudioDeviceGain"]
  fn set_audio_device_gain(devid AudioDeviceID, gain libc::c_float) bool

  @[link_name="SDL_CloseAudioDevice"]
  fn close_audio_device(devid AudioDeviceID)

  @[link_name="SDL_BindAudioStreams"]
  fn bind_audio_streams(devid AudioDeviceID, streams const **AudioStream, num_streams libc::c_int) bool

  @[link_name="SDL_BindAudioStream"]
  fn bind_audio_stream(devid AudioDeviceID, stream *AudioStream) bool

  @[link_name="SDL_UnbindAudioStreams"]
  fn unbind_audio_streams(streams const **AudioStream, num_streams libc::c_int)

  @[link_name="SDL_UnbindAudioStream"]
  fn unbind_audio_stream(stream *AudioStream)

  @[link_name="SDL_GetAudioStreamDevice"]
  fn get_audio_stream_device(stream *AudioStream) AudioDeviceID

  @[link_name="SDL_CreateAudioStream"]
  fn create_audio_stream(src_spec const *AudioSpec, dst_spec const *AudioSpec) *AudioStream

  @[link_name="SDL_GetAudioStreamProperties"]
  fn get_audio_stream_properties(stream *AudioStream) PropertiesID

  @[link_name="SDL_GetAudioStreamFormat"]
  fn get_audio_stream_format(stream *AudioStream, src_spec *AudioSpec, dst_spec *AudioSpec) bool

  @[link_name="SDL_SetAudioStreamFormat"]
  fn set_audio_stream_format(stream *AudioStream, src_spec const *AudioSpec, dst_spec const *AudioSpec) bool

  @[link_name="SDL_GetAudioStreamFrequencyRatio"]
  fn get_audio_stream_frequency_ratio(stream *AudioStream) libc::c_float

  @[link_name="SDL_SetAudioStreamFrequencyRatio"]
  fn set_audio_stream_frequency_ratio(stream *AudioStream, ratio libc::c_float) bool

  @[link_name="SDL_GetAudioStreamGain"]
  fn get_audio_stream_gain(stream *AudioStream) libc::c_float

  @[link_name="SDL_SetAudioStreamGain"]
  fn set_audio_stream_gain(stream *AudioStream, gain libc::c_float) bool

  @[link_name="SDL_GetAudioStreamInputChannelMap"]
  fn get_audio_stream_input_channel_map(stream *AudioStream, count *libc::c_int) *libc::c_int

  @[link_name="SDL_GetAudioStreamOutputChannelMap"]
  fn get_audio_stream_output_channel_map(stream *AudioStream, count *libc::c_int) *libc::c_int

  @[link_name="SDL_SetAudioStreamInputChannelMap"]
  fn set_audio_stream_input_channel_map(stream *AudioStream, chmap const *libc::c_int, count libc::c_int) bool

  @[link_name="SDL_SetAudioStreamOutputChannelMap"]
  fn set_audio_stream_output_channel_map(stream *AudioStream, chmap const *libc::c_int, count libc::c_int) bool

  @[link_name="SDL_PutAudioStreamData"]
  fn put_audio_stream_data(stream *AudioStream, buf const *Void, len libc::c_int) bool

  @[link_name="SDL_GetAudioStreamData"]
  fn get_audio_stream_data(stream *AudioStream, buf rawptr, len libc::c_int) libc::c_int

  @[link_name="SDL_GetAudioStreamAvailable"]
  fn get_audio_stream_available(stream *AudioStream) libc::c_int

  @[link_name="SDL_GetAudioStreamQueued"]
  fn get_audio_stream_queued(stream *AudioStream) libc::c_int

  @[link_name="SDL_FlushAudioStream"]
  fn flush_audio_stream(stream *AudioStream) bool

  @[link_name="SDL_ClearAudioStream"]
  fn clear_audio_stream(stream *AudioStream) bool

  @[link_name="SDL_PauseAudioStreamDevice"]
  fn pause_audio_stream_device(stream *AudioStream) bool

  @[link_name="SDL_ResumeAudioStreamDevice"]
  fn resume_audio_stream_device(stream *AudioStream) bool

  @[link_name="SDL_AudioStreamDevicePaused"]
  fn audio_stream_device_paused(stream *AudioStream) bool

  @[link_name="SDL_LockAudioStream"]
  fn lock_audio_stream(stream *AudioStream) bool

  @[link_name="SDL_UnlockAudioStream"]
  fn unlock_audio_stream(stream *AudioStream) bool

  @[link_name="SDL_SetAudioStreamGetCallback"]
  fn set_audio_stream_get_callback(stream *AudioStream, callback AudioStreamCallback, userdata rawptr) bool

  @[link_name="SDL_SetAudioStreamPutCallback"]
  fn set_audio_stream_put_callback(stream *AudioStream, callback AudioStreamCallback, userdata rawptr) bool

  @[link_name="SDL_DestroyAudioStream"]
  fn destroy_audio_stream(stream *AudioStream)

  @[link_name="SDL_OpenAudioDeviceStream"]
  fn open_audio_device_stream(devid AudioDeviceID, spec const *AudioSpec, callback AudioStreamCallback, userdata rawptr) *AudioStream

  @[link_name="SDL_SetAudioPostmixCallback"]
  fn set_audio_postmix_callback(devid AudioDeviceID, callback AudioPostmixCallback, userdata rawptr) bool

  @[link_name="SDL_LoadWAV_IO"]
  fn load_wav_io(src *IOStream, closeio bool, spec *AudioSpec, audio_buf **u8, audio_len *u32) bool

  @[link_name="SDL_LoadWAV"]
  fn load_wav(path const cstring, spec *AudioSpec, audio_buf **u8, audio_len *u32) bool

  @[link_name="SDL_MixAudio"]
  fn mix_audio(dst *u8, src const *u8, format AudioFormat, len u32, volume libc::c_float) bool

  @[link_name="SDL_ConvertAudioSamples"]
  fn convert_audio_samples(src_spec const *AudioSpec, src_data const *u8, src_len libc::c_int, dst_spec const *AudioSpec, dst_data **u8, dst_len *libc::c_int) bool

  @[link_name="SDL_GetAudioFormatName"]
  fn get_audio_format_name(format AudioFormat) const cstring

  @[link_name="SDL_GetSilenceValueForFormat"]
  fn get_silence_value_for_format(format AudioFormat) libc::c_int

  @[link_name="SDL_MostSignificantBitIndex32"]
  fn most_significant_bit_index32(x u32) libc::c_int

  @[link_name="SDL_HasExactlyOneBitSet32"]
  fn has_exactly_one_bit_set32(x u32) bool

  @[link_name="SDL_ComposeCustomBlendMode"]
  fn compose_custom_blend_mode(srcColorFactor BlendFactor, dstColorFactor BlendFactor, colorOperation BlendOperation, srcAlphaFactor BlendFactor, dstAlphaFactor BlendFactor, alphaOperation BlendOperation) BlendMode

  @[link_name="SDL_GetPixelFormatName"]
  fn get_pixel_format_name(format PixelFormat) const cstring

  @[link_name="SDL_GetMasksForPixelFormat"]
  fn get_masks_for_pixel_format(format PixelFormat, bpp *libc::c_int, Rmask *u32, Gmask *u32, Bmask *u32, Amask *u32) bool

  @[link_name="SDL_GetPixelFormatForMasks"]
  fn get_pixel_format_for_masks(bpp libc::c_int, Rmask u32, Gmask u32, Bmask u32, Amask u32) PixelFormat

  @[link_name="SDL_GetPixelFormatDetails"]
  fn get_pixel_format_details(format PixelFormat) const *PixelFormatDetails

  @[link_name="SDL_CreatePalette"]
  fn create_palette(ncolors libc::c_int) *Palette

  @[link_name="SDL_SetPaletteColors"]
  fn set_palette_colors(palette *Palette, colors const *Color, firstcolor libc::c_int, ncolors libc::c_int) bool

  @[link_name="SDL_DestroyPalette"]
  fn destroy_palette(palette *Palette)

  @[link_name="SDL_MapRGB"]
  fn map_rgb(format const *PixelFormatDetails, palette const *Palette, r u8, g u8, b u8) u32

  @[link_name="SDL_MapRGBA"]
  fn map_rgba(format const *PixelFormatDetails, palette const *Palette, r u8, g u8, b u8, a u8) u32

  @[link_name="SDL_GetRGB"]
  fn get_rgb(pixel u32, format const *PixelFormatDetails, palette const *Palette, r *u8, g *u8, b *u8)

  @[link_name="SDL_GetRGBA"]
  fn get_rgba(pixel u32, format const *PixelFormatDetails, palette const *Palette, r *u8, g *u8, b *u8, a *u8)

  @[link_name="SDL_RectToFRect"]
  fn rect_to_f_rect(rect const *Rect, frect *FRect)

  @[link_name="SDL_PointInRect"]
  fn point_in_rect(p const *Point, r const *Rect) bool

  @[link_name="SDL_RectEmpty"]
  fn rect_empty(r const *Rect) bool

  @[link_name="SDL_RectsEqual"]
  fn rects_equal(a const *Rect, b const *Rect) bool

  @[link_name="SDL_HasRectIntersection"]
  fn has_rect_intersection(A const *Rect, B const *Rect) bool

  @[link_name="SDL_GetRectIntersection"]
  fn get_rect_intersection(A const *Rect, B const *Rect, result *Rect) bool

  @[link_name="SDL_GetRectUnion"]
  fn get_rect_union(A const *Rect, B const *Rect, result *Rect) bool

  @[link_name="SDL_GetRectEnclosingPoints"]
  fn get_rect_enclosing_points(points const *Point, count libc::c_int, clip const *Rect, result *Rect) bool

  @[link_name="SDL_GetRectAndLineIntersection"]
  fn get_rect_and_line_intersection(rect const *Rect, X1 *libc::c_int, Y1 *libc::c_int, X2 *libc::c_int, Y2 *libc::c_int) bool

  @[link_name="SDL_PointInRectFloat"]
  fn point_in_rect_float(p const *FPoint, r const *FRect) bool

  @[link_name="SDL_RectEmptyFloat"]
  fn rect_empty_float(r const *FRect) bool

  @[link_name="SDL_RectsEqualEpsilon"]
  fn rects_equal_epsilon(a const *FRect, b const *FRect, epsilon libc::c_float) bool

  @[link_name="SDL_RectsEqualFloat"]
  fn rects_equal_float(a const *FRect, b const *FRect) bool

  @[link_name="SDL_HasRectIntersectionFloat"]
  fn has_rect_intersection_float(A const *FRect, B const *FRect) bool

  @[link_name="SDL_GetRectIntersectionFloat"]
  fn get_rect_intersection_float(A const *FRect, B const *FRect, result *FRect) bool

  @[link_name="SDL_GetRectUnionFloat"]
  fn get_rect_union_float(A const *FRect, B const *FRect, result *FRect) bool

  @[link_name="SDL_GetRectEnclosingPointsFloat"]
  fn get_rect_enclosing_points_float(points const *FPoint, count libc::c_int, clip const *FRect, result *FRect) bool

  @[link_name="SDL_GetRectAndLineIntersectionFloat"]
  fn get_rect_and_line_intersection_float(rect const *FRect, X1 *libc::c_float, Y1 *libc::c_float, X2 *libc::c_float, Y2 *libc::c_float) bool

  @[link_name="SDL_CreateSurface"]
  fn create_surface(width libc::c_int, height libc::c_int, format PixelFormat) *Surface

  @[link_name="SDL_CreateSurfaceFrom"]
  fn create_surface_from(width libc::c_int, height libc::c_int, format PixelFormat, pixels rawptr, pitch libc::c_int) *Surface

  @[link_name="SDL_DestroySurface"]
  fn destroy_surface(surface *Surface)

  @[link_name="SDL_GetSurfaceProperties"]
  fn get_surface_properties(surface *Surface) PropertiesID

  @[link_name="SDL_SetSurfaceColorspace"]
  fn set_surface_colorspace(surface *Surface, colorspace Colorspace) bool

  @[link_name="SDL_GetSurfaceColorspace"]
  fn get_surface_colorspace(surface *Surface) Colorspace

  @[link_name="SDL_CreateSurfacePalette"]
  fn create_surface_palette(surface *Surface) *Palette

  @[link_name="SDL_SetSurfacePalette"]
  fn set_surface_palette(surface *Surface, palette *Palette) bool

  @[link_name="SDL_GetSurfacePalette"]
  fn get_surface_palette(surface *Surface) *Palette

  @[link_name="SDL_AddSurfaceAlternateImage"]
  fn add_surface_alternate_image(surface *Surface, image *Surface) bool

  @[link_name="SDL_SurfaceHasAlternateImages"]
  fn surface_has_alternate_images(surface *Surface) bool

  @[link_name="SDL_GetSurfaceImages"]
  fn get_surface_images(surface *Surface, count *libc::c_int) **Surface

  @[link_name="SDL_RemoveSurfaceAlternateImages"]
  fn remove_surface_alternate_images(surface *Surface)

  @[link_name="SDL_LockSurface"]
  fn lock_surface(surface *Surface) bool

  @[link_name="SDL_UnlockSurface"]
  fn unlock_surface(surface *Surface)

  @[link_name="SDL_LoadBMP_IO"]
  fn load_bmp_io(src *IOStream, closeio bool) *Surface

  @[link_name="SDL_LoadBMP"]
  fn load_bmp(file const cstring) *Surface

  @[link_name="SDL_SaveBMP_IO"]
  fn save_bmp_io(surface *Surface, dst *IOStream, closeio bool) bool

  @[link_name="SDL_SaveBMP"]
  fn save_bmp(surface *Surface, file const cstring) bool

  @[link_name="SDL_SetSurfaceRLE"]
  fn set_surface_rle(surface *Surface, enabled bool) bool

  @[link_name="SDL_SurfaceHasRLE"]
  fn surface_has_rle(surface *Surface) bool

  @[link_name="SDL_SetSurfaceColorKey"]
  fn set_surface_color_key(surface *Surface, enabled bool, key u32) bool

  @[link_name="SDL_SurfaceHasColorKey"]
  fn surface_has_color_key(surface *Surface) bool

  @[link_name="SDL_GetSurfaceColorKey"]
  fn get_surface_color_key(surface *Surface, key *u32) bool

  @[link_name="SDL_SetSurfaceColorMod"]
  fn set_surface_color_mod(surface *Surface, r u8, g u8, b u8) bool

  @[link_name="SDL_GetSurfaceColorMod"]
  fn get_surface_color_mod(surface *Surface, r *u8, g *u8, b *u8) bool

  @[link_name="SDL_SetSurfaceAlphaMod"]
  fn set_surface_alpha_mod(surface *Surface, alpha u8) bool

  @[link_name="SDL_GetSurfaceAlphaMod"]
  fn get_surface_alpha_mod(surface *Surface, alpha *u8) bool

  @[link_name="SDL_SetSurfaceBlendMode"]
  fn set_surface_blend_mode(surface *Surface, blendMode BlendMode) bool

  @[link_name="SDL_GetSurfaceBlendMode"]
  fn get_surface_blend_mode(surface *Surface, blendMode *BlendMode) bool

  @[link_name="SDL_SetSurfaceClipRect"]
  fn set_surface_clip_rect(surface *Surface, rect const *Rect) bool

  @[link_name="SDL_GetSurfaceClipRect"]
  fn get_surface_clip_rect(surface *Surface, rect *Rect) bool

  @[link_name="SDL_FlipSurface"]
  fn flip_surface(surface *Surface, flip FlipMode) bool

  @[link_name="SDL_DuplicateSurface"]
  fn duplicate_surface(surface *Surface) *Surface

  @[link_name="SDL_ScaleSurface"]
  fn scale_surface(surface *Surface, width libc::c_int, height libc::c_int, scaleMode ScaleMode) *Surface

  @[link_name="SDL_ConvertSurface"]
  fn convert_surface(surface *Surface, format PixelFormat) *Surface

  @[link_name="SDL_ConvertSurfaceAndColorspace"]
  fn convert_surface_and_colorspace(surface *Surface, format PixelFormat, palette *Palette, colorspace Colorspace, props PropertiesID) *Surface

  @[link_name="SDL_ConvertPixels"]
  fn convert_pixels(width libc::c_int, height libc::c_int, src_format PixelFormat, src const *Void, src_pitch libc::c_int, dst_format PixelFormat, dst rawptr, dst_pitch libc::c_int) bool

  @[link_name="SDL_ConvertPixelsAndColorspace"]
  fn convert_pixels_and_colorspace(width libc::c_int, height libc::c_int, src_format PixelFormat, src_colorspace Colorspace, src_properties PropertiesID, src const *Void, src_pitch libc::c_int, dst_format PixelFormat, dst_colorspace Colorspace, dst_properties PropertiesID, dst rawptr, dst_pitch libc::c_int) bool

  @[link_name="SDL_PremultiplyAlpha"]
  fn premultiply_alpha(width libc::c_int, height libc::c_int, src_format PixelFormat, src const *Void, src_pitch libc::c_int, dst_format PixelFormat, dst rawptr, dst_pitch libc::c_int, linear bool) bool

  @[link_name="SDL_PremultiplySurfaceAlpha"]
  fn premultiply_surface_alpha(surface *Surface, linear bool) bool

  @[link_name="SDL_ClearSurface"]
  fn clear_surface(surface *Surface, r libc::c_float, g libc::c_float, b libc::c_float, a libc::c_float) bool

  @[link_name="SDL_FillSurfaceRect"]
  fn fill_surface_rect(dst *Surface, rect const *Rect, color u32) bool

  @[link_name="SDL_FillSurfaceRects"]
  fn fill_surface_rects(dst *Surface, rects const *Rect, count libc::c_int, color u32) bool

  @[link_name="SDL_BlitSurface"]
  fn blit_surface(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect) bool

  @[link_name="SDL_BlitSurfaceUnchecked"]
  fn blit_surface_unchecked(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect) bool

  @[link_name="SDL_BlitSurfaceScaled"]
  fn blit_surface_scaled(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect, scaleMode ScaleMode) bool

  @[link_name="SDL_BlitSurfaceUncheckedScaled"]
  fn blit_surface_unchecked_scaled(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect, scaleMode ScaleMode) bool

  @[link_name="SDL_StretchSurface"]
  fn stretch_surface(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect, scaleMode ScaleMode) bool

  @[link_name="SDL_BlitSurfaceTiled"]
  fn blit_surface_tiled(src *Surface, srcrect const *Rect, dst *Surface, dstrect const *Rect) bool

  @[link_name="SDL_BlitSurfaceTiledWithScale"]
  fn blit_surface_tiled_with_scale(src *Surface, srcrect const *Rect, scale libc::c_float, scaleMode ScaleMode, dst *Surface, dstrect const *Rect) bool

  @[link_name="SDL_BlitSurface9Grid"]
  fn blit_surface9_grid(src *Surface, srcrect const *Rect, left_width libc::c_int, right_width libc::c_int, top_height libc::c_int, bottom_height libc::c_int, scale libc::c_float, scaleMode ScaleMode, dst *Surface, dstrect const *Rect) bool

  @[link_name="SDL_MapSurfaceRGB"]
  fn map_surface_rgb(surface *Surface, r u8, g u8, b u8) u32

  @[link_name="SDL_MapSurfaceRGBA"]
  fn map_surface_rgba(surface *Surface, r u8, g u8, b u8, a u8) u32

  @[link_name="SDL_ReadSurfacePixel"]
  fn read_surface_pixel(surface *Surface, x libc::c_int, y libc::c_int, r *u8, g *u8, b *u8, a *u8) bool

  @[link_name="SDL_ReadSurfacePixelFloat"]
  fn read_surface_pixel_float(surface *Surface, x libc::c_int, y libc::c_int, r *libc::c_float, g *libc::c_float, b *libc::c_float, a *libc::c_float) bool

  @[link_name="SDL_WriteSurfacePixel"]
  fn write_surface_pixel(surface *Surface, x libc::c_int, y libc::c_int, r u8, g u8, b u8, a u8) bool

  @[link_name="SDL_WriteSurfacePixelFloat"]
  fn write_surface_pixel_float(surface *Surface, x libc::c_int, y libc::c_int, r libc::c_float, g libc::c_float, b libc::c_float, a libc::c_float) bool

  @[link_name="SDL_GetNumCameraDrivers"]
  fn get_num_camera_drivers() libc::c_int

  @[link_name="SDL_GetCameraDriver"]
  fn get_camera_driver(index libc::c_int) const cstring

  @[link_name="SDL_GetCurrentCameraDriver"]
  fn get_current_camera_driver() const cstring

  @[link_name="SDL_GetCameras"]
  fn get_cameras(count *libc::c_int) *CameraID

  @[link_name="SDL_GetCameraSupportedFormats"]
  fn get_camera_supported_formats(instance_id CameraID, count *libc::c_int) **CameraSpec

  @[link_name="SDL_GetCameraName"]
  fn get_camera_name(instance_id CameraID) const cstring

  @[link_name="SDL_GetCameraPosition"]
  fn get_camera_position(instance_id CameraID) CameraPosition

  @[link_name="SDL_OpenCamera"]
  fn open_camera(instance_id CameraID, spec const *CameraSpec) *Camera

  @[link_name="SDL_GetCameraPermissionState"]
  fn get_camera_permission_state(camera *Camera) libc::c_int

  @[link_name="SDL_GetCameraID"]
  fn get_camera_id(camera *Camera) CameraID

  @[link_name="SDL_GetCameraProperties"]
  fn get_camera_properties(camera *Camera) PropertiesID

  @[link_name="SDL_GetCameraFormat"]
  fn get_camera_format(camera *Camera, spec *CameraSpec) bool

  @[link_name="SDL_AcquireCameraFrame"]
  fn acquire_camera_frame(camera *Camera, timestampNS *u64) *Surface

  @[link_name="SDL_ReleaseCameraFrame"]
  fn release_camera_frame(camera *Camera, frame *Surface)

  @[link_name="SDL_CloseCamera"]
  fn close_camera(camera *Camera)

  @[link_name="SDL_SetClipboardText"]
  fn set_clipboard_text(text const cstring) bool

  @[link_name="SDL_GetClipboardText"]
  fn get_clipboard_text() cstring

  @[link_name="SDL_HasClipboardText"]
  fn has_clipboard_text() bool

  @[link_name="SDL_SetPrimarySelectionText"]
  fn set_primary_selection_text(text const cstring) bool

  @[link_name="SDL_GetPrimarySelectionText"]
  fn get_primary_selection_text() cstring

  @[link_name="SDL_HasPrimarySelectionText"]
  fn has_primary_selection_text() bool

  @[link_name="SDL_SetClipboardData"]
  fn set_clipboard_data(callback ClipboardDataCallback, cleanup ClipboardCleanupCallback, userdata rawptr, mime_types *const cstring, num_mime_types libc::c_size_t) bool

  @[link_name="SDL_ClearClipboardData"]
  fn clear_clipboard_data() bool

  @[link_name="SDL_GetClipboardData"]
  fn get_clipboard_data(mime_type const cstring, size *libc::c_size_t) rawptr

  @[link_name="SDL_HasClipboardData"]
  fn has_clipboard_data(mime_type const cstring) bool

  @[link_name="SDL_GetClipboardMimeTypes"]
  fn get_clipboard_mime_types(num_mime_types *libc::c_size_t) *cstring

  @[link_name="SDL_GetNumLogicalCPUCores"]
  fn get_num_logical_cpu_cores() libc::c_int

  @[link_name="SDL_GetCPUCacheLineSize"]
  fn get_cpu_cache_line_size() libc::c_int

  @[link_name="SDL_HasAltiVec"]
  fn has_alti_vec() bool

  @[link_name="SDL_HasMMX"]
  fn has_mmx() bool

  @[link_name="SDL_HasSSE"]
  fn has_sse() bool

  @[link_name="SDL_HasSSE2"]
  fn has_sse2() bool

  @[link_name="SDL_HasSSE3"]
  fn has_sse3() bool

  @[link_name="SDL_HasSSE41"]
  fn has_sse41() bool

  @[link_name="SDL_HasSSE42"]
  fn has_sse42() bool

  @[link_name="SDL_HasAVX"]
  fn has_avx() bool

  @[link_name="SDL_HasAVX2"]
  fn has_avx2() bool

  @[link_name="SDL_HasAVX512F"]
  fn has_avx512_f() bool

  @[link_name="SDL_HasARMSIMD"]
  fn has_armsimd() bool

  @[link_name="SDL_HasNEON"]
  fn has_neon() bool

  @[link_name="SDL_HasLSX"]
  fn has_lsx() bool

  @[link_name="SDL_HasLASX"]
  fn has_lasx() bool

  @[link_name="SDL_GetSystemRAM"]
  fn get_system_ram() libc::c_int

  @[link_name="SDL_GetSIMDAlignment"]
  fn get_simd_alignment() libc::c_size_t

  @[link_name="SDL_GetNumVideoDrivers"]
  fn get_num_video_drivers() libc::c_int

  @[link_name="SDL_GetVideoDriver"]
  fn get_video_driver(index libc::c_int) const cstring

  @[link_name="SDL_GetCurrentVideoDriver"]
  fn get_current_video_driver() const cstring

  @[link_name="SDL_GetSystemTheme"]
  fn get_system_theme() SystemTheme

  @[link_name="SDL_GetDisplays"]
  fn get_displays(count *libc::c_int) *DisplayID

  @[link_name="SDL_GetPrimaryDisplay"]
  fn get_primary_display() DisplayID

  @[link_name="SDL_GetDisplayProperties"]
  fn get_display_properties(displayID DisplayID) PropertiesID

  @[link_name="SDL_GetDisplayName"]
  fn get_display_name(displayID DisplayID) const cstring

  @[link_name="SDL_GetDisplayBounds"]
  fn get_display_bounds(displayID DisplayID, rect *Rect) bool

  @[link_name="SDL_GetDisplayUsableBounds"]
  fn get_display_usable_bounds(displayID DisplayID, rect *Rect) bool

  @[link_name="SDL_GetNaturalDisplayOrientation"]
  fn get_natural_display_orientation(displayID DisplayID) DisplayOrientation

  @[link_name="SDL_GetCurrentDisplayOrientation"]
  fn get_current_display_orientation(displayID DisplayID) DisplayOrientation

  @[link_name="SDL_GetDisplayContentScale"]
  fn get_display_content_scale(displayID DisplayID) libc::c_float

  @[link_name="SDL_GetFullscreenDisplayModes"]
  fn get_fullscreen_display_modes(displayID DisplayID, count *libc::c_int) **DisplayMode

  @[link_name="SDL_GetClosestFullscreenDisplayMode"]
  fn get_closest_fullscreen_display_mode(displayID DisplayID, w libc::c_int, h libc::c_int, refresh_rate libc::c_float, include_high_density_modes bool, closest *DisplayMode) bool

  @[link_name="SDL_GetDesktopDisplayMode"]
  fn get_desktop_display_mode(displayID DisplayID) const *DisplayMode

  @[link_name="SDL_GetCurrentDisplayMode"]
  fn get_current_display_mode(displayID DisplayID) const *DisplayMode

  @[link_name="SDL_GetDisplayForPoint"]
  fn get_display_for_point(point const *Point) DisplayID

  @[link_name="SDL_GetDisplayForRect"]
  fn get_display_for_rect(rect const *Rect) DisplayID

  @[link_name="SDL_GetDisplayForWindow"]
  fn get_display_for_window(window *Window) DisplayID

  @[link_name="SDL_GetWindowPixelDensity"]
  fn get_window_pixel_density(window *Window) libc::c_float

  @[link_name="SDL_GetWindowDisplayScale"]
  fn get_window_display_scale(window *Window) libc::c_float

  @[link_name="SDL_SetWindowFullscreenMode"]
  fn set_window_fullscreen_mode(window *Window, mode const *DisplayMode) bool

  @[link_name="SDL_GetWindowFullscreenMode"]
  fn get_window_fullscreen_mode(window *Window) const *DisplayMode

  @[link_name="SDL_GetWindowICCProfile"]
  fn get_window_icc_profile(window *Window, size *libc::c_size_t) rawptr

  @[link_name="SDL_GetWindowPixelFormat"]
  fn get_window_pixel_format(window *Window) PixelFormat

  @[link_name="SDL_GetWindows"]
  fn get_windows(count *libc::c_int) **Window

  @[link_name="SDL_CreateWindow"]
  fn create_window(title const cstring, w libc::c_int, h libc::c_int, flags WindowFlags) *Window

  @[link_name="SDL_CreatePopupWindow"]
  fn create_popup_window(parent *Window, offset_x libc::c_int, offset_y libc::c_int, w libc::c_int, h libc::c_int, flags WindowFlags) *Window

  @[link_name="SDL_CreateWindowWithProperties"]
  fn create_window_with_properties(props PropertiesID) *Window

  @[link_name="SDL_GetWindowID"]
  fn get_window_id(window *Window) WindowID

  @[link_name="SDL_GetWindowFromID"]
  fn get_window_from_id(id WindowID) *Window

  @[link_name="SDL_GetWindowParent"]
  fn get_window_parent(window *Window) *Window

  @[link_name="SDL_GetWindowProperties"]
  fn get_window_properties(window *Window) PropertiesID

  @[link_name="SDL_GetWindowFlags"]
  fn get_window_flags(window *Window) WindowFlags

  @[link_name="SDL_SetWindowTitle"]
  fn set_window_title(window *Window, title const cstring) bool

  @[link_name="SDL_GetWindowTitle"]
  fn get_window_title(window *Window) const cstring

  @[link_name="SDL_SetWindowIcon"]
  fn set_window_icon(window *Window, icon *Surface) bool

  @[link_name="SDL_SetWindowPosition"]
  fn set_window_position(window *Window, x libc::c_int, y libc::c_int) bool

  @[link_name="SDL_GetWindowPosition"]
  fn get_window_position(window *Window, x *libc::c_int, y *libc::c_int) bool

  @[link_name="SDL_SetWindowSize"]
  fn set_window_size(window *Window, w libc::c_int, h libc::c_int) bool

  @[link_name="SDL_GetWindowSize"]
  fn get_window_size(window *Window, w *libc::c_int, h *libc::c_int) bool

  @[link_name="SDL_GetWindowSafeArea"]
  fn get_window_safe_area(window *Window, rect *Rect) bool

  @[link_name="SDL_SetWindowAspectRatio"]
  fn set_window_aspect_ratio(window *Window, min_aspect libc::c_float, max_aspect libc::c_float) bool

  @[link_name="SDL_GetWindowAspectRatio"]
  fn get_window_aspect_ratio(window *Window, min_aspect *libc::c_float, max_aspect *libc::c_float) bool

  @[link_name="SDL_GetWindowBordersSize"]
  fn get_window_borders_size(window *Window, top *libc::c_int, left *libc::c_int, bottom *libc::c_int, right *libc::c_int) bool

  @[link_name="SDL_GetWindowSizeInPixels"]
  fn get_window_size_in_pixels(window *Window, w *libc::c_int, h *libc::c_int) bool

  @[link_name="SDL_SetWindowMinimumSize"]
  fn set_window_minimum_size(window *Window, min_w libc::c_int, min_h libc::c_int) bool

  @[link_name="SDL_GetWindowMinimumSize"]
  fn get_window_minimum_size(window *Window, w *libc::c_int, h *libc::c_int) bool

  @[link_name="SDL_SetWindowMaximumSize"]
  fn set_window_maximum_size(window *Window, max_w libc::c_int, max_h libc::c_int) bool

  @[link_name="SDL_GetWindowMaximumSize"]
  fn get_window_maximum_size(window *Window, w *libc::c_int, h *libc::c_int) bool

  @[link_name="SDL_SetWindowBordered"]
  fn set_window_bordered(window *Window, bordered bool) bool

  @[link_name="SDL_SetWindowResizable"]
  fn set_window_resizable(window *Window, resizable bool) bool

  @[link_name="SDL_SetWindowAlwaysOnTop"]
  fn set_window_always_on_top(window *Window, on_top bool) bool

  @[link_name="SDL_ShowWindow"]
  fn show_window(window *Window) bool

  @[link_name="SDL_HideWindow"]
  fn hide_window(window *Window) bool

  @[link_name="SDL_RaiseWindow"]
  fn raise_window(window *Window) bool

  @[link_name="SDL_MaximizeWindow"]
  fn maximize_window(window *Window) bool

  @[link_name="SDL_MinimizeWindow"]
  fn minimize_window(window *Window) bool

  @[link_name="SDL_RestoreWindow"]
  fn restore_window(window *Window) bool

  @[link_name="SDL_SetWindowFullscreen"]
  fn set_window_fullscreen(window *Window, fullscreen bool) bool

  @[link_name="SDL_SyncWindow"]
  fn sync_window(window *Window) bool

  @[link_name="SDL_WindowHasSurface"]
  fn window_has_surface(window *Window) bool

  @[link_name="SDL_GetWindowSurface"]
  fn get_window_surface(window *Window) *Surface

  @[link_name="SDL_SetWindowSurfaceVSync"]
  fn set_window_surface_v_sync(window *Window, vsync libc::c_int) bool

  @[link_name="SDL_GetWindowSurfaceVSync"]
  fn get_window_surface_v_sync(window *Window, vsync *libc::c_int) bool

  @[link_name="SDL_UpdateWindowSurface"]
  fn update_window_surface(window *Window) bool

  @[link_name="SDL_UpdateWindowSurfaceRects"]
  fn update_window_surface_rects(window *Window, rects const *Rect, numrects libc::c_int) bool

  @[link_name="SDL_DestroyWindowSurface"]
  fn destroy_window_surface(window *Window) bool

  @[link_name="SDL_SetWindowKeyboardGrab"]
  fn set_window_keyboard_grab(window *Window, grabbed bool) bool

  @[link_name="SDL_SetWindowMouseGrab"]
  fn set_window_mouse_grab(window *Window, grabbed bool) bool

  @[link_name="SDL_GetWindowKeyboardGrab"]
  fn get_window_keyboard_grab(window *Window) bool

  @[link_name="SDL_GetWindowMouseGrab"]
  fn get_window_mouse_grab(window *Window) bool

  @[link_name="SDL_GetGrabbedWindow"]
  fn get_grabbed_window() *Window

  @[link_name="SDL_SetWindowMouseRect"]
  fn set_window_mouse_rect(window *Window, rect const *Rect) bool

  @[link_name="SDL_GetWindowMouseRect"]
  fn get_window_mouse_rect(window *Window) const *Rect

  @[link_name="SDL_SetWindowOpacity"]
  fn set_window_opacity(window *Window, opacity libc::c_float) bool

  @[link_name="SDL_GetWindowOpacity"]
  fn get_window_opacity(window *Window) libc::c_float

  @[link_name="SDL_SetWindowParent"]
  fn set_window_parent(window *Window, parent *Window) bool

  @[link_name="SDL_SetWindowModal"]
  fn set_window_modal(window *Window, modal bool) bool

  @[link_name="SDL_SetWindowFocusable"]
  fn set_window_focusable(window *Window, focusable bool) bool

  @[link_name="SDL_ShowWindowSystemMenu"]
  fn show_window_system_menu(window *Window, x libc::c_int, y libc::c_int) bool

  @[link_name="SDL_SetWindowHitTest"]
  fn set_window_hit_test(window *Window, callback HitTest, callback_data rawptr) bool

  @[link_name="SDL_SetWindowShape"]
  fn set_window_shape(window *Window, shape *Surface) bool

  @[link_name="SDL_FlashWindow"]
  fn flash_window(window *Window, operation FlashOperation) bool

  @[link_name="SDL_DestroyWindow"]
  fn destroy_window(window *Window)

  @[link_name="SDL_ScreenSaverEnabled"]
  fn screen_saver_enabled() bool

  @[link_name="SDL_EnableScreenSaver"]
  fn enable_screen_saver() bool

  @[link_name="SDL_DisableScreenSaver"]
  fn disable_screen_saver() bool

  @[link_name="SDL_GL_LoadLibrary"]
  fn gl_load_library(path const cstring) bool

  @[link_name="SDL_GL_GetProcAddress"]
  fn gl_get_proc_address(proc const cstring) FunctionPointer

  @[link_name="SDL_EGL_GetProcAddress"]
  fn egl_get_proc_address(proc const cstring) FunctionPointer

  @[link_name="SDL_GL_UnloadLibrary"]
  fn gl_unload_library()

  @[link_name="SDL_GL_ExtensionSupported"]
  fn gl_extension_supported(extension const cstring) bool

  @[link_name="SDL_GL_ResetAttributes"]
  fn gl_reset_attributes()

  @[link_name="SDL_GL_SetAttribute"]
  fn gl_set_attribute(attr GLAttr, value libc::c_int) bool

  @[link_name="SDL_GL_GetAttribute"]
  fn gl_get_attribute(attr GLAttr, value *libc::c_int) bool

  @[link_name="SDL_GL_CreateContext"]
  fn gl_create_context(window *Window) GLContext

  @[link_name="SDL_GL_MakeCurrent"]
  fn gl_make_current(window *Window, context GLContext) bool

  @[link_name="SDL_GL_GetCurrentWindow"]
  fn gl_get_current_window() *Window

  @[link_name="SDL_GL_GetCurrentContext"]
  fn gl_get_current_context() GLContext

  @[link_name="SDL_EGL_GetCurrentDisplay"]
  fn egl_get_current_display() EGLDisplay

  @[link_name="SDL_EGL_GetCurrentConfig"]
  fn egl_get_current_config() EGLConfig

  @[link_name="SDL_EGL_GetWindowSurface"]
  fn egl_get_window_surface(window *Window) EGLSurface

  @[link_name="SDL_EGL_SetAttributeCallbacks"]
  fn egl_set_attribute_callbacks(platformAttribCallback EGLAttribArrayCallback, surfaceAttribCallback EGLIntArrayCallback, contextAttribCallback EGLIntArrayCallback, userdata rawptr)

  @[link_name="SDL_GL_SetSwapInterval"]
  fn gl_set_swap_interval(interval libc::c_int) bool

  @[link_name="SDL_GL_GetSwapInterval"]
  fn gl_get_swap_interval(interval *libc::c_int) bool

  @[link_name="SDL_GL_SwapWindow"]
  fn gl_swap_window(window *Window) bool

  @[link_name="SDL_GL_DestroyContext"]
  fn gl_destroy_context(context GLContext) bool

  @[link_name="SDL_ShowOpenFileDialog"]
  fn show_open_file_dialog(callback DialogFileCallback, userdata rawptr, window *Window, filters const *DialogFileFilter, nfilters libc::c_int, default_location const cstring, allow_many bool)

  @[link_name="SDL_ShowSaveFileDialog"]
  fn show_save_file_dialog(callback DialogFileCallback, userdata rawptr, window *Window, filters const *DialogFileFilter, nfilters libc::c_int, default_location const cstring)

  @[link_name="SDL_ShowOpenFolderDialog"]
  fn show_open_folder_dialog(callback DialogFileCallback, userdata rawptr, window *Window, default_location const cstring, allow_many bool)

  @[link_name="SDL_ShowFileDialogWithProperties"]
  fn show_file_dialog_with_properties(type_ FileDialogType, callback DialogFileCallback, userdata rawptr, props PropertiesID)

  @[link_name="SDL_GUIDToString"]
  fn guid_to_string(guid GUID, pszGUID cstring, cbGUID libc::c_int)

  @[link_name="SDL_StringToGUID"]
  fn string_to_guid(pchGUID const cstring) GUID

  @[link_name="SDL_GetPowerInfo"]
  fn get_power_info(seconds *libc::c_int, percent *libc::c_int) PowerState

  @[link_name="SDL_GetSensors"]
  fn get_sensors(count *libc::c_int) *SensorID

  @[link_name="SDL_GetSensorNameForID"]
  fn get_sensor_name_for_id(instance_id SensorID) const cstring

  @[link_name="SDL_GetSensorTypeForID"]
  fn get_sensor_type_for_id(instance_id SensorID) SensorType

  @[link_name="SDL_GetSensorNonPortableTypeForID"]
  fn get_sensor_non_portable_type_for_id(instance_id SensorID) libc::c_int

  @[link_name="SDL_OpenSensor"]
  fn open_sensor(instance_id SensorID) *Sensor

  @[link_name="SDL_GetSensorFromID"]
  fn get_sensor_from_id(instance_id SensorID) *Sensor

  @[link_name="SDL_GetSensorProperties"]
  fn get_sensor_properties(sensor *Sensor) PropertiesID

  @[link_name="SDL_GetSensorName"]
  fn get_sensor_name(sensor *Sensor) const cstring

  @[link_name="SDL_GetSensorType"]
  fn get_sensor_type(sensor *Sensor) SensorType

  @[link_name="SDL_GetSensorNonPortableType"]
  fn get_sensor_non_portable_type(sensor *Sensor) libc::c_int

  @[link_name="SDL_GetSensorID"]
  fn get_sensor_id(sensor *Sensor) SensorID

  @[link_name="SDL_GetSensorData"]
  fn get_sensor_data(sensor *Sensor, data *libc::c_float, num_values libc::c_int) bool

  @[link_name="SDL_CloseSensor"]
  fn close_sensor(sensor *Sensor)

  @[link_name="SDL_UpdateSensors"]
  fn update_sensors()

  @[link_name="SDL_LockJoysticks"]
  fn lock_joysticks()

  @[link_name="SDL_UnlockJoysticks"]
  fn unlock_joysticks()

  @[link_name="SDL_HasJoystick"]
  fn has_joystick() bool

  @[link_name="SDL_GetJoysticks"]
  fn get_joysticks(count *libc::c_int) *JoystickID

  @[link_name="SDL_GetJoystickNameForID"]
  fn get_joystick_name_for_id(instance_id JoystickID) const cstring

  @[link_name="SDL_GetJoystickPathForID"]
  fn get_joystick_path_for_id(instance_id JoystickID) const cstring

  @[link_name="SDL_GetJoystickPlayerIndexForID"]
  fn get_joystick_player_index_for_id(instance_id JoystickID) libc::c_int

  @[link_name="SDL_GetJoystickGUIDForID"]
  fn get_joystick_guid_for_id(instance_id JoystickID) GUID

  @[link_name="SDL_GetJoystickVendorForID"]
  fn get_joystick_vendor_for_id(instance_id JoystickID) u16

  @[link_name="SDL_GetJoystickProductForID"]
  fn get_joystick_product_for_id(instance_id JoystickID) u16

  @[link_name="SDL_GetJoystickProductVersionForID"]
  fn get_joystick_product_version_for_id(instance_id JoystickID) u16

  @[link_name="SDL_GetJoystickTypeForID"]
  fn get_joystick_type_for_id(instance_id JoystickID) JoystickType

  @[link_name="SDL_OpenJoystick"]
  fn open_joystick(instance_id JoystickID) *Joystick

  @[link_name="SDL_GetJoystickFromID"]
  fn get_joystick_from_id(instance_id JoystickID) *Joystick

  @[link_name="SDL_GetJoystickFromPlayerIndex"]
  fn get_joystick_from_player_index(player_index libc::c_int) *Joystick

  @[link_name="SDL_AttachVirtualJoystick"]
  fn attach_virtual_joystick(desc const *VirtualJoystickDesc) JoystickID

  @[link_name="SDL_DetachVirtualJoystick"]
  fn detach_virtual_joystick(instance_id JoystickID) bool

  @[link_name="SDL_IsJoystickVirtual"]
  fn is_joystick_virtual(instance_id JoystickID) bool

  @[link_name="SDL_SetJoystickVirtualAxis"]
  fn set_joystick_virtual_axis(joystick *Joystick, axis libc::c_int, value i16) bool

  @[link_name="SDL_SetJoystickVirtualBall"]
  fn set_joystick_virtual_ball(joystick *Joystick, ball libc::c_int, xrel i16, yrel i16) bool

  @[link_name="SDL_SetJoystickVirtualButton"]
  fn set_joystick_virtual_button(joystick *Joystick, button libc::c_int, down bool) bool

  @[link_name="SDL_SetJoystickVirtualHat"]
  fn set_joystick_virtual_hat(joystick *Joystick, hat libc::c_int, value u8) bool

  @[link_name="SDL_SetJoystickVirtualTouchpad"]
  fn set_joystick_virtual_touchpad(joystick *Joystick, touchpad libc::c_int, finger libc::c_int, down bool, x libc::c_float, y libc::c_float, pressure libc::c_float) bool

  @[link_name="SDL_SendJoystickVirtualSensorData"]
  fn send_joystick_virtual_sensor_data(joystick *Joystick, type_ SensorType, sensor_timestamp u64, data const *libc::c_float, num_values libc::c_int) bool

  @[link_name="SDL_GetJoystickProperties"]
  fn get_joystick_properties(joystick *Joystick) PropertiesID

  @[link_name="SDL_GetJoystickName"]
  fn get_joystick_name(joystick *Joystick) const cstring

  @[link_name="SDL_GetJoystickPath"]
  fn get_joystick_path(joystick *Joystick) const cstring

  @[link_name="SDL_GetJoystickPlayerIndex"]
  fn get_joystick_player_index(joystick *Joystick) libc::c_int

  @[link_name="SDL_SetJoystickPlayerIndex"]
  fn set_joystick_player_index(joystick *Joystick, player_index libc::c_int) bool

  @[link_name="SDL_GetJoystickGUID"]
  fn get_joystick_guid(joystick *Joystick) GUID

  @[link_name="SDL_GetJoystickVendor"]
  fn get_joystick_vendor(joystick *Joystick) u16

  @[link_name="SDL_GetJoystickProduct"]
  fn get_joystick_product(joystick *Joystick) u16

  @[link_name="SDL_GetJoystickProductVersion"]
  fn get_joystick_product_version(joystick *Joystick) u16

  @[link_name="SDL_GetJoystickFirmwareVersion"]
  fn get_joystick_firmware_version(joystick *Joystick) u16

  @[link_name="SDL_GetJoystickSerial"]
  fn get_joystick_serial(joystick *Joystick) const cstring

  @[link_name="SDL_GetJoystickType"]
  fn get_joystick_type(joystick *Joystick) JoystickType

  @[link_name="SDL_GetJoystickGUIDInfo"]
  fn get_joystick_guid_info(guid GUID, vendor *u16, product *u16, version *u16, crc16 *u16)

  @[link_name="SDL_JoystickConnected"]
  fn joystick_connected(joystick *Joystick) bool

  @[link_name="SDL_GetJoystickID"]
  fn get_joystick_id(joystick *Joystick) JoystickID

  @[link_name="SDL_GetNumJoystickAxes"]
  fn get_num_joystick_axes(joystick *Joystick) libc::c_int

  @[link_name="SDL_GetNumJoystickBalls"]
  fn get_num_joystick_balls(joystick *Joystick) libc::c_int

  @[link_name="SDL_GetNumJoystickHats"]
  fn get_num_joystick_hats(joystick *Joystick) libc::c_int

  @[link_name="SDL_GetNumJoystickButtons"]
  fn get_num_joystick_buttons(joystick *Joystick) libc::c_int

  @[link_name="SDL_SetJoystickEventsEnabled"]
  fn set_joystick_events_enabled(enabled bool)

  @[link_name="SDL_JoystickEventsEnabled"]
  fn joystick_events_enabled() bool

  @[link_name="SDL_UpdateJoysticks"]
  fn update_joysticks()

  @[link_name="SDL_GetJoystickAxis"]
  fn get_joystick_axis(joystick *Joystick, axis libc::c_int) i16

  @[link_name="SDL_GetJoystickAxisInitialState"]
  fn get_joystick_axis_initial_state(joystick *Joystick, axis libc::c_int, state *i16) bool

  @[link_name="SDL_GetJoystickBall"]
  fn get_joystick_ball(joystick *Joystick, ball libc::c_int, dx *libc::c_int, dy *libc::c_int) bool

  @[link_name="SDL_GetJoystickHat"]
  fn get_joystick_hat(joystick *Joystick, hat libc::c_int) u8

  @[link_name="SDL_GetJoystickButton"]
  fn get_joystick_button(joystick *Joystick, button libc::c_int) bool

  @[link_name="SDL_RumbleJoystick"]
  fn rumble_joystick(joystick *Joystick, low_frequency_rumble u16, high_frequency_rumble u16, duration_ms u32) bool

  @[link_name="SDL_RumbleJoystickTriggers"]
  fn rumble_joystick_triggers(joystick *Joystick, left_rumble u16, right_rumble u16, duration_ms u32) bool

  @[link_name="SDL_SetJoystickLED"]
  fn set_joystick_led(joystick *Joystick, red u8, green u8, blue u8) bool

  @[link_name="SDL_SendJoystickEffect"]
  fn send_joystick_effect(joystick *Joystick, data const *Void, size libc::c_int) bool

  @[link_name="SDL_CloseJoystick"]
  fn close_joystick(joystick *Joystick)

  @[link_name="SDL_GetJoystickConnectionState"]
  fn get_joystick_connection_state(joystick *Joystick) JoystickConnectionState

  @[link_name="SDL_GetJoystickPowerInfo"]
  fn get_joystick_power_info(joystick *Joystick, percent *libc::c_int) PowerState

  @[link_name="SDL_AddGamepadMapping"]
  fn add_gamepad_mapping(mapping const cstring) libc::c_int

  @[link_name="SDL_AddGamepadMappingsFromIO"]
  fn add_gamepad_mappings_from_io(src *IOStream, closeio bool) libc::c_int

  @[link_name="SDL_AddGamepadMappingsFromFile"]
  fn add_gamepad_mappings_from_file(file const cstring) libc::c_int

  @[link_name="SDL_ReloadGamepadMappings"]
  fn reload_gamepad_mappings() bool

  @[link_name="SDL_GetGamepadMappings"]
  fn get_gamepad_mappings(count *libc::c_int) *cstring

  @[link_name="SDL_GetGamepadMappingForGUID"]
  fn get_gamepad_mapping_for_guid(guid GUID) cstring

  @[link_name="SDL_GetGamepadMapping"]
  fn get_gamepad_mapping(gamepad *Gamepad) cstring

  @[link_name="SDL_SetGamepadMapping"]
  fn set_gamepad_mapping(instance_id JoystickID, mapping const cstring) bool

  @[link_name="SDL_HasGamepad"]
  fn has_gamepad() bool

  @[link_name="SDL_GetGamepads"]
  fn get_gamepads(count *libc::c_int) *JoystickID

  @[link_name="SDL_IsGamepad"]
  fn is_gamepad(instance_id JoystickID) bool

  @[link_name="SDL_GetGamepadNameForID"]
  fn get_gamepad_name_for_id(instance_id JoystickID) const cstring

  @[link_name="SDL_GetGamepadPathForID"]
  fn get_gamepad_path_for_id(instance_id JoystickID) const cstring

  @[link_name="SDL_GetGamepadPlayerIndexForID"]
  fn get_gamepad_player_index_for_id(instance_id JoystickID) libc::c_int

  @[link_name="SDL_GetGamepadGUIDForID"]
  fn get_gamepad_guid_for_id(instance_id JoystickID) GUID

  @[link_name="SDL_GetGamepadVendorForID"]
  fn get_gamepad_vendor_for_id(instance_id JoystickID) u16

  @[link_name="SDL_GetGamepadProductForID"]
  fn get_gamepad_product_for_id(instance_id JoystickID) u16

  @[link_name="SDL_GetGamepadProductVersionForID"]
  fn get_gamepad_product_version_for_id(instance_id JoystickID) u16

  @[link_name="SDL_GetGamepadTypeForID"]
  fn get_gamepad_type_for_id(instance_id JoystickID) GamepadType

  @[link_name="SDL_GetRealGamepadTypeForID"]
  fn get_real_gamepad_type_for_id(instance_id JoystickID) GamepadType

  @[link_name="SDL_GetGamepadMappingForID"]
  fn get_gamepad_mapping_for_id(instance_id JoystickID) cstring

  @[link_name="SDL_OpenGamepad"]
  fn open_gamepad(instance_id JoystickID) *Gamepad

  @[link_name="SDL_GetGamepadFromID"]
  fn get_gamepad_from_id(instance_id JoystickID) *Gamepad

  @[link_name="SDL_GetGamepadFromPlayerIndex"]
  fn get_gamepad_from_player_index(player_index libc::c_int) *Gamepad

  @[link_name="SDL_GetGamepadProperties"]
  fn get_gamepad_properties(gamepad *Gamepad) PropertiesID

  @[link_name="SDL_GetGamepadID"]
  fn get_gamepad_id(gamepad *Gamepad) JoystickID

  @[link_name="SDL_GetGamepadName"]
  fn get_gamepad_name(gamepad *Gamepad) const cstring

  @[link_name="SDL_GetGamepadPath"]
  fn get_gamepad_path(gamepad *Gamepad) const cstring

  @[link_name="SDL_GetGamepadType"]
  fn get_gamepad_type(gamepad *Gamepad) GamepadType

  @[link_name="SDL_GetRealGamepadType"]
  fn get_real_gamepad_type(gamepad *Gamepad) GamepadType

  @[link_name="SDL_GetGamepadPlayerIndex"]
  fn get_gamepad_player_index(gamepad *Gamepad) libc::c_int

  @[link_name="SDL_SetGamepadPlayerIndex"]
  fn set_gamepad_player_index(gamepad *Gamepad, player_index libc::c_int) bool

  @[link_name="SDL_GetGamepadVendor"]
  fn get_gamepad_vendor(gamepad *Gamepad) u16

  @[link_name="SDL_GetGamepadProduct"]
  fn get_gamepad_product(gamepad *Gamepad) u16

  @[link_name="SDL_GetGamepadProductVersion"]
  fn get_gamepad_product_version(gamepad *Gamepad) u16

  @[link_name="SDL_GetGamepadFirmwareVersion"]
  fn get_gamepad_firmware_version(gamepad *Gamepad) u16

  @[link_name="SDL_GetGamepadSerial"]
  fn get_gamepad_serial(gamepad *Gamepad) const cstring

  @[link_name="SDL_GetGamepadSteamHandle"]
  fn get_gamepad_steam_handle(gamepad *Gamepad) u64

  @[link_name="SDL_GetGamepadConnectionState"]
  fn get_gamepad_connection_state(gamepad *Gamepad) JoystickConnectionState

  @[link_name="SDL_GetGamepadPowerInfo"]
  fn get_gamepad_power_info(gamepad *Gamepad, percent *libc::c_int) PowerState

  @[link_name="SDL_GamepadConnected"]
  fn gamepad_connected(gamepad *Gamepad) bool

  @[link_name="SDL_GetGamepadJoystick"]
  fn get_gamepad_joystick(gamepad *Gamepad) *Joystick

  @[link_name="SDL_SetGamepadEventsEnabled"]
  fn set_gamepad_events_enabled(enabled bool)

  @[link_name="SDL_GamepadEventsEnabled"]
  fn gamepad_events_enabled() bool

  @[link_name="SDL_GetGamepadBindings"]
  fn get_gamepad_bindings(gamepad *Gamepad, count *libc::c_int) **GamepadBinding

  @[link_name="SDL_UpdateGamepads"]
  fn update_gamepads()

  @[link_name="SDL_GetGamepadTypeFromString"]
  fn get_gamepad_type_from_string(str const cstring) GamepadType

  @[link_name="SDL_GetGamepadStringForType"]
  fn get_gamepad_string_for_type(type_ GamepadType) const cstring

  @[link_name="SDL_GetGamepadAxisFromString"]
  fn get_gamepad_axis_from_string(str const cstring) GamepadAxis

  @[link_name="SDL_GetGamepadStringForAxis"]
  fn get_gamepad_string_for_axis(axis GamepadAxis) const cstring

  @[link_name="SDL_GamepadHasAxis"]
  fn gamepad_has_axis(gamepad *Gamepad, axis GamepadAxis) bool

  @[link_name="SDL_GetGamepadAxis"]
  fn get_gamepad_axis(gamepad *Gamepad, axis GamepadAxis) i16

  @[link_name="SDL_GetGamepadButtonFromString"]
  fn get_gamepad_button_from_string(str const cstring) GamepadButton

  @[link_name="SDL_GetGamepadStringForButton"]
  fn get_gamepad_string_for_button(button GamepadButton) const cstring

  @[link_name="SDL_GamepadHasButton"]
  fn gamepad_has_button(gamepad *Gamepad, button GamepadButton) bool

  @[link_name="SDL_GetGamepadButton"]
  fn get_gamepad_button(gamepad *Gamepad, button GamepadButton) bool

  @[link_name="SDL_GetGamepadButtonLabelForType"]
  fn get_gamepad_button_label_for_type(type_ GamepadType, button GamepadButton) GamepadButtonLabel

  @[link_name="SDL_GetGamepadButtonLabel"]
  fn get_gamepad_button_label(gamepad *Gamepad, button GamepadButton) GamepadButtonLabel

  @[link_name="SDL_GetNumGamepadTouchpads"]
  fn get_num_gamepad_touchpads(gamepad *Gamepad) libc::c_int

  @[link_name="SDL_GetNumGamepadTouchpadFingers"]
  fn get_num_gamepad_touchpad_fingers(gamepad *Gamepad, touchpad libc::c_int) libc::c_int

  @[link_name="SDL_GetGamepadTouchpadFinger"]
  fn get_gamepad_touchpad_finger(gamepad *Gamepad, touchpad libc::c_int, finger libc::c_int, down *bool, x *libc::c_float, y *libc::c_float, pressure *libc::c_float) bool

  @[link_name="SDL_GamepadHasSensor"]
  fn gamepad_has_sensor(gamepad *Gamepad, type_ SensorType) bool

  @[link_name="SDL_SetGamepadSensorEnabled"]
  fn set_gamepad_sensor_enabled(gamepad *Gamepad, type_ SensorType, enabled bool) bool

  @[link_name="SDL_GamepadSensorEnabled"]
  fn gamepad_sensor_enabled(gamepad *Gamepad, type_ SensorType) bool

  @[link_name="SDL_GetGamepadSensorDataRate"]
  fn get_gamepad_sensor_data_rate(gamepad *Gamepad, type_ SensorType) libc::c_float

  @[link_name="SDL_GetGamepadSensorData"]
  fn get_gamepad_sensor_data(gamepad *Gamepad, type_ SensorType, data *libc::c_float, num_values libc::c_int) bool

  @[link_name="SDL_RumbleGamepad"]
  fn rumble_gamepad(gamepad *Gamepad, low_frequency_rumble u16, high_frequency_rumble u16, duration_ms u32) bool

  @[link_name="SDL_RumbleGamepadTriggers"]
  fn rumble_gamepad_triggers(gamepad *Gamepad, left_rumble u16, right_rumble u16, duration_ms u32) bool

  @[link_name="SDL_SetGamepadLED"]
  fn set_gamepad_led(gamepad *Gamepad, red u8, green u8, blue u8) bool

  @[link_name="SDL_SendGamepadEffect"]
  fn send_gamepad_effect(gamepad *Gamepad, data const *Void, size libc::c_int) bool

  @[link_name="SDL_CloseGamepad"]
  fn close_gamepad(gamepad *Gamepad)

  @[link_name="SDL_GetGamepadAppleSFSymbolsNameForButton"]
  fn get_gamepad_apple_sf_symbols_name_for_button(gamepad *Gamepad, button GamepadButton) const cstring

  @[link_name="SDL_GetGamepadAppleSFSymbolsNameForAxis"]
  fn get_gamepad_apple_sf_symbols_name_for_axis(gamepad *Gamepad, axis GamepadAxis) const cstring

  @[link_name="SDL_HasKeyboard"]
  fn has_keyboard() bool

  @[link_name="SDL_GetKeyboards"]
  fn get_keyboards(count *libc::c_int) *KeyboardID

  @[link_name="SDL_GetKeyboardNameForID"]
  fn get_keyboard_name_for_id(instance_id KeyboardID) const cstring

  @[link_name="SDL_GetKeyboardFocus"]
  fn get_keyboard_focus() *Window

  @[link_name="SDL_GetKeyboardState"]
  fn get_keyboard_state(numkeys *libc::c_int) const *bool

  @[link_name="SDL_ResetKeyboard"]
  fn reset_keyboard()

  @[link_name="SDL_GetModState"]
  fn get_mod_state() Keymod

  @[link_name="SDL_SetModState"]
  fn set_mod_state(modstate Keymod)

  @[link_name="SDL_GetKeyFromScancode"]
  fn get_key_from_scancode(scancode Scancode, modstate Keymod, key_event bool) Keycode

  @[link_name="SDL_GetScancodeFromKey"]
  fn get_scancode_from_key(key Keycode, modstate *Keymod) Scancode

  @[link_name="SDL_SetScancodeName"]
  fn set_scancode_name(scancode Scancode, name const cstring) bool

  @[link_name="SDL_GetScancodeName"]
  fn get_scancode_name(scancode Scancode) const cstring

  @[link_name="SDL_GetScancodeFromName"]
  fn get_scancode_from_name(name const cstring) Scancode

  @[link_name="SDL_GetKeyName"]
  fn get_key_name(key Keycode) const cstring

  @[link_name="SDL_GetKeyFromName"]
  fn get_key_from_name(name const cstring) Keycode

  @[link_name="SDL_StartTextInput"]
  fn start_text_input(window *Window) bool

  @[link_name="SDL_StartTextInputWithProperties"]
  fn start_text_input_with_properties(window *Window, props PropertiesID) bool

  @[link_name="SDL_TextInputActive"]
  fn text_input_active(window *Window) bool

  @[link_name="SDL_StopTextInput"]
  fn stop_text_input(window *Window) bool

  @[link_name="SDL_ClearComposition"]
  fn clear_composition(window *Window) bool

  @[link_name="SDL_SetTextInputArea"]
  fn set_text_input_area(window *Window, rect const *Rect, cursor libc::c_int) bool

  @[link_name="SDL_GetTextInputArea"]
  fn get_text_input_area(window *Window, rect *Rect, cursor *libc::c_int) bool

  @[link_name="SDL_HasScreenKeyboardSupport"]
  fn has_screen_keyboard_support() bool

  @[link_name="SDL_ScreenKeyboardShown"]
  fn screen_keyboard_shown(window *Window) bool

  @[link_name="SDL_HasMouse"]
  fn has_mouse() bool

  @[link_name="SDL_GetMice"]
  fn get_mice(count *libc::c_int) *MouseID

  @[link_name="SDL_GetMouseNameForID"]
  fn get_mouse_name_for_id(instance_id MouseID) const cstring

  @[link_name="SDL_GetMouseFocus"]
  fn get_mouse_focus() *Window

  @[link_name="SDL_GetMouseState"]
  fn get_mouse_state(x *libc::c_float, y *libc::c_float) MouseButtonFlags

  @[link_name="SDL_GetGlobalMouseState"]
  fn get_global_mouse_state(x *libc::c_float, y *libc::c_float) MouseButtonFlags

  @[link_name="SDL_GetRelativeMouseState"]
  fn get_relative_mouse_state(x *libc::c_float, y *libc::c_float) MouseButtonFlags

  @[link_name="SDL_WarpMouseInWindow"]
  fn warp_mouse_in_window(window *Window, x libc::c_float, y libc::c_float)

  @[link_name="SDL_WarpMouseGlobal"]
  fn warp_mouse_global(x libc::c_float, y libc::c_float) bool

  @[link_name="SDL_SetWindowRelativeMouseMode"]
  fn set_window_relative_mouse_mode(window *Window, enabled bool) bool

  @[link_name="SDL_GetWindowRelativeMouseMode"]
  fn get_window_relative_mouse_mode(window *Window) bool

  @[link_name="SDL_CaptureMouse"]
  fn capture_mouse(enabled bool) bool

  @[link_name="SDL_CreateCursor"]
  fn create_cursor(data const *u8, mask const *u8, w libc::c_int, h libc::c_int, hot_x libc::c_int, hot_y libc::c_int) *Cursor

  @[link_name="SDL_CreateColorCursor"]
  fn create_color_cursor(surface *Surface, hot_x libc::c_int, hot_y libc::c_int) *Cursor

  @[link_name="SDL_CreateSystemCursor"]
  fn create_system_cursor(id SystemCursor) *Cursor

  @[link_name="SDL_SetCursor"]
  fn set_cursor(cursor *Cursor) bool

  @[link_name="SDL_GetCursor"]
  fn get_cursor() *Cursor

  @[link_name="SDL_GetDefaultCursor"]
  fn get_default_cursor() *Cursor

  @[link_name="SDL_DestroyCursor"]
  fn destroy_cursor(cursor *Cursor)

  @[link_name="SDL_ShowCursor"]
  fn show_cursor() bool

  @[link_name="SDL_HideCursor"]
  fn hide_cursor() bool

  @[link_name="SDL_CursorVisible"]
  fn cursor_visible() bool

  @[link_name="SDL_GetTouchDevices"]
  fn get_touch_devices(count *libc::c_int) *TouchID

  @[link_name="SDL_GetTouchDeviceName"]
  fn get_touch_device_name(touchID TouchID) const cstring

  @[link_name="SDL_GetTouchDeviceType"]
  fn get_touch_device_type(touchID TouchID) TouchDeviceType

  @[link_name="SDL_GetTouchFingers"]
  fn get_touch_fingers(touchID TouchID, count *libc::c_int) **Finger

  @[link_name="SDL_PumpEvents"]
  fn pump_events()

  @[link_name="SDL_PeepEvents"]
  fn peep_events(events *Event, numevents libc::c_int, action EventAction, minType u32, maxType u32) libc::c_int

  @[link_name="SDL_HasEvent"]
  fn has_event(type_ u32) bool

  @[link_name="SDL_HasEvents"]
  fn has_events(minType u32, maxType u32) bool

  @[link_name="SDL_FlushEvent"]
  fn flush_event(type_ u32)

  @[link_name="SDL_FlushEvents"]
  fn flush_events(minType u32, maxType u32)

  @[link_name="SDL_PollEvent"]
  fn poll_event(event *Event) bool

  @[link_name="SDL_WaitEvent"]
  fn wait_event(event *Event) bool

  @[link_name="SDL_WaitEventTimeout"]
  fn wait_event_timeout(event *Event, timeoutMS i32) bool

  @[link_name="SDL_PushEvent"]
  fn push_event(event *Event) bool

  @[link_name="SDL_SetEventFilter"]
  fn set_event_filter(filter EventFilter, userdata rawptr)

  @[link_name="SDL_GetEventFilter"]
  fn get_event_filter(filter *EventFilter, userdata *rawptr) bool

  @[link_name="SDL_AddEventWatch"]
  fn add_event_watch(filter EventFilter, userdata rawptr) bool

  @[link_name="SDL_RemoveEventWatch"]
  fn remove_event_watch(filter EventFilter, userdata rawptr)

  @[link_name="SDL_FilterEvents"]
  fn filter_events(filter EventFilter, userdata rawptr)

  @[link_name="SDL_SetEventEnabled"]
  fn set_event_enabled(type_ u32, enabled bool)

  @[link_name="SDL_EventEnabled"]
  fn event_enabled(type_ u32) bool

  @[link_name="SDL_RegisterEvents"]
  fn register_events(numevents libc::c_int) u32

  @[link_name="SDL_GetWindowFromEvent"]
  fn get_window_from_event(event const *Event) *Window

  @[link_name="SDL_GetBasePath"]
  fn get_base_path() const cstring

  @[link_name="SDL_GetPrefPath"]
  fn get_pref_path(org const cstring, app const cstring) cstring

  @[link_name="SDL_GetUserFolder"]
  fn get_user_folder(folder Folder) const cstring

  @[link_name="SDL_CreateDirectory"]
  fn create_directory(path const cstring) bool

  @[link_name="SDL_EnumerateDirectory"]
  fn enumerate_directory(path const cstring, callback EnumerateDirectoryCallback, userdata rawptr) bool

  @[link_name="SDL_RemovePath"]
  fn remove_path(path const cstring) bool

  @[link_name="SDL_RenamePath"]
  fn rename_path(oldpath const cstring, newpath const cstring) bool

  @[link_name="SDL_CopyFile"]
  fn copy_file(oldpath const cstring, newpath const cstring) bool

  @[link_name="SDL_GetPathInfo"]
  fn get_path_info(path const cstring, info *PathInfo) bool

  @[link_name="SDL_GlobDirectory"]
  fn glob_directory(path const cstring, pattern const cstring, flags GlobFlags, count *libc::c_int) *cstring

  @[link_name="SDL_GetCurrentDirectory"]
  fn get_current_directory() cstring

  @[link_name="SDL_GPUSupportsShaderFormats"]
  fn gpu_supports_shader_formats(format_flags GPUShaderFormat, name const cstring) bool

  @[link_name="SDL_GPUSupportsProperties"]
  fn gpu_supports_properties(props PropertiesID) bool

  @[link_name="SDL_CreateGPUDevice"]
  fn create_gpu_device(format_flags GPUShaderFormat, debug_mode bool, name const cstring) *GPUDevice

  @[link_name="SDL_CreateGPUDeviceWithProperties"]
  fn create_gpu_device_with_properties(props PropertiesID) *GPUDevice

  @[link_name="SDL_DestroyGPUDevice"]
  fn destroy_gpu_device(device *GPUDevice)

  @[link_name="SDL_GetNumGPUDrivers"]
  fn get_num_gpu_drivers() libc::c_int

  @[link_name="SDL_GetGPUDriver"]
  fn get_gpu_driver(index libc::c_int) const cstring

  @[link_name="SDL_GetGPUDeviceDriver"]
  fn get_gpu_device_driver(device *GPUDevice) const cstring

  @[link_name="SDL_GetGPUShaderFormats"]
  fn get_gpu_shader_formats(device *GPUDevice) GPUShaderFormat

  @[link_name="SDL_CreateGPUComputePipeline"]
  fn create_gpu_compute_pipeline(device *GPUDevice, createinfo const *GPUComputePipelineCreateInfo) *GPUComputePipeline

  @[link_name="SDL_CreateGPUGraphicsPipeline"]
  fn create_gpu_graphics_pipeline(device *GPUDevice, createinfo const *GPUGraphicsPipelineCreateInfo) *GPUGraphicsPipeline

  @[link_name="SDL_CreateGPUSampler"]
  fn create_gpu_sampler(device *GPUDevice, createinfo const *GPUSamplerCreateInfo) *GPUSampler

  @[link_name="SDL_CreateGPUShader"]
  fn create_gpu_shader(device *GPUDevice, createinfo const *GPUShaderCreateInfo) *GPUShader

  @[link_name="SDL_CreateGPUTexture"]
  fn create_gpu_texture(device *GPUDevice, createinfo const *GPUTextureCreateInfo) *GPUTexture

  @[link_name="SDL_CreateGPUBuffer"]
  fn create_gpu_buffer(device *GPUDevice, createinfo const *GPUBufferCreateInfo) *GPUBuffer

  @[link_name="SDL_CreateGPUTransferBuffer"]
  fn create_gpu_transfer_buffer(device *GPUDevice, createinfo const *GPUTransferBufferCreateInfo) *GPUTransferBuffer

  @[link_name="SDL_SetGPUBufferName"]
  fn set_gpu_buffer_name(device *GPUDevice, buffer *GPUBuffer, text const cstring)

  @[link_name="SDL_SetGPUTextureName"]
  fn set_gpu_texture_name(device *GPUDevice, texture *GPUTexture, text const cstring)

  @[link_name="SDL_InsertGPUDebugLabel"]
  fn insert_gpu_debug_label(command_buffer *GPUCommandBuffer, text const cstring)

  @[link_name="SDL_PushGPUDebugGroup"]
  fn push_gpu_debug_group(command_buffer *GPUCommandBuffer, name const cstring)

  @[link_name="SDL_PopGPUDebugGroup"]
  fn pop_gpu_debug_group(command_buffer *GPUCommandBuffer)

  @[link_name="SDL_ReleaseGPUTexture"]
  fn release_gpu_texture(device *GPUDevice, texture *GPUTexture)

  @[link_name="SDL_ReleaseGPUSampler"]
  fn release_gpu_sampler(device *GPUDevice, sampler *GPUSampler)

  @[link_name="SDL_ReleaseGPUBuffer"]
  fn release_gpu_buffer(device *GPUDevice, buffer *GPUBuffer)

  @[link_name="SDL_ReleaseGPUTransferBuffer"]
  fn release_gpu_transfer_buffer(device *GPUDevice, transfer_buffer *GPUTransferBuffer)

  @[link_name="SDL_ReleaseGPUComputePipeline"]
  fn release_gpu_compute_pipeline(device *GPUDevice, compute_pipeline *GPUComputePipeline)

  @[link_name="SDL_ReleaseGPUShader"]
  fn release_gpu_shader(device *GPUDevice, shader *GPUShader)

  @[link_name="SDL_ReleaseGPUGraphicsPipeline"]
  fn release_gpu_graphics_pipeline(device *GPUDevice, graphics_pipeline *GPUGraphicsPipeline)

  @[link_name="SDL_AcquireGPUCommandBuffer"]
  fn acquire_gpu_command_buffer(device *GPUDevice) *GPUCommandBuffer

  @[link_name="SDL_PushGPUVertexUniformData"]
  fn push_gpu_vertex_uniform_data(command_buffer *GPUCommandBuffer, slot_index u32, data const *Void, length u32)

  @[link_name="SDL_PushGPUFragmentUniformData"]
  fn push_gpu_fragment_uniform_data(command_buffer *GPUCommandBuffer, slot_index u32, data const *Void, length u32)

  @[link_name="SDL_PushGPUComputeUniformData"]
  fn push_gpu_compute_uniform_data(command_buffer *GPUCommandBuffer, slot_index u32, data const *Void, length u32)

  @[link_name="SDL_BeginGPURenderPass"]
  fn begin_gpu_render_pass(command_buffer *GPUCommandBuffer, color_target_infos const *GPUColorTargetInfo, num_color_targets u32, depth_stencil_target_info const *GPUDepthStencilTargetInfo) *GPURenderPass

  @[link_name="SDL_BindGPUGraphicsPipeline"]
  fn bind_gpu_graphics_pipeline(render_pass *GPURenderPass, graphics_pipeline *GPUGraphicsPipeline)

  @[link_name="SDL_SetGPUViewport"]
  fn set_gpu_viewport(render_pass *GPURenderPass, viewport const *GPUViewport)

  @[link_name="SDL_SetGPUScissor"]
  fn set_gpu_scissor(render_pass *GPURenderPass, scissor const *Rect)

  @[link_name="SDL_SetGPUBlendConstants"]
  fn set_gpu_blend_constants(render_pass *GPURenderPass, blend_constants FColor)

  @[link_name="SDL_SetGPUStencilReference"]
  fn set_gpu_stencil_reference(render_pass *GPURenderPass, reference u8)

  @[link_name="SDL_BindGPUVertexBuffers"]
  fn bind_gpu_vertex_buffers(render_pass *GPURenderPass, first_slot u32, bindings const *GPUBufferBinding, num_bindings u32)

  @[link_name="SDL_BindGPUIndexBuffer"]
  fn bind_gpu_index_buffer(render_pass *GPURenderPass, binding const *GPUBufferBinding, index_element_size GPUIndexElementSize)

  @[link_name="SDL_BindGPUVertexSamplers"]
  fn bind_gpu_vertex_samplers(render_pass *GPURenderPass, first_slot u32, texture_sampler_bindings const *GPUTextureSamplerBinding, num_bindings u32)

  @[link_name="SDL_BindGPUVertexStorageTextures"]
  fn bind_gpu_vertex_storage_textures(render_pass *GPURenderPass, first_slot u32, storage_textures const **GPUTexture, num_bindings u32)

  @[link_name="SDL_BindGPUVertexStorageBuffers"]
  fn bind_gpu_vertex_storage_buffers(render_pass *GPURenderPass, first_slot u32, storage_buffers const **GPUBuffer, num_bindings u32)

  @[link_name="SDL_BindGPUFragmentSamplers"]
  fn bind_gpu_fragment_samplers(render_pass *GPURenderPass, first_slot u32, texture_sampler_bindings const *GPUTextureSamplerBinding, num_bindings u32)

  @[link_name="SDL_BindGPUFragmentStorageTextures"]
  fn bind_gpu_fragment_storage_textures(render_pass *GPURenderPass, first_slot u32, storage_textures const **GPUTexture, num_bindings u32)

  @[link_name="SDL_BindGPUFragmentStorageBuffers"]
  fn bind_gpu_fragment_storage_buffers(render_pass *GPURenderPass, first_slot u32, storage_buffers const **GPUBuffer, num_bindings u32)

  @[link_name="SDL_DrawGPUIndexedPrimitives"]
  fn draw_gpu_indexed_primitives(render_pass *GPURenderPass, num_indices u32, num_instances u32, first_index u32, vertex_offset i32, first_instance u32)

  @[link_name="SDL_DrawGPUPrimitives"]
  fn draw_gpu_primitives(render_pass *GPURenderPass, num_vertices u32, num_instances u32, first_vertex u32, first_instance u32)

  @[link_name="SDL_DrawGPUPrimitivesIndirect"]
  fn draw_gpu_primitives_indirect(render_pass *GPURenderPass, buffer *GPUBuffer, offset u32, draw_count u32)

  @[link_name="SDL_DrawGPUIndexedPrimitivesIndirect"]
  fn draw_gpu_indexed_primitives_indirect(render_pass *GPURenderPass, buffer *GPUBuffer, offset u32, draw_count u32)

  @[link_name="SDL_EndGPURenderPass"]
  fn end_gpu_render_pass(render_pass *GPURenderPass)

  @[link_name="SDL_BeginGPUComputePass"]
  fn begin_gpu_compute_pass(command_buffer *GPUCommandBuffer, storage_texture_bindings const *GPUStorageTextureReadWriteBinding, num_storage_texture_bindings u32, storage_buffer_bindings const *GPUStorageBufferReadWriteBinding, num_storage_buffer_bindings u32) *GPUComputePass

  @[link_name="SDL_BindGPUComputePipeline"]
  fn bind_gpu_compute_pipeline(compute_pass *GPUComputePass, compute_pipeline *GPUComputePipeline)

  @[link_name="SDL_BindGPUComputeSamplers"]
  fn bind_gpu_compute_samplers(compute_pass *GPUComputePass, first_slot u32, texture_sampler_bindings const *GPUTextureSamplerBinding, num_bindings u32)

  @[link_name="SDL_BindGPUComputeStorageTextures"]
  fn bind_gpu_compute_storage_textures(compute_pass *GPUComputePass, first_slot u32, storage_textures const **GPUTexture, num_bindings u32)

  @[link_name="SDL_BindGPUComputeStorageBuffers"]
  fn bind_gpu_compute_storage_buffers(compute_pass *GPUComputePass, first_slot u32, storage_buffers const **GPUBuffer, num_bindings u32)

  @[link_name="SDL_DispatchGPUCompute"]
  fn dispatch_gpu_compute(compute_pass *GPUComputePass, groupcount_x u32, groupcount_y u32, groupcount_z u32)

  @[link_name="SDL_DispatchGPUComputeIndirect"]
  fn dispatch_gpu_compute_indirect(compute_pass *GPUComputePass, buffer *GPUBuffer, offset u32)

  @[link_name="SDL_EndGPUComputePass"]
  fn end_gpu_compute_pass(compute_pass *GPUComputePass)

  @[link_name="SDL_MapGPUTransferBuffer"]
  fn map_gpu_transfer_buffer(device *GPUDevice, transfer_buffer *GPUTransferBuffer, cycle bool) rawptr

  @[link_name="SDL_UnmapGPUTransferBuffer"]
  fn unmap_gpu_transfer_buffer(device *GPUDevice, transfer_buffer *GPUTransferBuffer)

  @[link_name="SDL_BeginGPUCopyPass"]
  fn begin_gpu_copy_pass(command_buffer *GPUCommandBuffer) *GPUCopyPass

  @[link_name="SDL_UploadToGPUTexture"]
  fn upload_to_gpu_texture(copy_pass *GPUCopyPass, source const *GPUTextureTransferInfo, destination const *GPUTextureRegion, cycle bool)

  @[link_name="SDL_UploadToGPUBuffer"]
  fn upload_to_gpu_buffer(copy_pass *GPUCopyPass, source const *GPUTransferBufferLocation, destination const *GPUBufferRegion, cycle bool)

  @[link_name="SDL_CopyGPUTextureToTexture"]
  fn copy_gpu_texture_to_texture(copy_pass *GPUCopyPass, source const *GPUTextureLocation, destination const *GPUTextureLocation, w u32, h u32, d u32, cycle bool)

  @[link_name="SDL_CopyGPUBufferToBuffer"]
  fn copy_gpu_buffer_to_buffer(copy_pass *GPUCopyPass, source const *GPUBufferLocation, destination const *GPUBufferLocation, size u32, cycle bool)

  @[link_name="SDL_DownloadFromGPUTexture"]
  fn download_from_gpu_texture(copy_pass *GPUCopyPass, source const *GPUTextureRegion, destination const *GPUTextureTransferInfo)

  @[link_name="SDL_DownloadFromGPUBuffer"]
  fn download_from_gpu_buffer(copy_pass *GPUCopyPass, source const *GPUBufferRegion, destination const *GPUTransferBufferLocation)

  @[link_name="SDL_EndGPUCopyPass"]
  fn end_gpu_copy_pass(copy_pass *GPUCopyPass)

  @[link_name="SDL_GenerateMipmapsForGPUTexture"]
  fn generate_mipmaps_for_gpu_texture(command_buffer *GPUCommandBuffer, texture *GPUTexture)

  @[link_name="SDL_BlitGPUTexture"]
  fn blit_gpu_texture(command_buffer *GPUCommandBuffer, info const *GPUBlitInfo)

  @[link_name="SDL_WindowSupportsGPUSwapchainComposition"]
  fn window_supports_gpu_swapchain_composition(device *GPUDevice, window *Window, swapchain_composition GPUSwapchainComposition) bool

  @[link_name="SDL_WindowSupportsGPUPresentMode"]
  fn window_supports_gpu_present_mode(device *GPUDevice, window *Window, present_mode GPUPresentMode) bool

  @[link_name="SDL_ClaimWindowForGPUDevice"]
  fn claim_window_for_gpu_device(device *GPUDevice, window *Window) bool

  @[link_name="SDL_ReleaseWindowFromGPUDevice"]
  fn release_window_from_gpu_device(device *GPUDevice, window *Window)

  @[link_name="SDL_SetGPUSwapchainParameters"]
  fn set_gpu_swapchain_parameters(device *GPUDevice, window *Window, swapchain_composition GPUSwapchainComposition, present_mode GPUPresentMode) bool

  @[link_name="SDL_SetGPUAllowedFramesInFlight"]
  fn set_gpu_allowed_frames_in_flight(device *GPUDevice, allowed_frames_in_flight u32) bool

  @[link_name="SDL_GetGPUSwapchainTextureFormat"]
  fn get_gpu_swapchain_texture_format(device *GPUDevice, window *Window) GPUTextureFormat

  @[link_name="SDL_AcquireGPUSwapchainTexture"]
  fn acquire_gpu_swapchain_texture(command_buffer *GPUCommandBuffer, window *Window, swapchain_texture **GPUTexture, swapchain_texture_width *u32, swapchain_texture_height *u32) bool

  @[link_name="SDL_WaitForGPUSwapchain"]
  fn wait_for_gpu_swapchain(device *GPUDevice, window *Window) bool

  @[link_name="SDL_WaitAndAcquireGPUSwapchainTexture"]
  fn wait_and_acquire_gpu_swapchain_texture(command_buffer *GPUCommandBuffer, window *Window, swapchain_texture **GPUTexture, swapchain_texture_width *u32, swapchain_texture_height *u32) bool

  @[link_name="SDL_SubmitGPUCommandBuffer"]
  fn submit_gpu_command_buffer(command_buffer *GPUCommandBuffer) bool

  @[link_name="SDL_SubmitGPUCommandBufferAndAcquireFence"]
  fn submit_gpu_command_buffer_and_acquire_fence(command_buffer *GPUCommandBuffer) *GPUFence

  @[link_name="SDL_CancelGPUCommandBuffer"]
  fn cancel_gpu_command_buffer(command_buffer *GPUCommandBuffer) bool

  @[link_name="SDL_WaitForGPUIdle"]
  fn wait_for_gpu_idle(device *GPUDevice) bool

  @[link_name="SDL_WaitForGPUFences"]
  fn wait_for_gpu_fences(device *GPUDevice, wait_all bool, fences const **GPUFence, num_fences u32) bool

  @[link_name="SDL_QueryGPUFence"]
  fn query_gpu_fence(device *GPUDevice, fence *GPUFence) bool

  @[link_name="SDL_ReleaseGPUFence"]
  fn release_gpu_fence(device *GPUDevice, fence *GPUFence)

  @[link_name="SDL_GPUTextureFormatTexelBlockSize"]
  fn gpu_texture_format_texel_block_size(format GPUTextureFormat) u32

  @[link_name="SDL_GPUTextureSupportsFormat"]
  fn gpu_texture_supports_format(device *GPUDevice, format GPUTextureFormat, type_ GPUTextureType, usage GPUTextureUsageFlags) bool

  @[link_name="SDL_GPUTextureSupportsSampleCount"]
  fn gpu_texture_supports_sample_count(device *GPUDevice, format GPUTextureFormat, sample_count GPUSampleCount) bool

  @[link_name="SDL_CalculateGPUTextureFormatSize"]
  fn calculate_gpu_texture_format_size(format GPUTextureFormat, width u32, height u32, depth_or_layer_count u32) u32

  @[link_name="SDL_GetHaptics"]
  fn get_haptics(count *libc::c_int) *HapticID

  @[link_name="SDL_GetHapticNameForID"]
  fn get_haptic_name_for_id(instance_id HapticID) const cstring

  @[link_name="SDL_OpenHaptic"]
  fn open_haptic(instance_id HapticID) *Haptic

  @[link_name="SDL_GetHapticFromID"]
  fn get_haptic_from_id(instance_id HapticID) *Haptic

  @[link_name="SDL_GetHapticID"]
  fn get_haptic_id(haptic *Haptic) HapticID

  @[link_name="SDL_GetHapticName"]
  fn get_haptic_name(haptic *Haptic) const cstring

  @[link_name="SDL_IsMouseHaptic"]
  fn is_mouse_haptic() bool

  @[link_name="SDL_OpenHapticFromMouse"]
  fn open_haptic_from_mouse() *Haptic

  @[link_name="SDL_IsJoystickHaptic"]
  fn is_joystick_haptic(joystick *Joystick) bool

  @[link_name="SDL_OpenHapticFromJoystick"]
  fn open_haptic_from_joystick(joystick *Joystick) *Haptic

  @[link_name="SDL_CloseHaptic"]
  fn close_haptic(haptic *Haptic)

  @[link_name="SDL_GetMaxHapticEffects"]
  fn get_max_haptic_effects(haptic *Haptic) libc::c_int

  @[link_name="SDL_GetMaxHapticEffectsPlaying"]
  fn get_max_haptic_effects_playing(haptic *Haptic) libc::c_int

  @[link_name="SDL_GetHapticFeatures"]
  fn get_haptic_features(haptic *Haptic) u32

  @[link_name="SDL_GetNumHapticAxes"]
  fn get_num_haptic_axes(haptic *Haptic) libc::c_int

  @[link_name="SDL_HapticEffectSupported"]
  fn haptic_effect_supported(haptic *Haptic, effect const *HapticEffect) bool

  @[link_name="SDL_CreateHapticEffect"]
  fn create_haptic_effect(haptic *Haptic, effect const *HapticEffect) libc::c_int

  @[link_name="SDL_UpdateHapticEffect"]
  fn update_haptic_effect(haptic *Haptic, effect libc::c_int, data const *HapticEffect) bool

  @[link_name="SDL_RunHapticEffect"]
  fn run_haptic_effect(haptic *Haptic, effect libc::c_int, iterations u32) bool

  @[link_name="SDL_StopHapticEffect"]
  fn stop_haptic_effect(haptic *Haptic, effect libc::c_int) bool

  @[link_name="SDL_DestroyHapticEffect"]
  fn destroy_haptic_effect(haptic *Haptic, effect libc::c_int)

  @[link_name="SDL_GetHapticEffectStatus"]
  fn get_haptic_effect_status(haptic *Haptic, effect libc::c_int) bool

  @[link_name="SDL_SetHapticGain"]
  fn set_haptic_gain(haptic *Haptic, gain libc::c_int) bool

  @[link_name="SDL_SetHapticAutocenter"]
  fn set_haptic_autocenter(haptic *Haptic, autocenter libc::c_int) bool

  @[link_name="SDL_PauseHaptic"]
  fn pause_haptic(haptic *Haptic) bool

  @[link_name="SDL_ResumeHaptic"]
  fn resume_haptic(haptic *Haptic) bool

  @[link_name="SDL_StopHapticEffects"]
  fn stop_haptic_effects(haptic *Haptic) bool

  @[link_name="SDL_HapticRumbleSupported"]
  fn haptic_rumble_supported(haptic *Haptic) bool

  @[link_name="SDL_InitHapticRumble"]
  fn init_haptic_rumble(haptic *Haptic) bool

  @[link_name="SDL_PlayHapticRumble"]
  fn play_haptic_rumble(haptic *Haptic, strength libc::c_float, length u32) bool

  @[link_name="SDL_StopHapticRumble"]
  fn stop_haptic_rumble(haptic *Haptic) bool

  @[link_name="SDL_hid_init"]
  fn hid_init() libc::c_int

  @[link_name="SDL_hid_exit"]
  fn hid_exit() libc::c_int

  @[link_name="SDL_hid_device_change_count"]
  fn hid_device_change_count() u32

  @[link_name="SDL_hid_enumerate"]
  fn hid_enumerate(vendor_id libc::c_ushort, product_id libc::c_ushort) *hid_device_info

  @[link_name="SDL_hid_free_enumeration"]
  fn hid_free_enumeration(devs *hid_device_info)

  @[link_name="SDL_hid_open"]
  fn hid_open(vendor_id libc::c_ushort, product_id libc::c_ushort, serial_number const *libc::c_wchar_t) *hid_device

  @[link_name="SDL_hid_open_path"]
  fn hid_open_path(path const cstring) *hid_device

  @[link_name="SDL_hid_write"]
  fn hid_write(dev *hid_device, data const *libc::c_uchar, length libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_read_timeout"]
  fn hid_read_timeout(dev *hid_device, data *libc::c_uchar, length libc::c_size_t, milliseconds libc::c_int) libc::c_int

  @[link_name="SDL_hid_read"]
  fn hid_read(dev *hid_device, data *libc::c_uchar, length libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_set_nonblocking"]
  fn hid_set_nonblocking(dev *hid_device, nonblock libc::c_int) libc::c_int

  @[link_name="SDL_hid_send_feature_report"]
  fn hid_send_feature_report(dev *hid_device, data const *libc::c_uchar, length libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_get_feature_report"]
  fn hid_get_feature_report(dev *hid_device, data *libc::c_uchar, length libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_get_input_report"]
  fn hid_get_input_report(dev *hid_device, data *libc::c_uchar, length libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_close"]
  fn hid_close(dev *hid_device) libc::c_int

  @[link_name="SDL_hid_get_manufacturer_string"]
  fn hid_get_manufacturer_string(dev *hid_device, string_ *libc::c_wchar_t, maxlen libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_get_product_string"]
  fn hid_get_product_string(dev *hid_device, string_ *libc::c_wchar_t, maxlen libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_get_serial_number_string"]
  fn hid_get_serial_number_string(dev *hid_device, string_ *libc::c_wchar_t, maxlen libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_get_indexed_string"]
  fn hid_get_indexed_string(dev *hid_device, string_index libc::c_int, string_ *libc::c_wchar_t, maxlen libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_get_device_info"]
  fn hid_get_device_info(dev *hid_device) *hid_device_info

  @[link_name="SDL_hid_get_report_descriptor"]
  fn hid_get_report_descriptor(dev *hid_device, buf *libc::c_uchar, buf_size libc::c_size_t) libc::c_int

  @[link_name="SDL_hid_ble_scan"]
  fn hid_ble_scan(active bool)

  @[link_name="SDL_SetHintWithPriority"]
  fn set_hint_with_priority(name const cstring, value const cstring, priority HintPriority) bool

  @[link_name="SDL_SetHint"]
  fn set_hint(name const cstring, value const cstring) bool

  @[link_name="SDL_ResetHint"]
  fn reset_hint(name const cstring) bool

  @[link_name="SDL_ResetHints"]
  fn reset_hints()

  @[link_name="SDL_GetHint"]
  fn get_hint(name const cstring) const cstring

  @[link_name="SDL_GetHintBoolean"]
  fn get_hint_boolean(name const cstring, default_value bool) bool

  @[link_name="SDL_AddHintCallback"]
  fn add_hint_callback(name const cstring, callback HintCallback, userdata rawptr) bool

  @[link_name="SDL_RemoveHintCallback"]
  fn remove_hint_callback(name const cstring, callback HintCallback, userdata rawptr)

  @[link_name="SDL_Init"]
  fn init(flags InitFlags) bool

  @[link_name="SDL_InitSubSystem"]
  fn init_sub_system(flags InitFlags) bool

  @[link_name="SDL_QuitSubSystem"]
  fn quit_sub_system(flags InitFlags)

  @[link_name="SDL_WasInit"]
  fn was_init(flags InitFlags) InitFlags

  @[link_name="SDL_Quit"]
  fn quit()

  @[link_name="SDL_IsMainThread"]
  fn is_main_thread() bool

  @[link_name="SDL_RunOnMainThread"]
  fn run_on_main_thread(callback MainThreadCallback, userdata rawptr, wait_complete bool) bool

  @[link_name="SDL_SetAppMetadata"]
  fn set_app_metadata(appname const cstring, appversion const cstring, appidentifier const cstring) bool

  @[link_name="SDL_SetAppMetadataProperty"]
  fn set_app_metadata_property(name const cstring, value const cstring) bool

  @[link_name="SDL_GetAppMetadataProperty"]
  fn get_app_metadata_property(name const cstring) const cstring

  @[link_name="SDL_LoadObject"]
  fn load_object(sofile const cstring) *SharedObject

  @[link_name="SDL_LoadFunction"]
  fn load_function(handle *SharedObject, name const cstring) FunctionPointer

  @[link_name="SDL_UnloadObject"]
  fn unload_object(handle *SharedObject)

  @[link_name="SDL_GetPreferredLocales"]
  fn get_preferred_locales(count *libc::c_int) **Locale

  @[link_name="SDL_SetLogPriorities"]
  fn set_log_priorities(priority LogPriority)

  @[link_name="SDL_SetLogPriority"]
  fn set_log_priority(category libc::c_int, priority LogPriority)

  @[link_name="SDL_GetLogPriority"]
  fn get_log_priority(category libc::c_int) LogPriority

  @[link_name="SDL_ResetLogPriorities"]
  fn reset_log_priorities()

  @[link_name="SDL_SetLogPriorityPrefix"]
  fn set_log_priority_prefix(priority LogPriority, prefix const cstring) bool

  @[link_name="SDL_Log"]
  fn log(fmt const cstring, args @c ...)

  @[link_name="SDL_LogTrace"]
  fn log_trace(category libc::c_int, fmt const cstring, args @c ...)

  @[link_name="SDL_LogVerbose"]
  fn log_verbose(category libc::c_int, fmt const cstring, args @c ...)

  @[link_name="SDL_LogDebug"]
  fn log_debug(category libc::c_int, fmt const cstring, args @c ...)

  @[link_name="SDL_LogInfo"]
  fn log_info(category libc::c_int, fmt const cstring, args @c ...)

  @[link_name="SDL_LogWarn"]
  fn log_warn(category libc::c_int, fmt const cstring, args @c ...)

  @[link_name="SDL_LogError"]
  fn log_error(category libc::c_int, fmt const cstring, args @c ...)

  @[link_name="SDL_LogCritical"]
  fn log_critical(category libc::c_int, fmt const cstring, args @c ...)

  @[link_name="SDL_LogMessage"]
  fn log_message(category libc::c_int, priority LogPriority, fmt const cstring, args @c ...)

  @[link_name="SDL_LogMessageV"]
  fn log_message_v(category libc::c_int, priority LogPriority, fmt const cstring, ap libc::c_va_list)

  @[link_name="SDL_GetDefaultLogOutputFunction"]
  fn get_default_log_output_function() LogOutputFunction

  @[link_name="SDL_GetLogOutputFunction"]
  fn get_log_output_function(callback *LogOutputFunction, userdata *rawptr)

  @[link_name="SDL_SetLogOutputFunction"]
  fn set_log_output_function(callback LogOutputFunction, userdata rawptr)

  @[link_name="SDL_ShowMessageBox"]
  fn show_message_box(messageboxdata const *MessageBoxData, buttonid *libc::c_int) bool

  @[link_name="SDL_ShowSimpleMessageBox"]
  fn show_simple_message_box(flags MessageBoxFlags, title const cstring, message const cstring, window *Window) bool

  @[link_name="SDL_Metal_CreateView"]
  fn metal_create_view(window *Window) MetalView

  @[link_name="SDL_Metal_DestroyView"]
  fn metal_destroy_view(view MetalView)

  @[link_name="SDL_Metal_GetLayer"]
  fn metal_get_layer(view MetalView) rawptr

  @[link_name="SDL_OpenURL"]
  fn open_url(url const cstring) bool

  @[link_name="SDL_GetPlatform"]
  fn get_platform() const cstring

  @[link_name="SDL_CreateProcess"]
  fn create_process(args const *const cstring, pipe_stdio bool) *Process

  @[link_name="SDL_CreateProcessWithProperties"]
  fn create_process_with_properties(props PropertiesID) *Process

  @[link_name="SDL_GetProcessProperties"]
  fn get_process_properties(process *Process) PropertiesID

  @[link_name="SDL_ReadProcess"]
  fn read_process(process *Process, datasize *libc::c_size_t, exitcode *libc::c_int) rawptr

  @[link_name="SDL_GetProcessInput"]
  fn get_process_input(process *Process) *IOStream

  @[link_name="SDL_GetProcessOutput"]
  fn get_process_output(process *Process) *IOStream

  @[link_name="SDL_KillProcess"]
  fn kill_process(process *Process, force bool) bool

  @[link_name="SDL_WaitProcess"]
  fn wait_process(process *Process, block bool, exitcode *libc::c_int) bool

  @[link_name="SDL_DestroyProcess"]
  fn destroy_process(process *Process)

  @[link_name="SDL_GetNumRenderDrivers"]
  fn get_num_render_drivers() libc::c_int

  @[link_name="SDL_GetRenderDriver"]
  fn get_render_driver(index libc::c_int) const cstring

  @[link_name="SDL_CreateWindowAndRenderer"]
  fn create_window_and_renderer(title const cstring, width libc::c_int, height libc::c_int, window_flags WindowFlags, window **Window, renderer **Renderer) bool

  @[link_name="SDL_CreateRenderer"]
  fn create_renderer(window *Window, name const cstring) *Renderer

  @[link_name="SDL_CreateRendererWithProperties"]
  fn create_renderer_with_properties(props PropertiesID) *Renderer

  @[link_name="SDL_CreateSoftwareRenderer"]
  fn create_software_renderer(surface *Surface) *Renderer

  @[link_name="SDL_GetRenderer"]
  fn get_renderer(window *Window) *Renderer

  @[link_name="SDL_GetRenderWindow"]
  fn get_render_window(renderer *Renderer) *Window

  @[link_name="SDL_GetRendererName"]
  fn get_renderer_name(renderer *Renderer) const cstring

  @[link_name="SDL_GetRendererProperties"]
  fn get_renderer_properties(renderer *Renderer) PropertiesID

  @[link_name="SDL_GetRenderOutputSize"]
  fn get_render_output_size(renderer *Renderer, w *libc::c_int, h *libc::c_int) bool

  @[link_name="SDL_GetCurrentRenderOutputSize"]
  fn get_current_render_output_size(renderer *Renderer, w *libc::c_int, h *libc::c_int) bool

  @[link_name="SDL_CreateTexture"]
  fn create_texture(renderer *Renderer, format PixelFormat, access TextureAccess, w libc::c_int, h libc::c_int) *Texture

  @[link_name="SDL_CreateTextureFromSurface"]
  fn create_texture_from_surface(renderer *Renderer, surface *Surface) *Texture

  @[link_name="SDL_CreateTextureWithProperties"]
  fn create_texture_with_properties(renderer *Renderer, props PropertiesID) *Texture

  @[link_name="SDL_GetTextureProperties"]
  fn get_texture_properties(texture *Texture) PropertiesID

  @[link_name="SDL_GetRendererFromTexture"]
  fn get_renderer_from_texture(texture *Texture) *Renderer

  @[link_name="SDL_GetTextureSize"]
  fn get_texture_size(texture *Texture, w *libc::c_float, h *libc::c_float) bool

  @[link_name="SDL_SetTextureColorMod"]
  fn set_texture_color_mod(texture *Texture, r u8, g u8, b u8) bool

  @[link_name="SDL_SetTextureColorModFloat"]
  fn set_texture_color_mod_float(texture *Texture, r libc::c_float, g libc::c_float, b libc::c_float) bool

  @[link_name="SDL_GetTextureColorMod"]
  fn get_texture_color_mod(texture *Texture, r *u8, g *u8, b *u8) bool

  @[link_name="SDL_GetTextureColorModFloat"]
  fn get_texture_color_mod_float(texture *Texture, r *libc::c_float, g *libc::c_float, b *libc::c_float) bool

  @[link_name="SDL_SetTextureAlphaMod"]
  fn set_texture_alpha_mod(texture *Texture, alpha u8) bool

  @[link_name="SDL_SetTextureAlphaModFloat"]
  fn set_texture_alpha_mod_float(texture *Texture, alpha libc::c_float) bool

  @[link_name="SDL_GetTextureAlphaMod"]
  fn get_texture_alpha_mod(texture *Texture, alpha *u8) bool

  @[link_name="SDL_GetTextureAlphaModFloat"]
  fn get_texture_alpha_mod_float(texture *Texture, alpha *libc::c_float) bool

  @[link_name="SDL_SetTextureBlendMode"]
  fn set_texture_blend_mode(texture *Texture, blendMode BlendMode) bool

  @[link_name="SDL_GetTextureBlendMode"]
  fn get_texture_blend_mode(texture *Texture, blendMode *BlendMode) bool

  @[link_name="SDL_SetTextureScaleMode"]
  fn set_texture_scale_mode(texture *Texture, scaleMode ScaleMode) bool

  @[link_name="SDL_GetTextureScaleMode"]
  fn get_texture_scale_mode(texture *Texture, scaleMode *ScaleMode) bool

  @[link_name="SDL_UpdateTexture"]
  fn update_texture(texture *Texture, rect const *Rect, pixels const *Void, pitch libc::c_int) bool

  @[link_name="SDL_UpdateYUVTexture"]
  fn update_yuv_texture(texture *Texture, rect const *Rect, Yplane const *u8, Ypitch libc::c_int, Uplane const *u8, Upitch libc::c_int, Vplane const *u8, Vpitch libc::c_int) bool

  @[link_name="SDL_UpdateNVTexture"]
  fn update_nv_texture(texture *Texture, rect const *Rect, Yplane const *u8, Ypitch libc::c_int, UVplane const *u8, UVpitch libc::c_int) bool

  @[link_name="SDL_LockTexture"]
  fn lock_texture(texture *Texture, rect const *Rect, pixels *rawptr, pitch *libc::c_int) bool

  @[link_name="SDL_LockTextureToSurface"]
  fn lock_texture_to_surface(texture *Texture, rect const *Rect, surface **Surface) bool

  @[link_name="SDL_UnlockTexture"]
  fn unlock_texture(texture *Texture)

  @[link_name="SDL_SetRenderTarget"]
  fn set_render_target(renderer *Renderer, texture *Texture) bool

  @[link_name="SDL_GetRenderTarget"]
  fn get_render_target(renderer *Renderer) *Texture

  @[link_name="SDL_SetRenderLogicalPresentation"]
  fn set_render_logical_presentation(renderer *Renderer, w libc::c_int, h libc::c_int, mode RendererLogicalPresentation) bool

  @[link_name="SDL_GetRenderLogicalPresentation"]
  fn get_render_logical_presentation(renderer *Renderer, w *libc::c_int, h *libc::c_int, mode *RendererLogicalPresentation) bool

  @[link_name="SDL_GetRenderLogicalPresentationRect"]
  fn get_render_logical_presentation_rect(renderer *Renderer, rect *FRect) bool

  @[link_name="SDL_RenderCoordinatesFromWindow"]
  fn render_coordinates_from_window(renderer *Renderer, window_x libc::c_float, window_y libc::c_float, x *libc::c_float, y *libc::c_float) bool

  @[link_name="SDL_RenderCoordinatesToWindow"]
  fn render_coordinates_to_window(renderer *Renderer, x libc::c_float, y libc::c_float, window_x *libc::c_float, window_y *libc::c_float) bool

  @[link_name="SDL_ConvertEventToRenderCoordinates"]
  fn convert_event_to_render_coordinates(renderer *Renderer, event *Event) bool

  @[link_name="SDL_SetRenderViewport"]
  fn set_render_viewport(renderer *Renderer, rect const *Rect) bool

  @[link_name="SDL_GetRenderViewport"]
  fn get_render_viewport(renderer *Renderer, rect *Rect) bool

  @[link_name="SDL_RenderViewportSet"]
  fn render_viewport_set(renderer *Renderer) bool

  @[link_name="SDL_GetRenderSafeArea"]
  fn get_render_safe_area(renderer *Renderer, rect *Rect) bool

  @[link_name="SDL_SetRenderClipRect"]
  fn set_render_clip_rect(renderer *Renderer, rect const *Rect) bool

  @[link_name="SDL_GetRenderClipRect"]
  fn get_render_clip_rect(renderer *Renderer, rect *Rect) bool

  @[link_name="SDL_RenderClipEnabled"]
  fn render_clip_enabled(renderer *Renderer) bool

  @[link_name="SDL_SetRenderScale"]
  fn set_render_scale(renderer *Renderer, scaleX libc::c_float, scaleY libc::c_float) bool

  @[link_name="SDL_GetRenderScale"]
  fn get_render_scale(renderer *Renderer, scaleX *libc::c_float, scaleY *libc::c_float) bool

  @[link_name="SDL_SetRenderDrawColor"]
  fn set_render_draw_color(renderer *Renderer, r u8, g u8, b u8, a u8) bool

  @[link_name="SDL_SetRenderDrawColorFloat"]
  fn set_render_draw_color_float(renderer *Renderer, r libc::c_float, g libc::c_float, b libc::c_float, a libc::c_float) bool

  @[link_name="SDL_GetRenderDrawColor"]
  fn get_render_draw_color(renderer *Renderer, r *u8, g *u8, b *u8, a *u8) bool

  @[link_name="SDL_GetRenderDrawColorFloat"]
  fn get_render_draw_color_float(renderer *Renderer, r *libc::c_float, g *libc::c_float, b *libc::c_float, a *libc::c_float) bool

  @[link_name="SDL_SetRenderColorScale"]
  fn set_render_color_scale(renderer *Renderer, scale libc::c_float) bool

  @[link_name="SDL_GetRenderColorScale"]
  fn get_render_color_scale(renderer *Renderer, scale *libc::c_float) bool

  @[link_name="SDL_SetRenderDrawBlendMode"]
  fn set_render_draw_blend_mode(renderer *Renderer, blendMode BlendMode) bool

  @[link_name="SDL_GetRenderDrawBlendMode"]
  fn get_render_draw_blend_mode(renderer *Renderer, blendMode *BlendMode) bool

  @[link_name="SDL_RenderClear"]
  fn render_clear(renderer *Renderer) bool

  @[link_name="SDL_RenderPoint"]
  fn render_point(renderer *Renderer, x libc::c_float, y libc::c_float) bool

  @[link_name="SDL_RenderPoints"]
  fn render_points(renderer *Renderer, points const *FPoint, count libc::c_int) bool

  @[link_name="SDL_RenderLine"]
  fn render_line(renderer *Renderer, x1 libc::c_float, y1 libc::c_float, x2 libc::c_float, y2 libc::c_float) bool

  @[link_name="SDL_RenderLines"]
  fn render_lines(renderer *Renderer, points const *FPoint, count libc::c_int) bool

  @[link_name="SDL_RenderRect"]
  fn render_rect(renderer *Renderer, rect const *FRect) bool

  @[link_name="SDL_RenderRects"]
  fn render_rects(renderer *Renderer, rects const *FRect, count libc::c_int) bool

  @[link_name="SDL_RenderFillRect"]
  fn render_fill_rect(renderer *Renderer, rect const *FRect) bool

  @[link_name="SDL_RenderFillRects"]
  fn render_fill_rects(renderer *Renderer, rects const *FRect, count libc::c_int) bool

  @[link_name="SDL_RenderTexture"]
  fn render_texture(renderer *Renderer, texture *Texture, srcrect const *FRect, dstrect const *FRect) bool

  @[link_name="SDL_RenderTextureRotated"]
  fn render_texture_rotated(renderer *Renderer, texture *Texture, srcrect const *FRect, dstrect const *FRect, angle libc::c_double, center const *FPoint, flip FlipMode) bool

  @[link_name="SDL_RenderTextureAffine"]
  fn render_texture_affine(renderer *Renderer, texture *Texture, srcrect const *FRect, origin const *FPoint, right const *FPoint, down const *FPoint) bool

  @[link_name="SDL_RenderTextureTiled"]
  fn render_texture_tiled(renderer *Renderer, texture *Texture, srcrect const *FRect, scale libc::c_float, dstrect const *FRect) bool

  @[link_name="SDL_RenderTexture9Grid"]
  fn render_texture9_grid(renderer *Renderer, texture *Texture, srcrect const *FRect, left_width libc::c_float, right_width libc::c_float, top_height libc::c_float, bottom_height libc::c_float, scale libc::c_float, dstrect const *FRect) bool

  @[link_name="SDL_RenderGeometry"]
  fn render_geometry(renderer *Renderer, texture *Texture, vertices const *Vertex, num_vertices libc::c_int, indices const *libc::c_int, num_indices libc::c_int) bool

  @[link_name="SDL_RenderGeometryRaw"]
  fn render_geometry_raw(renderer *Renderer, texture *Texture, xy const *libc::c_float, xy_stride libc::c_int, color const *FColor, color_stride libc::c_int, uv const *libc::c_float, uv_stride libc::c_int, num_vertices libc::c_int, indices const *Void, num_indices libc::c_int, size_indices libc::c_int) bool

  @[link_name="SDL_RenderReadPixels"]
  fn render_read_pixels(renderer *Renderer, rect const *Rect) *Surface

  @[link_name="SDL_RenderPresent"]
  fn render_present(renderer *Renderer) bool

  @[link_name="SDL_DestroyTexture"]
  fn destroy_texture(texture *Texture)

  @[link_name="SDL_DestroyRenderer"]
  fn destroy_renderer(renderer *Renderer)

  @[link_name="SDL_FlushRenderer"]
  fn flush_renderer(renderer *Renderer) bool

  @[link_name="SDL_GetRenderMetalLayer"]
  fn get_render_metal_layer(renderer *Renderer) rawptr

  @[link_name="SDL_GetRenderMetalCommandEncoder"]
  fn get_render_metal_command_encoder(renderer *Renderer) rawptr

  @[link_name="SDL_AddVulkanRenderSemaphores"]
  fn add_vulkan_render_semaphores(renderer *Renderer, wait_stage_mask u32, wait_semaphore i64, signal_semaphore i64) bool

  @[link_name="SDL_SetRenderVSync"]
  fn set_render_v_sync(renderer *Renderer, vsync libc::c_int) bool

  @[link_name="SDL_GetRenderVSync"]
  fn get_render_v_sync(renderer *Renderer, vsync *libc::c_int) bool

  @[link_name="SDL_RenderDebugText"]
  fn render_debug_text(renderer *Renderer, x libc::c_float, y libc::c_float, str const cstring) bool

  @[link_name="SDL_RenderDebugTextFormat"]
  fn render_debug_text_format(renderer *Renderer, x libc::c_float, y libc::c_float, fmt const cstring, args @c ...) bool

  @[link_name="SDL_OpenTitleStorage"]
  fn open_title_storage(override const cstring, props PropertiesID) *Storage

  @[link_name="SDL_OpenUserStorage"]
  fn open_user_storage(org const cstring, app const cstring, props PropertiesID) *Storage

  @[link_name="SDL_OpenFileStorage"]
  fn open_file_storage(path const cstring) *Storage

  @[link_name="SDL_OpenStorage"]
  fn open_storage(iface const *StorageInterface, userdata rawptr) *Storage

  @[link_name="SDL_CloseStorage"]
  fn close_storage(storage *Storage) bool

  @[link_name="SDL_StorageReady"]
  fn storage_ready(storage *Storage) bool

  @[link_name="SDL_GetStorageFileSize"]
  fn get_storage_file_size(storage *Storage, path const cstring, length *u64) bool

  @[link_name="SDL_ReadStorageFile"]
  fn read_storage_file(storage *Storage, path const cstring, destination rawptr, length u64) bool

  @[link_name="SDL_WriteStorageFile"]
  fn write_storage_file(storage *Storage, path const cstring, source const *Void, length u64) bool

  @[link_name="SDL_CreateStorageDirectory"]
  fn create_storage_directory(storage *Storage, path const cstring) bool

  @[link_name="SDL_EnumerateStorageDirectory"]
  fn enumerate_storage_directory(storage *Storage, path const cstring, callback EnumerateDirectoryCallback, userdata rawptr) bool

  @[link_name="SDL_RemoveStoragePath"]
  fn remove_storage_path(storage *Storage, path const cstring) bool

  @[link_name="SDL_RenameStoragePath"]
  fn rename_storage_path(storage *Storage, oldpath const cstring, newpath const cstring) bool

  @[link_name="SDL_CopyStorageFile"]
  fn copy_storage_file(storage *Storage, oldpath const cstring, newpath const cstring) bool

  @[link_name="SDL_GetStoragePathInfo"]
  fn get_storage_path_info(storage *Storage, path const cstring, info *PathInfo) bool

  @[link_name="SDL_GetStorageSpaceRemaining"]
  fn get_storage_space_remaining(storage *Storage) u64

  @[link_name="SDL_GlobStorageDirectory"]
  fn glob_storage_directory(storage *Storage, path const cstring, pattern const cstring, flags GlobFlags, count *libc::c_int) *cstring

  @[link_name="SDL_SetX11EventHook"]
  fn set_x11_event_hook(callback X11EventHook, userdata rawptr)

  @[link_name="SDL_IsTablet"]
  fn is_tablet() bool

  @[link_name="SDL_IsTV"]
  fn is_tv() bool

  @[link_name="SDL_GetSandbox"]
  fn get_sandbox() Sandbox

  @[link_name="SDL_OnApplicationWillTerminate"]
  fn on_application_will_terminate()

  @[link_name="SDL_OnApplicationDidReceiveMemoryWarning"]
  fn on_application_did_receive_memory_warning()

  @[link_name="SDL_OnApplicationWillEnterBackground"]
  fn on_application_will_enter_background()

  @[link_name="SDL_OnApplicationDidEnterBackground"]
  fn on_application_did_enter_background()

  @[link_name="SDL_OnApplicationWillEnterForeground"]
  fn on_application_will_enter_foreground()

  @[link_name="SDL_OnApplicationDidEnterForeground"]
  fn on_application_did_enter_foreground()

  @[link_name="SDL_GetDateTimeLocalePreferences"]
  fn get_date_time_locale_preferences(dateFormat *DateFormat, timeFormat *TimeFormat) bool

  @[link_name="SDL_GetCurrentTime"]
  fn get_current_time(ticks *Time) bool

  @[link_name="SDL_TimeToDateTime"]
  fn time_to_date_time(ticks Time, dt *DateTime, localTime bool) bool

  @[link_name="SDL_DateTimeToTime"]
  fn date_time_to_time(dt const *DateTime, ticks *Time) bool

  @[link_name="SDL_TimeToWindows"]
  fn time_to_windows(ticks Time, dwLowDateTime *u32, dwHighDateTime *u32)

  @[link_name="SDL_TimeFromWindows"]
  fn time_from_windows(dwLowDateTime u32, dwHighDateTime u32) Time

  @[link_name="SDL_GetDaysInMonth"]
  fn get_days_in_month(year libc::c_int, month libc::c_int) libc::c_int

  @[link_name="SDL_GetDayOfYear"]
  fn get_day_of_year(year libc::c_int, month libc::c_int, day libc::c_int) libc::c_int

  @[link_name="SDL_GetDayOfWeek"]
  fn get_day_of_week(year libc::c_int, month libc::c_int, day libc::c_int) libc::c_int

  @[link_name="SDL_GetTicks"]
  fn get_ticks() u64

  @[link_name="SDL_GetTicksNS"]
  fn get_ticks_ns() u64

  @[link_name="SDL_GetPerformanceCounter"]
  fn get_performance_counter() u64

  @[link_name="SDL_GetPerformanceFrequency"]
  fn get_performance_frequency() u64

  @[link_name="SDL_Delay"]
  fn delay(ms u32)

  @[link_name="SDL_DelayNS"]
  fn delay_ns(ns u64)

  @[link_name="SDL_DelayPrecise"]
  fn delay_precise(ns u64)

  @[link_name="SDL_AddTimer"]
  fn add_timer(interval u32, callback TimerCallback, userdata rawptr) TimerID

  @[link_name="SDL_AddTimerNS"]
  fn add_timer_ns(interval u64, callback NSTimerCallback, userdata rawptr) TimerID

  @[link_name="SDL_RemoveTimer"]
  fn remove_timer(id TimerID) bool

  @[link_name="SDL_CreateTray"]
  fn create_tray(icon *Surface, tooltip const cstring) *Tray

  @[link_name="SDL_SetTrayIcon"]
  fn set_tray_icon(tray *Tray, icon *Surface)

  @[link_name="SDL_SetTrayTooltip"]
  fn set_tray_tooltip(tray *Tray, tooltip const cstring)

  @[link_name="SDL_CreateTrayMenu"]
  fn create_tray_menu(tray *Tray) *TrayMenu

  @[link_name="SDL_CreateTraySubmenu"]
  fn create_tray_submenu(entry *TrayEntry) *TrayMenu

  @[link_name="SDL_GetTrayMenu"]
  fn get_tray_menu(tray *Tray) *TrayMenu

  @[link_name="SDL_GetTraySubmenu"]
  fn get_tray_submenu(entry *TrayEntry) *TrayMenu

  @[link_name="SDL_GetTrayEntries"]
  fn get_tray_entries(menu *TrayMenu, count *libc::c_int) *const *TrayEntry

  @[link_name="SDL_RemoveTrayEntry"]
  fn remove_tray_entry(entry *TrayEntry)

  @[link_name="SDL_InsertTrayEntryAt"]
  fn insert_tray_entry_at(menu *TrayMenu, pos libc::c_int, label const cstring, flags TrayEntryFlags) *TrayEntry

  @[link_name="SDL_SetTrayEntryLabel"]
  fn set_tray_entry_label(entry *TrayEntry, label const cstring)

  @[link_name="SDL_GetTrayEntryLabel"]
  fn get_tray_entry_label(entry *TrayEntry) const cstring

  @[link_name="SDL_SetTrayEntryChecked"]
  fn set_tray_entry_checked(entry *TrayEntry, checked bool)

  @[link_name="SDL_GetTrayEntryChecked"]
  fn get_tray_entry_checked(entry *TrayEntry) bool

  @[link_name="SDL_SetTrayEntryEnabled"]
  fn set_tray_entry_enabled(entry *TrayEntry, enabled bool)

  @[link_name="SDL_GetTrayEntryEnabled"]
  fn get_tray_entry_enabled(entry *TrayEntry) bool

  @[link_name="SDL_SetTrayEntryCallback"]
  fn set_tray_entry_callback(entry *TrayEntry, callback TrayCallback, userdata rawptr)

  @[link_name="SDL_ClickTrayEntry"]
  fn click_tray_entry(entry *TrayEntry)

  @[link_name="SDL_DestroyTray"]
  fn destroy_tray(tray *Tray)

  @[link_name="SDL_GetTrayEntryParent"]
  fn get_tray_entry_parent(entry *TrayEntry) *TrayMenu

  @[link_name="SDL_GetTrayMenuParentEntry"]
  fn get_tray_menu_parent_entry(menu *TrayMenu) *TrayEntry

  @[link_name="SDL_GetTrayMenuParentTray"]
  fn get_tray_menu_parent_tray(menu *TrayMenu) *Tray

  @[link_name="SDL_UpdateTrays"]
  fn update_trays()

  @[link_name="SDL_GetVersion"]
  fn get_version() libc::c_int

  @[link_name="SDL_GetRevision"]
  fn get_revision() const cstring

  @[link_name="SDL_main"]
  fn main(argc libc::c_int, argv *cstring) libc::c_int

  @[link_name="SDL_SetMainReady"]
  fn set_main_ready()

  @[link_name="SDL_RunApp"]
  fn run_app(argc libc::c_int, argv *cstring, mainFunction main_func, reserved rawptr) libc::c_int

  @[link_name="SDL_EnterAppMainCallbacks"]
  fn enter_app_main_callbacks(argc libc::c_int, argv *cstring, appinit AppInit_func, appiter AppIterate_func, appevent AppEvent_func, appquit AppQuit_func) libc::c_int

  @[link_name="SDL_GDKSuspendComplete"]
  fn gdk_suspend_complete()

  @[link_name="SDL_Vulkan_LoadLibrary"]
  fn vulkan_load_library(path const cstring) bool

  @[link_name="SDL_Vulkan_GetVkGetInstanceProcAddr"]
  fn vulkan_get_vk_get_instance_proc_addr() FunctionPointer

  @[link_name="SDL_Vulkan_UnloadLibrary"]
  fn vulkan_unload_library()

  @[link_name="SDL_Vulkan_GetInstanceExtensions"]
  fn vulkan_get_instance_extensions(count *u32) const *const cstring

  @[link_name="SDL_Vulkan_CreateSurface"]
  fn vulkan_create_surface(window *Window, instance VkInstance, allocator const *VkAllocationCallbacks, surface *VkSurfaceKHR) bool

  @[link_name="SDL_Vulkan_DestroySurface"]
  fn vulkan_destroy_surface(instance VkInstance, surface VkSurfaceKHR, allocator const *VkAllocationCallbacks)

  @[link_name="SDL_Vulkan_GetPresentationSupport"]
  fn vulkan_get_presentation_support(instance VkInstance, physicalDevice VkPhysicalDevice, queueFamilyIndex u32) bool

  @[link_name="SDL_CreateThread"]
  fn create_thread(entry ThreadFunction, name const cstring, data rawptr) *Thread

  @[link_name="SDL_CreateThreadWithProperties"]
  fn create_thread_with_properties(props PropertiesID) *Thread
}
