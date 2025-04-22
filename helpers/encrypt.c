#include <stdio.h>
#include <stdlib.h>

#define XOR_KEY 0xAA
#define PADDING_SIZE (101 * 1024 * 1024) // 101 MB

int main() {
    FILE *in, *out;
    const char *inputFile = "app.exe";
    const char *outputFile = "app.enc";

    in = fopen(inputFile, "rb");
    if(!in) {
        perror("Error opening input file");
        return 1;
    }

    out = fopen(outputFile, "wb");
    if(!out) {
        perror("Error opening output file");
        fclose(in);
        return 1;
    }

    int byte;
    while ((byte = fgetc(in)) != EOF) {
        fputc(byte ^ XOR_KEY, out);
    }

    printf("Encryption of %s > %s complete. XOR key: 0x%X\n", inputFile, outputFile, XOR_KEY);

    for (size_t i = 0; i < PADDING_SIZE; i++) {
        fputc(0xAA, out);
    }

    printf("Padding added to the end of the file.\n");

    fclose(in);
    fclose(out);
    return 0;
}
