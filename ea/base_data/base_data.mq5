//+------------------------------------------------------------------+
//|                                                          tp1.mq5 |
//|                                    Copyright 2021, Dollarkiller. |
//|                           https://github.com/dollarkillerx/Quant |
//+------------------------------------------------------------------+
#property copyright "Copyright 2021, Dollarkiller."
#property link      "https://github.com/dollarkillerx/Quant"
#property version   "1.00"
#include <dollarkiller\currency.mqh>

//+------------------------------------------------------------------+
//| Expert initialization function                                   |
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

      currency_price info = Currency::GetCurrencyPrice();
      datetime st = SymbolInfoInteger(Symbol(),SYMBOL_TIME);
      FileWrite(file_handler,info.ask,info.ask_hight,info.ask_low,info.bid,info.bid_hight,info.bid_low,st);
      // Print("bin: ",last_tick.bid, "ask: ",last_tick.ask, "r: ",r," file: ",file_handler);
  }
//+------------------------------------------------------------------+
