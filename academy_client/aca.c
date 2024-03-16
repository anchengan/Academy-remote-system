#include <stdio.h>
#include "yamlreader.h"
#include "codereader.h"
#include <string.h>
//#define DEBUGMODE
void show_copyright(void){
    printf("@copyright2024 Anchengan\n");
}
void show_help(void){
    printf("Help document for aca software:\n");
    printf("aca [authkey of remote client] [path to settings.yaml] [password of remote client]\n");
    printf("aca listenmode [path to settings.yaml]\n");
}
int is_path_in_ld_library_path(const char *path_to_check) {
    const char *ld_library_path = getenv("LD_LIBRARY_PATH");
    if (ld_library_path == NULL) {
        // 如果LD_LIBRARY_PATH未设置，则直接返回0（不在其中）
        return 0;
    }

    // 分割LD_LIBRARY_PATH并检查每个路径
    char *path_copy = strdup(ld_library_path); // 复制环境变量值
    char *token = strtok(path_copy, ":");
    while (token != NULL) {
        // 比较当前路径和要检查的路径
        if (strcmp(token, path_to_check) == 0) {
            free(path_copy); // 释放内存
            return 1; // 找到路径，返回1
        }
        token = strtok(NULL, ":"); // 继续查找下一个路径
    }

    free(path_copy); // 释放内存
    return 0; // 未找到路径，返回0
}
char* remove_char(const char* str, char c) {
    int count = 0;
    for (const char* p = str; *p; p++) {
        if (*p != c) {
            count++;
        }
    }

    // 分配足够的空间来存储新字符串
    char* result = (char*)malloc(count + 1);  // +1 为 '\0'
    if (!result) {
        return NULL;  // 内存分配失败
    }

    int i = 0;
    for (const char* p = str; *p; p++) {
        if (*p != c) {
            result[i++] = *p;
        }
    }
    result[i] = '\0';  // 确保新字符串以 '\0' 结尾

    return result;
}
int main(int argc, char *argv[]) {
    // 打印传递的参数
    if(argc-1==0){
	show_copyright();
        printf("Enter -h or --help after aca for help.\n");
	printf("Example:\naca -h\naca --help\n");
	return 0;
    }
    for(int i=1;i<=argc-1;i++){
        if(strcmp(argv[i],"-h")==0 || strcmp(argv[i],"--help")==0){
            show_copyright();
            show_help();
            return 0;
        }
    }
    if(argc-1>3){
        printf("Unknown usage of aca...\n");
        show_help();
        return 0;
    }
    int yamllength=strlen(argv[2]);
    if(yamllength<5 || argv[2][yamllength-5]!='.' || argv[2][yamllength-4]!='y' || argv[2][yamllength-3]!='a' || argv[2][yamllength-2]!='m' || argv[2][yamllength-1]!='l'){
        printf("Please check your path to settings.yaml!\n");
        printf("Example:./settings.yaml\n");
        show_help();
	return 0;
    }
    int islistenmode=0;
    if (strcmp(argv[1],"listenmode")==0){
        if(argc-1>2){
            printf("Unknown usage of aca...\n");
            show_help();
            return 0;
        }
	printf("listenmode setup successfully...\n");
	islistenmode=1;
    }else{
        printf("your authkey of remote client is:%s\n",argv[1]);
    }
    printf("setting env...\n");
    int env_return_value;
    const char *path_to_check = "./"; // 要检查的路径  
    const char *bashrc_path = getenv("HOME"); // 获取用户主目录
    if (bashrc_path == NULL) {
        perror("Error getting HOME environment variable");
        return EXIT_FAILURE;
    }

    // 拼接.bashrc文件的完整路径
    char full_path[1024];
    snprintf(full_path, sizeof(full_path), "%s/.bashrc", bashrc_path);

    // 打开.bashrc文件
    FILE *file = fopen(full_path, "r");
    if (file == NULL) {
        perror("Error opening ~/.bashrc");
        return EXIT_FAILURE;
    }

    // 要搜索的内容
    const char *search_content = "export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./";
    char line[1024]; // 假设每行不会超过1023个字符

    // 逐行读取并搜索内容
    int found = 0; // 标记是否找到内容
    while (fgets(line, sizeof(line), file)) {
        if (strstr(line, search_content) != NULL) {
#ifdef DEBUGMODE
            printf("Found the content in ~/.bashrc: %s", line);
#endif
	    found = 1; // 找到内容，设置标记
            break; // 如果只需要找到一次就退出循环
        }
    }

    

    // 关闭文件
    fclose(file);

    if (!is_path_in_ld_library_path(path_to_check) && !found) {  
          
 
        // 使用system()函数运行命令行命令
        env_return_value = system("echo 'export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./' >> ~/.bashrc");

        // 检查命令的退出码
#ifdef DEBUGMODE

        if (env_return_value == -1) {
            // system()函数调用失败
            // 在这种情况下，你可以查看errno来获取错误详情（需要包含errno.h头文件）
            perror("system() failed");
        } else if (WIFEXITED(env_return_value)) {
            // 命令正常退出，WEXITSTATUS(env_return_value)可以获取命令的退出状态
            printf("Command exited with status %d\n", WEXITSTATUS(env_return_value));
        } else if (WIFSIGNALED(env_return_value)) {
            // 命令因为接收到信号而终止，WTERMSIG(env_return_value)可以获取信号的编号
            printf("Command terminated by signal %d\n", WTERMSIG(env_return_value));
        }
#endif
        printf("env set successfully...\n");
        printf("please run: source ~/.bashrc\n");
	return 0;
    }else if(!is_path_in_ld_library_path(path_to_check) && found){
        printf("please run: source ~/.bashrc\n");
	return 0;
    }
    #include "libsender.h"
    #include "libreceiver.h"
    printf("your path to settings.yaml is:%s\n",argv[2]);    
    
    printf("Reading settings.yaml...\n");
    struct item_ret* yaml_data = get_item(argv[2]);
    struct ip* ip_reader;
    struct port* port_reader;
    ip_reader=yaml_data->ip_items;
    port_reader=yaml_data->port_items;
    char mac_device[20];

    char c = ':';   
    strcpy(mac_device,yaml_data->mac_device);
    int iplength = 0;
    struct ip* ip_pin = ip_reader;
    while (ip_pin != NULL) {
        iplength++;
        ip_pin = ip_pin->next;
    }
    
    char send_ip_info[iplength*100];
    char send_port_info[iplength*100];
    while (ip_reader != NULL){
        printf("ip:%s\n", ip_reader->ip);
	strcat(send_ip_info,ip_reader->ip);
	strcat(send_ip_info,"-");
        ip_reader = ip_reader->next;
	for(int i=0;i<2;i++){
            printf("port:%s\n", port_reader->port);
	    strcat(send_port_info,port_reader->port);
	    strcat(send_port_info,"-");
            port_reader = port_reader->next;
        }
    }
    printf("settings.yaml read completed...\n");
    if(islistenmode){
	//waiting for connection...
        printf("Getting connecting code of this device...\n");
	struct device_ret* device_data=get_device_code(mac_device);
        char* device_code=device_data->mac_device;
	char* address_code=device_data->mac_address;
        printf("Connecting code of this device is:%s\n",device_code);
        char* password = remove_char(address_code, c);
	printf("Your password is:%s\n",password);
        Receiver(device_code,password,send_ip_info,send_port_info);
    }else{
	Sender(argv[1],argv[3],send_ip_info,send_port_info);
        //connecting...
    }
    return 0; 
}
