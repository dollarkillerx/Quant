//+------------------------------------------------------------------+
//|                                                    base_data.mq5 |
//|                                    Copyright 2021, Dollarkiller. |
//|                           https://github.com/dollarkillerx/Quant |
//+------------------------------------------------------------------+
#property copyright "Copyright 2021, Dollarkiller."
#property link      "https://github.com/dollarkillerx/Quant"
#property version   "1.00"
//+------------------------------------------------------------------+
//| 获取基础数据 使用 GO OR PYTHON做数据分析
//+------------------------------------------------------------------+

int file_handler;
int OnInit()
  {

   string terminal_data_path=TerminalInfoString(TERMINAL_DATA_PATH);
   string filename2 =terminal_data_path+"\\MQL5\\Files\\"+"fractals.csv";
   string filename = "fractals.csv";
   file_handler=FileOpen(filename,FILE_WRITE|FILE_CSV);
   if (file_handler <= 0) {
      Print("File Error: ", GetLastError());
   }

   printf("r: %s f: %d",filename2,file_handler);
//---
   return(INIT_SUCCEEDED);
  }
//+------------------------------------------------------------------+
//| Expert deinitialization function                                 |
//+------------------------------------------------------------------+
void OnDeinit(const int reason)
  {
//---
   printf("close");
   FileFlush(file_handler);
   FileClose(file_handler);
  }
//+------------------------------------------------------------------+
//| Expert tick function                                             |
//+------------------------------------------------------------------+
void OnTick()
  {
      MqlTick last_tick;
      SymbolInfoTick(Symbol(),last_tick);

      int r = FileWrite(file_handler,last_tick.bid,last_tick.ask);
      // Print("bin: ",last_tick.bid, "ask: ",last_tick.ask, "r: ",r," file: ",file_handler);
  }
//+------------------------------------------------------------------+
