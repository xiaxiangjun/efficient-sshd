package system

/*
#cgo LDFLAGS: -lUserenv -static

#include <windows.h>
#include <tlhelp32.h>
#include <userenv.h>

DWORD GetLogonPid(DWORD dwSessionId, BOOL as_user)
{
    DWORD dwLogonPid = 0;
    HANDLE hSnap = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
    if (hSnap != INVALID_HANDLE_VALUE)
    {
        PROCESSENTRY32W procEntry;
        procEntry.dwSize = sizeof procEntry;

        if (Process32FirstW(hSnap, &procEntry))
            do
            {
                DWORD dwLogonSessionId = 0;
                if (_wcsicmp(procEntry.szExeFile, as_user ? L"explorer.exe" : L"winlogon.exe") == 0 &&
                    ProcessIdToSessionId(procEntry.th32ProcessID, &dwLogonSessionId) &&
                    dwLogonSessionId == dwSessionId)
                {
                    dwLogonPid = procEntry.th32ProcessID;
                    break;
                }
            } while (Process32NextW(hSnap, &procEntry));
        CloseHandle(hSnap);
    }
    return dwLogonPid;
}

BOOL LaunchProcessWin(LPWSTR cmd, DWORD logonPID)
{
	// open process
    HANDLE hProcess = OpenProcess(MAXIMUM_ALLOWED, FALSE, logonPID);
    if(NULL == hProcess)
    {
        return FALSE;
    }

    HANDLE hToken = INVALID_HANDLE_VALUE;
    BOOL ret = OpenProcessToken(hProcess, TOKEN_DUPLICATE, &hToken);
	CloseHandle(hProcess);
    if(FALSE == ret)
    {
        return FALSE;
    }

    // copy token
    HANDLE token = NULL;
    SECURITY_ATTRIBUTES sa = {0};
    sa.nLength = sizeof(sa);
    ret = DuplicateTokenEx(hToken, MAXIMUM_ALLOWED, &sa, SecurityIdentification, TokenPrimary, &token);
	CloseHandle(hToken);
    if(FALSE == ret)
	{
		return FALSE;
	}

    // open station
    STARTUPINFOW si = {0};
    PROCESS_INFORMATION pi = {0};

    si.lpDesktop = L"winsta0\\default";

    ret = CreateProcessAsUserW(token, NULL, cmd, &sa, &sa, FALSE,
        NORMAL_PRIORITY_CLASS | CREATE_NEW_CONSOLE, NULL, NULL, &si, &pi);
    CloseHandle(token);
	if(FALSE == ret)
	{
		return FALSE;
	}

	CloseHandle(pi.hThread);
	CloseHandle(pi.hProcess);
	return TRUE;
}

// ???????????????
LPCWSTR utf8_to_widechar(LPCSTR str)
{
	int len = strlen(str);
	int max = (len + 1) * sizeof(wchar_t);
	wchar_t *out = (wchar_t *) malloc(max);

	memset(out, 0, max);
	len = MultiByteToWideChar(CP_UTF8, 0, str, -1, out, max);
	return out;
}

*/
import "C"
import (
	"fmt"
	"log"
	"strings"
	"time"
	"unsafe"
)

// ????????????????????????
func LaunchProcessWithUser(name string, arg ...string) error {
	var cmd string
	args := append([]string{name}, arg...)

	for i := 0; i < len(args); i++ {
		if len(cmd) > 0 {
			cmd += "\x20"
		}

		a := strings.Trim(args[i], "\x20")
		// ????????????????????????
		if strings.HasPrefix(a, "\"") && strings.HasSuffix(a, "\"") {
			cmd += a
			continue
		}

		// ???????????????????????????
		if strings.Index(a, "\x20") > 0 || strings.Index(a, "\"") > 0 {
			cmd += "\""
			cmd += strings.ReplaceAll(a, "\"", "\\\"")
			cmd += "\""
		} else {
			cmd += a
		}
	}

	// ???????????????????????????
	go startProcess(cmd)
	return nil
}

// ???????????????????????????
func startProcess(cmd string) {
	for {
		// ????????????
		time.Sleep(time.Second * 1)

		// ???????????????ID
		sessionId := C.WTSGetActiveConsoleSessionId()
		if 0xFFFFFFFF == sessionId {
			log.Println("WTSGetActiveConsoleSessionId error")
			continue
		}

		log.Println("WTSGetActiveConsoleSessionId success: ", sessionId)
		// ?????????????????????PID
		logonPid := C.GetLogonPid(sessionId, C.FALSE)
		if 0 == logonPid {
			log.Println("GetLogonPid error")
			continue
		}

		log.Println("GetLogonPid success: ", logonPid)
		// ??????????????????
		err := launchProcess(cmd, logonPid)
		if nil == err {
			return
		}
	}
}

// ????????????
func launchProcess(cmd string, logonPid C.DWORD) error {
	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))
	wCmd := C.utf8_to_widechar(cCmd)
	defer C.free(unsafe.Pointer(wCmd))

	ret := C.LaunchProcessWin(wCmd, logonPid)
	if 0 == ret {
		return fmt.Errorf("start error")
	}

	return nil
}
