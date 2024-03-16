#include <stdio.h>  
#include <openssl/sha.h>  
#include <string.h>
#include <stdlib.h>
char* get_sha256_sum(const char *filename) {  
    FILE *file = fopen(filename, "rb"); // 以二进制模式打开文件  
    if (!file) {  
        perror("无法打开文件");  
        exit(1);  
    }  
   
    // 初始化SHA-256上下文  
    SHA256_CTX sha256;  
    SHA256_Init(&sha256);  
  
    // 读取文件并更新SHA-256上下文  
    const int bufSize = 4096;  
    unsigned char buffer[bufSize];  
    size_t bytesRead;  
    while ((bytesRead = fread(buffer, 1, bufSize, file)) != 0) {  
        SHA256_Update(&sha256, buffer, bytesRead);  
    }  
  
    // 计算最终的SHA-256哈希值  
    unsigned char hash[SHA256_DIGEST_LENGTH];  
    SHA256_Final(hash, &sha256);  
    int ret_size=5;
    char** ret_array = (char**)malloc(ret_size * sizeof(char*));  
    char* ret_result;
    if (ret_array == NULL) {
        printf("Memory allocation failed for string array.\n");
        exit(1);
    }
    for (int i = 0; i < ret_size; i++) {  
        ret_array[i] = (char*)malloc(50 * sizeof(char));
    	sprintf(ret_array[i],"%02x", hash[i]);
    	if(i==0){
            ret_result=ret_array[i];
    	    continue;
    	} 
    	strcat(ret_result,ret_array[i]);
    }  
    free(ret_array); 
    // 关闭文件  
    fclose(file);  
    return ret_result;
}  
  
/*int main() {  
    const char *filename = "your_file.txt"; // 替换为你要计算校验值的文件名  
    print_sha256_sum(filename);  
    return 0;  
}*/
