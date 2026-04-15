package locker

import (
	"syscall"
	"unsafe"
)

// системная библиотека Windows
var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procGetAsyncKeyState    = user32.NewProc("GetAsyncKeyState") // Для проверки Ctrl
)

// Системные константы кодов клавиш
const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	WM_SYSKEYDOWN  = 0x0104
	VK_LWIN        = 0x5B // Левый Win
	VK_RWIN        = 0x5C // Правый Win
	VK_TAB         = 0x09 // Tab
	VK_ESCAPE      = 0x1B // Esc
	VK_CONTROL     = 0x11 // Ctrl
	LLKHF_ALTDOWN  = 0x20 // Флаг зажатого Alt
)

// Структура, которую нам передает Windows при нажатии клавиши
type kbdHookStruct struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

var hookHandle uintptr // Глобальный идентификатор хука

// Сама функция-шпион, которая решает судьбу нажатия
func keyboardHookProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		kbd := (*kbdHookStruct)(unsafe.Pointer(lParam))
		isKeyDown := wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN

		if isKeyDown {
			// 1. Блокируем клавиши Win (Пуск и системные шорткаты)
			if kbd.VkCode == VK_LWIN || kbd.VkCode == VK_RWIN {
				return 1 
			}

			// 2. Блокируем Alt + Tab (Переключение окон)
			if kbd.VkCode == VK_TAB && (kbd.Flags&LLKHF_ALTDOWN) != 0 {
				return 1
			}

			// 3. Блокируем Alt + Esc (Сворачивание/переключение окон)
			if kbd.VkCode == VK_ESCAPE && (kbd.Flags&LLKHF_ALTDOWN) != 0 {
				return 1
			}
			
			// 4. Блокируем Ctrl + Esc (Пуск) и Ctrl + Shift + Esc (Диспетчер задач)
			// Запрашиваем у системы, зажат ли сейчас Ctrl
			ctrlState, _, _ := procGetAsyncKeyState.Call(uintptr(VK_CONTROL))
			// Если старший бит равен 1 (0x8000), значит клавиша физически зажата
			isCtrlDown := (ctrlState & 0x8000) != 0

			if kbd.VkCode == VK_ESCAPE && isCtrlDown {
				return 1 // Съедаем нажатие
			}
		}
	}

	ret, _, _ := procCallNextHookEx.Call(hookHandle, uintptr(nCode), wParam, lParam)
	return ret
}

// StartKeyboardBlocker запускает хук (Вызывать строго как горутину!)
func StartKeyboardBlocker() {
	if hookHandle != 0 {
		return 
	}

	cb := syscall.NewCallback(keyboardHookProc)
	hookHandle, _, _ = procSetWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		cb,
		0, 0,
	)

	// Windows требует, чтобы поток с хуком имел цикл обработки сообщений
	var msg struct {
		hwnd    uintptr
		message uint32
		wParam  uintptr
		lParam  uintptr
		time    uint32
		pt      struct{ x, y int32 }
	}
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 || ret == ^uintptr(0) || hookHandle == 0 {
			break // Выходим, если хук сняли
		}
	}
}

// StopKeyboardBlocker снимает блокировку и отдает управление обратно
func StopKeyboardBlocker() {
	if hookHandle != 0 {
		procUnhookWindowsHookEx.Call(hookHandle)
		hookHandle = 0
	}
}