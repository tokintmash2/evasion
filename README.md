# Malware

## Description

This project is an educational demonstration of how ransomware encryption/decryption works. It was created strictly for academic purposes to understand security concepts and encryption techniques in a controlled environment.

*ONLY RUN THIS CODE ON A VIRTUAL MACHINE*

### Key Generation
The program calls generateKey() to create a cryptographically secure random 256-bit key.
This key is stored in memory and will be used for both encryption and creating the decryptor.

### User Detection
The program identifies the current Windows user via getUser().
This is used to correctly locate the desktop path for file operations.

### File Encryption
The target file at C:\Users\[username]\Desktop\dont_hurt_me.victim is read into memory.
Data is padded with 100 MB to waste time for antivirus software scan (using PKCS#7 padding to ensure it's a multiple of the AES block size).
A random Initialization Vector (IV) is generated for secure encryption.
The file is encrypted using AES-256 in CBC mode.
The IV is prepended to the encrypted data.
The original file is overwritten with the encrypted version.

### Ransom Note Creation
A file named "note.txt" is created on the desktop.
The note contains a simulated ransom message along with the encryption key in hexadecimal format.
In an actual ransomware, this key would be kept private, but for educational purposes it's made visible.

### Decryptor Generation
A Go program is dynamically generated using text templates.
The encryption key is embedded in the decryptor code as a hexadecimal string.
The program is stored in the "cmd\decrypt" directory.
The Go compiler is invoked to build an executable named "decryptor.exe".
The executable is placed on the desktop for easy access.

## Setup

Set up a Windows 10 VM with Golang installed. Copy the "dont_hurt_me.victim" file to the Desktop.

### Build the malware executable

```bash
go build -o app.exe
```

### Build the encryptor executable

```bash
gcc helpers\encrypt.c -o encryptor.exe
```

Then run it:

```bash
./encryptor.exe
```

This will encypt the "app.exe" and save it as "app.enc"

### Build loader and malware executable

Find the size of the encrypted file:

```bash
go run helpers\fileSize.go app.enc
```

Copy the number and use it in helpers\loader.c:

```bash
#define ORIGINAL_PAYLOAD_SIZE <size_of_app.enc>
```

Then build the loader:

```bash
gcc helpers\loader.c -o loader.exe
```

Then use this command to create the final executable and embed the payload:

```bash
cmd /c copy /b loader.exe + app.enc <filename.exe>
```

## Usage

Run the program:

```
./filename.exe
```

This will start the ransomware simulation.

1. The program executes and triggers anti-virus.
2. The program waits for 101 seconds until anti-virus has finished scanning.
3. The program will then decrypt the encryptor and execute it in memory.
4. This will encrypt the "dont_hurt_me.victim" file on the Desktop.
5. With that, a ransom note "note.txt" is created on the Desktop, containing the key for decrytion.
6. Open the "dont_hurt_me.victim" file to verify it's encrypted.
7. The same key is used to create decryptor.exe which is also saved on the Desktop.
8. Run the decryptor.exe to decrypt the file.
9. Open the "dont_hurt_me.victim" to verify it's decrypted.
