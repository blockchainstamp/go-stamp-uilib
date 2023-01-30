#ifndef _CALL_BACK_HEADER_
#define _CALL_BACK_HEADER_

        typedef char* (*UserInterfaceAPI) (const char*);
        typedef void (*SetLastErr) (const char*);

        char* bridge_func(UserInterfaceAPI f, const char* v);
        void bridge_Error(SetLastErr f,  const char* v);

#endif