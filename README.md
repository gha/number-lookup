#number-lookup

##Installation

Downlad one of the pre-compiled binaries for your system and either extract it to a directory

    //Unix/OSX
    $ numberlookup -input /path/to/number/list
    //Windows
    $ numberlookup.exe -input C:\path\to\number\list

##PHP Script Usage

    //Ensure numberLookup is in your $PATH
    //Supply numbers to lookup via STDIN
    echo "01234567891" | numberLookup
