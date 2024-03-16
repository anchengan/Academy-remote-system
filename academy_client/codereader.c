#include "codereader.h"
#include "macreader.h"
#include "sha256tool.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#define MAX_BUF 1024
struct device_ret * get_device_code(const char *interface){
    char mac_addr[MAX_BUF];
    get_mac_address(mac_addr,interface);
    FILE *file;  
    const char* device_code_path="/tmp/academy_device_code.txt"; 
    // 打开文件，如果文件不存在则创建它
    file = fopen(device_code_path, "w");
    
    // 检查文件是否成功打开
    if (file == NULL) {
        printf("无法打开文件\n");
        exit(1);
    }
    
    // 将字符串写入文件
    fputs(mac_addr, file);
    fclose(file);
    
    char* device_code=get_sha256_sum(device_code_path); 
    struct device_ret* ret_devices=(struct device_ret*)malloc(sizeof(struct device_ret)); 
    strcpy(ret_devices-> mac_address,mac_addr);
    strcpy(ret_devices->mac_device,device_code);
    return ret_devices;
}

