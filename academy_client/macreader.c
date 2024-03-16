#include <stdio.h>  
#include <stdlib.h>  
#include <string.h>  
#include <unistd.h>  
#include "macreader.h"
#define MAX_BUF 1024
#define MAC_ADDRESS_LENGTH 18 
  
void get_mac_address(char* mac_address,const char *interface) {     
    char path[MAX_BUF];  
    FILE *fp;  
    snprintf(path, sizeof(path), "/sys/class/net/%s/address", interface);  
    
    fp = fopen(path, "r");  
    
    if (fp == NULL) {  
        perror("fopen");  
        exit(1);  
    }  
    if (fgets(mac_address, MAC_ADDRESS_LENGTH, fp) == NULL) {  
        fclose(fp);  
         exit(1);  
    }  
  
    // Remove newline character from fgets() 
    mac_address[strcspn(mac_address, "\n")] = 0;  
    fclose(fp);
    return;  
}  
  

