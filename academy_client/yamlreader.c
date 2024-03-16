#include <stdio.h>
#include <yaml.h>
#include <string.h>
#include "yamlreader.h"
enum itemtypes{
    INIT_ITEM,
    IP,
    PORT,
    MAC
};

struct item_ret* get_item(const char* filename) {
    FILE* file = fopen(filename, "r");
    if (file == NULL) {
        printf("Failed to open the YAML file.\n");
        exit(1);
    }
    yaml_parser_t parser;
    yaml_event_t event;
    // Initialize the parser and set input file
    if (!yaml_parser_initialize(&parser)) {
        printf("Failed to initialize the YAML parser.\n");
        exit(1);
    }
    enum itemtypes typeofitem=INIT_ITEM;
    yaml_parser_set_input_file(&parser, file);
    struct ip* oldip;
    struct ip* newip=NULL;
    struct port* oldport;
    struct port* newport=NULL;
    int Loop=1;
    char mac_device[20];
    do {
        // Parse events from the stream until we reach the end of the doc
        if (!yaml_parser_parse(&parser, &event)) {
           
           break;
        }
        	
        switch (event.type) {
	 
            case YAML_SCALAR_EVENT:
	       if(strcmp(event.data.scalar.value,"--address")==0){
		    typeofitem=IP;
		    break;
		}
	       if(strcmp(event.data.scalar.value,"--addressend")==0){
                    typeofitem=INIT_ITEM;
                    break;
                }
		
		if(strcmp(event.data.scalar.value,"--port")==0){
                    typeofitem=PORT;
                    break;
                }
		if(strcmp(event.data.scalar.value,"--portend")==0){
                    typeofitem=INIT_ITEM;
                    break;
                }
		if(strcmp(event.data.scalar.value,"-macdevice")==0){
                    typeofitem=MAC;
                    break;
                }
		if(strcmp(event.data.scalar.value,"-macdeviceend")==0){
                    typeofitem=INIT_ITEM;
                    break;
                }
		if(strcmp(event.data.scalar.value,"-end")==0){
                    Loop=0;
                    break;
                }
         	if (typeofitem==IP){
		   oldip=newip;
                   newip = (struct ip*)malloc(sizeof(struct ip));
                   if (newip == NULL) {  
                       printf("Memory allocation failed.\n");  
                       exit(1);  
		   }  
                   strcpy(newip->ip , event.data.scalar.value);  
                   newip->next = oldip;  
		   break;
		};
                
		if (typeofitem==PORT){
                    oldport=newport;
                    newport = (struct port*)malloc(sizeof(struct port));
                    if (newport == NULL) {
                        printf("Memory allocation failed.\n");
                        exit(1);
                    }
                    strcpy(newport->port , event.data.scalar.value);
                    newport->next = oldport;
                    break;
                };
		if (typeofitem==MAC){
                    strcpy(mac_device, event.data.scalar.value);
                    break;
                };
            default:
                break;
        }
        
        // Free memory for the current event
        yaml_event_delete(&event);
    } while (Loop);
    // Clean up resources
    struct item_ret* ret_items=(struct item_ret*)malloc(sizeof(struct item_ret));
    
    ret_items->ip_items=newip;
    ret_items->port_items=newport;
    strcpy(ret_items->mac_device,mac_device);
    yaml_parser_delete(&parser);
    return ret_items;
}
