/* Forked from https://github.com/jabb3rd/rslinux_test/blob/master/loader.h */

#ifndef __LOADER_H
#define __LOADER_H

/* Some basic definitions */
#define FALSE 0
#define TRUE !FALSE

typedef int bool;
typedef unsigned char byte;
typedef unsigned short word;
typedef unsigned int dword;

/* A module description structure */
typedef struct {
        bool enabled;
        char name[512];
        char desc[1024];
} t_module_desc;

/* GetParam/SetParam parameters */
enum _st_enum {
	stEnableDebug,
	stDebugVerbosity,
	stWriteLogCallback,
	stSetTableDataCallback,
	stUserAgent,
	stUseCustomPage,
	stCustomPage,
	stDualAuthCheck,
	stPairsBasic,
	stPairsDigest,
	stProxyType,
	stProxyIP,
	stProxyPort,
	stUseCredentials,
	stCredentialsUsername,
	stCredentialsPassword,
	stPairsForm,
	stFilterRules,
	stProxyUseAuth,
	stProxyUser,
	stProxyPass
};

/* Shared library loader */
bool lib_loader(void *handle);

extern void tableDataCallback(dword row, char *name, char *value);


/* Functions' declaration */
bool Initialize();
bool GetModuleCount(dword *count);
bool GetModuleInfoW(dword index, t_module_desc *description);
bool SwitchModule(dword index, bool enabled);
bool SetParamW(dword st, void *pointer);
bool PrepareRouter(dword row, dword ip, word port, void *hrouter);
bool ScanRouter(void *hrouter);
bool FreeRouter(void *hrouter);

typedef bool (*StopRouter_t)(void *hrouter);
typedef bool (*IsRouterStopping_t)(void *hrouter);

typedef bool (*GetParam_DWord_t)(dword st, dword *value, dword size, dword *out_length);
typedef bool (*GetParam_Pointer_t)(dword st, void *pointer, dword size, dword *out_length);
typedef bool (*GetParam_Bool_t)(dword st, bool *value, dword size, dword *out_length);

typedef bool (*SetParam_Pointer_t)(dword st, void *pointer);

StopRouter_t StopRouter;
IsRouterStopping_t IsRouterStopping;

GetParam_DWord_t GetParam_DWord;
GetParam_Pointer_t GetParam_Pointer;
GetParam_Bool_t GetParam_Bool;
SetParam_Pointer_t SetParam_Pointer;


struct globalArgs_t {
    char *inputWordList; 	/* Arg -w */
    char *target_ip;		/* Arg -t */
} globalArgs;

static const char *optString = "w:t:h?";

#endif