Convert code page of ASCII text files
=====================================

Command line utilites to convert code page of text files

programm automatically detected code page of input files

(c) softlandia@gmail.com

dependences: 
	"github.com/softlandia/xLib"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

using  
>cpc 866 1251 .las
all files with extention "las" founded in current folder and all subfolders will be converted from IBM CodePage 866 to Windows 1251 code page
if file already in 1251 code page then nothing is done 
 

history  
0.0.1  
* init commit
