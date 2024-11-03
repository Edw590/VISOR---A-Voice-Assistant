// Got it from https://superuser.com/a/1436051/719714.

#include <tchar.h>
#include <Windows.h>

#ifndef __CRT_STRINGIZE
#define __CRT_STRINGIZE(Value) #Value
#endif
#ifndef _CRT_STRINGIZE
#define _CRT_STRINGIZE(Value) __CRT_STRINGIZE(Value)
#endif

enum { WLAN_NOTIFICATION_SOURCE_ACM = 0x00000008 };
typedef enum _WLAN_NOTIFICATION_ACM { wlan_notification_acm_start, wlan_notification_acm_autoconf_enabled, wlan_notification_acm_autoconf_disabled, wlan_notification_acm_background_scan_enabled, wlan_notification_acm_background_scan_disabled, wlan_notification_acm_bss_type_change, wlan_notification_acm_power_setting_change, wlan_notification_acm_scan_complete, wlan_notification_acm_scan_fail, wlan_notification_acm_connection_start, wlan_notification_acm_connection_complete, wlan_notification_acm_connection_attempt_fail, wlan_notification_acm_filter_list_change, wlan_notification_acm_interface_arrival, wlan_notification_acm_interface_removal, wlan_notification_acm_profile_change, wlan_notification_acm_profile_name_change, wlan_notification_acm_profiles_exhausted, wlan_notification_acm_network_not_available, wlan_notification_acm_network_available, wlan_notification_acm_disconnecting, wlan_notification_acm_disconnected, wlan_notification_acm_adhoc_network_state_change, wlan_notification_acm_profile_unblocked, wlan_notification_acm_screen_power_change, wlan_notification_acm_profile_blocked, wlan_notification_acm_scan_list_refresh, wlan_notification_acm_end } WLAN_NOTIFICATION_ACM, *PWLAN_NOTIFICATION_ACM;
typedef enum _WLAN_INTERFACE_STATE { wlan_interface_state_not_ready, wlan_interface_state_connected, wlan_interface_state_ad_hoc_network_formed, wlan_interface_state_disconnecting, wlan_interface_state_disconnected, wlan_interface_state_associating, wlan_interface_state_discovering, wlan_interface_state_authenticating } WLAN_INTERFACE_STATE, *PWLAN_INTERFACE_STATE;
typedef struct _WLAN_INTERFACE_INFO { GUID InterfaceGuid; WCHAR strInterfaceDescription[256]; WLAN_INTERFACE_STATE isState; } WLAN_INTERFACE_INFO;
typedef struct _WLAN_INTERFACE_INFO_LIST { DWORD dwNumberOfItems; DWORD dwIndex; WLAN_INTERFACE_INFO InterfaceInfo[1]; } WLAN_INTERFACE_INFO_LIST;
typedef struct _WLAN_NOTIFICATION_DATA { DWORD NotificationSource; DWORD NotificationCode; GUID InterfaceGuid; DWORD dwDataSize; void *pData; } WLAN_NOTIFICATION_DATA, *PWLAN_NOTIFICATION_DATA;
typedef void WINAPI WLAN_NOTIFICATION_CALLBACK(WLAN_NOTIFICATION_DATA *, void *);

typedef struct wlan_scan_finished_context { WLAN_INTERFACE_INFO_LIST *interface_list; HANDLE semaphore; } wlan_scan_finished_context;
static void WINAPI wlan_notification_callback(WLAN_NOTIFICATION_DATA *data, void *context)
{
    if ((data->NotificationSource | WLAN_NOTIFICATION_SOURCE_ACM) == WLAN_NOTIFICATION_SOURCE_ACM)
    {
        if (data->NotificationCode == wlan_notification_acm_power_setting_change || data->NotificationCode == wlan_notification_acm_scan_complete)
        {
            wlan_scan_finished_context *const ctx = (wlan_scan_finished_context *)context;
            if (ctx)
            {
                for (unsigned int i = 0; i != (ctx->interface_list ? ctx->interface_list->dwNumberOfItems : 0); ++i)
                {
                    if (memcmp(&ctx->interface_list->InterfaceInfo[i].InterfaceGuid, &data->InterfaceGuid, sizeof(data->InterfaceGuid)) == 0)
                    {
                        ctx->interface_list->InterfaceInfo[i].isState = (WLAN_INTERFACE_STATE)(ctx->interface_list->InterfaceInfo[i].isState | (data->NotificationCode << 16));
                    }
                }
                if (ctx->semaphore) { long prev; ReleaseSemaphore(ctx->semaphore, 1, &prev); }
            }
        }
    }
}

int _tmain(int argc, TCHAR *argv[])
{
    (void)argv;
    unsigned int result;
    if (argc > 1) { result = ERROR_INVALID_PARAMETER; }
    else
    {
        HMODULE const wlanapi = LoadLibrary(TEXT("wlanapi.dll"));
        if (wlanapi)
        {
#define X(Module, Return, Name, Params) typedef Return Name##_t Params; Name##_t *Name = (Name##_t *)GetProcAddress(Module, __CRT_STRINGIZE(Name));
            X(wlanapi, DWORD WINAPI, WlanCloseHandle, (HANDLE hClientHandle, void *pReserved));
            X(wlanapi, DWORD WINAPI, WlanEnumInterfaces, (HANDLE hClientHandle, void *pReserved, struct _WLAN_INTERFACE_INFO_LIST **ppInterfaceList));
            X(wlanapi, void WINAPI, WlanFreeMemory, (void *pMemory));
            X(wlanapi, DWORD WINAPI, WlanOpenHandle, (DWORD dwClientVersion, void *pReserved, DWORD *pdwNegotiatedVersion, HANDLE *phClientHandle));
            X(wlanapi, DWORD WINAPI, WlanRegisterNotification, (HANDLE hClientHandle, DWORD dwNotifSource, BOOL bIgnoreDuplicate, WLAN_NOTIFICATION_CALLBACK funcCallback, void *pCallbackContext, void *pReserved, DWORD *pdwPrevNotifSource));
            X(wlanapi, DWORD WINAPI, WlanScan, (HANDLE hClientHandle, const GUID *pInterfaceGuid, const struct _DOT11_SSID *pDot11Ssid, const struct _WLAN_RAW_DATA *pIeData, void *pReserved));
#undef X
            DWORD version;
            HANDLE handle = NULL;
            result = WlanOpenHandle ? WlanOpenHandle(1, NULL, &version, &handle) : ERROR_PROC_NOT_FOUND;
            if (result == ERROR_SUCCESS)
            {
                WLAN_INTERFACE_INFO_LIST *interface_list = NULL;
                result = WlanEnumInterfaces ? WlanEnumInterfaces(handle, NULL, &interface_list) : ERROR_PROC_NOT_FOUND;
                if (result == ERROR_SUCCESS)
                {
                    wlan_scan_finished_context context = { interface_list, CreateSemaphore(NULL, 0, (LONG)interface_list->dwNumberOfItems, NULL) };
                    DWORD prev;
                    unsigned int nwait = 0;
                    unsigned long const register_notification_result = WlanRegisterNotification ? WlanRegisterNotification(handle, WLAN_NOTIFICATION_SOURCE_ACM, FALSE, wlan_notification_callback, &context, NULL, &prev) : ERROR_PROC_NOT_FOUND;
                    for (unsigned int i = 0; i != interface_list->dwNumberOfItems; ++i)
                    {
                        unsigned int const result_i = WlanScan ? WlanScan(handle, &interface_list->InterfaceInfo[i].InterfaceGuid, NULL, NULL, NULL) : ERROR_PROC_NOT_FOUND;
                        interface_list->InterfaceInfo[i].isState = (WLAN_INTERFACE_STATE)(interface_list->InterfaceInfo[i].isState & 0xFFFF);
                        if (result_i == ERROR_SUCCESS) { ++nwait; }
                        else { interface_list->InterfaceInfo[i].isState = (WLAN_INTERFACE_STATE)(interface_list->InterfaceInfo[i].isState | (result_i << 16)); }
                    }
                    while (nwait > 0)
                    {
                        WaitForSingleObject(context.semaphore, INFINITE);
                        --nwait;
                    }
                    for (unsigned int i = 0; i != interface_list->dwNumberOfItems; ++i)
                    {
                        unsigned int const result_i = interface_list->InterfaceInfo[i].isState >> 16;
                        if (register_notification_result == ERROR_SUCCESS && result_i != ERROR_SUCCESS)
                        { result = result_i; }
                    }
                    if (result == ERROR_SUCCESS && register_notification_result != ERROR_SUCCESS) { result = ERROR_IO_PENDING; }
                    if (context.semaphore) { CloseHandle(context.semaphore); }
                    WlanFreeMemory ? WlanFreeMemory(interface_list) : ERROR_PROC_NOT_FOUND;
                }
                WlanCloseHandle ? WlanCloseHandle(handle, NULL) : ERROR_PROC_NOT_FOUND;
            }
        }
        else { result = (unsigned int)GetLastError(); }
    }
    return (int)result;
}