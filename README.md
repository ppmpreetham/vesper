# Vesper                                                               
           =       *                               +        +                
         @     %    =@                  @ @          @        :@@@@          
     @+:@       @@   @@@@@   @@@          @+@@%      %@@@@      @@@          
      +@@       @@*  @@  @@@  @@@:        @@@@@ @    @@  @@@:   %@:@@@       
       @@        @@@@@   :     @@@@* :    :@@ @@--@@@@=  : @    @@%  @@      
       %@       =@@@@@@+       @@ @@@  @@@@@@  @@@  @@@@     @@@@@#@@@@@*   +
       @@@      @@     @@@@ # %@%  #     @@@@@@@@@@    @@@  =:@@@@@#:        
        @@    %@@      @@@  @@@@          @@@    #     @@@   #@@@ @*@@@      
         @@   @@@    @@@:      @@@@  @#   @@@        @@@%     @@  +:@@@      
        @@@@  @@@   *@@         =*@@@-    @@#       =@@     -@@     :@@@     
      +@= @@  @@    @@=    @@@@@@:  @@ #   @        @@@     *@@      @@      
     %    @@@@@    @@=@@ @@@@   @@@@@@@    @@      @@%@% @@@@@@     @@+      
           @@@#    @@ +@@# @      -@@@     @@*     @@  @@@ @ @@   @=         
          #:@ -     @             +@@       @@      @         %              
                    :           @*             @                =            
                                       

> Blazingly-fast, OSINT engine to track down people that fuses the raw power of [WhatsMyName](https://github.com/WebBreacher/WhatsMyName), [Sherlock](https://github.com/sherlock-project/sherlock), and [Maigret](https://github.com/soxoj/maigret) into one unified, high performance tool. Designed for speed, precision, and scale, it hunts down usernames across the internet

> [!WARNING]  
> This tool is intended for educational purposes only. Use responsibly and ethically. The author is not responsible for any misuse or illegal activities.

> [!NOTE]  
> Only WhatsMyname is stable at the moment. Sherlock and Maigret give false positives.

## Performance

| Database      | Average Execution Time | 
|---------------|------------------------|
| WhatsMyName   | ~21 seconds           |
| Sherlock      | ~12 seconds           |
| Maigret       | ~33 seconds           |

> Benchmarked on Arch Linux with AMD Ryzen 7 5000 Series CPU and 16GB RAM
> 
> ⚡ **Up to 3x faster** than competing tools like [blackbird](https://github.com/p1ngul1n0/blackbird)
  
## Installation
### Build from Source
> [!NOTE]
> Requires Go 1.20 or later
```bash
git clone https://github.com/ppmpreetham/vesper.git
cd vesper
go build -o vesper
./vesper (for Windows, use vesper.exe)
```

## Usage

### Basic Usage
```bash
# Search for a username across all databases
./vesper

# Search using a specific database
./vesper --database whatsmyname john_doe

# Search with custom timeout (10 seconds) (higher for better results, more time)
./vesper --timeout 10 john_doe
```

### Performance Tips
- **Lower timeout (3-5s)**: Faster results, may miss slow-responding sites
- **Higher timeout (10-15s)**: More thorough, catches slow sites but takes longer
- **Database selection**: Use specific databases for targeted searches
  - `whatsmyname`: Balanced coverage with good performance
  - `sherlock`: Best for general social media platforms (more false-positives)
  - `maigret`: Most comprehensive, includes niche sites (more false-positives)