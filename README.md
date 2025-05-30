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
> Only WhatsMyname is supported at the moment. Sherlock and Maigret will be added in the future.

## Performance

- ~ 21 seconds for whatsmyname
- ~ 12 seconds for sherlock
- ~ 33 seconds for maigret
- Tested on Arch Linux with 16GB RAM and AMD Ryzen 7 5000 Series
- Thrice as fast as [blackbird](https://github.com/p1ngul1n0/blackbird)
  
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