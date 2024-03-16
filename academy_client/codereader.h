#ifndef CALLCODEREADER_H
#define CALLCODEREADER_H
typedef struct device_ret{
    char mac_address[20];
    char mac_device[20];
}_device_ret;
struct device_ret* get_device_code(const char *interface);
#endif
