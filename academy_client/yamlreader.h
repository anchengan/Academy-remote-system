#ifndef CALLYAMLREADER_H
#define CALLYAMLREADER_H
#include <yaml.h>

struct item_ret* get_item(const char* filename);
typedef struct item_ret{
    struct ip* ip_items;
    struct port* port_items;
    char mac_device[20];
}_item_ret;
typedef struct ip{
    char ip[20];
    struct ip* next;
}_ip;
typedef struct port{
    char port[20];
    struct port* next;
}_port;

#endif
