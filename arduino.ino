#include <NewPing.h>

#define RIGHT_TRIGGER_PIN  13
#define RIGHT_ECHO_PIN     12

#define LEFT_TRIGGER_PIN  11
#define LEFT_ECHO_PIN     10

#define TOP_TRIGGER_PIN   9
#define TOP_ECHO_PIN      8

#define BACK_TRIGGER_PIN   7
#define BACK_ECHO_PIN      6

#define TOP_MAX          73
#define WIDTH_MAX        100
#define LENGTH_MAX       60

boolean passWidthBox = false;
boolean passHeightBox = false;
boolean passLengthBox = false;

int lastWidthBox;
int lastHeightBox;
int lastLengthBox;

int widthBox;
int heightBox;
int lengthBox;

int queueIndicationsWidth = 0;
int queueIndicationsHeight = 0;
int queueIndicationsLength = 0;

boolean debug = false;

NewPing rightSonar(RIGHT_TRIGGER_PIN, RIGHT_ECHO_PIN, WIDTH_MAX);
NewPing leftSonar(LEFT_TRIGGER_PIN, LEFT_ECHO_PIN, WIDTH_MAX);
NewPing topSonar(TOP_TRIGGER_PIN, TOP_ECHO_PIN, TOP_MAX);
NewPing backSonar(BACK_TRIGGER_PIN, BACK_ECHO_PIN, LENGTH_MAX);
void setup() {
  Serial.begin(9600);
}
   
void loop() {
  Indication();
  
  if(Serial.available()) {
    
    byte incomingByte = Serial.read();
    
    if(incomingByte == 0x95) {
      Serial.write(0x7F);
      Serial.write(0x7F);
      Serial.write(0x7F);
      Serial.write(0x7F);
      Serial.write(0x7F);
      Serial.flush();
    }

    // запрос ширины
    if(incomingByte == 0x88) {
      Serial.write(0x2D);
      
      if(passWidthBox) {
        Serial.write(0x7F);
      } else {
        Serial.write(0x7E);
      }
      
      Serial.write(0x0B);
      Serial.write(widthBox);
      Serial.write(0x7B);
    }
    
    // запрос высоты
    if(incomingByte == 0x99) {
      Serial.write(0x2D);
      
      if(passHeightBox) {
        Serial.write(0x7F);
      } else {
        Serial.write(0x7E);
      }      
      
      Serial.write(0x16);
      Serial.write(heightBox);
      Serial.write(0x7B);
    }

    // запрос длины
    if(incomingByte == 0x77) {
      Serial.write(0x2D);
      
      if(passLengthBox) {
        Serial.write(0x7F);
      } else {
        Serial.write(0x7E);
      }
      
      Serial.write(0x21);
      Serial.write(lengthBox);
      Serial.write(0x7B);
    } 
  }
}

int SearchAvg (int indications[], int countIndications) {
  
  int result = 0;
  int maxcount = 0;
  
  for (int i = 0; i < countIndications; i++) {
    int count = 0;
    for (int j = 0; j < countIndications; j++) {
      if (indications[i] == indications[j]) {
        count++;
      }
 
      if (maxcount < count) {
        maxcount = count;
        result = i;
      }
    }
  }

  return indications[result];
}

void Indication() {
  int countIndications = 5;
  
  int rightIndications[countIndications];
  int leftIndications[countIndications];
  int topIndications[countIndications];
  int backIndications[countIndications];


  for (int i = 0; i < countIndications; i++){
      
      rightIndications[i] = rightSonar.ping_cm();
      leftIndications[i] = leftSonar.ping_cm();
      topIndications[i] = topSonar.ping_cm();
      backIndications[i] = backSonar.ping_cm();
      
  }

  int right = SearchAvg(rightIndications, countIndications);
  int left = SearchAvg(leftIndications, countIndications);
  int top = SearchAvg(topIndications, countIndications);
  int back = SearchAvg(backIndications, countIndications);
  
  widthBox = WIDTH_MAX - (right + left);
  heightBox = TOP_MAX - top;
  lengthBox = LENGTH_MAX - back;

  passWidthBox = PassedIndication(lastWidthBox, widthBox, queueIndicationsWidth);
  passHeightBox = PassedIndication(lastHeightBox, heightBox, queueIndicationsHeight);
  passLengthBox = PassedIndication(lastLengthBox, lengthBox, queueIndicationsLength);
  
  if (debug) {
    Serial.print("Right_ping: ");
    Serial.print(right);
    Serial.println("cm");
    
    Serial.print("Left_ping: ");
    Serial.print(left);
    Serial.println("cm");
  
    Serial.print("Top_ping: ");
    Serial.print(top);
    Serial.println("cm");

    Serial.print("Back_ping: ");
    Serial.print(back);
    Serial.println("cm");
    
    Serial.println("///////////////////////////////////////");
  
    Serial.print("widthBox ");
    Serial.print(passWidthBox);
    Serial.print(" ");
    Serial.println(widthBox);
    
    Serial.print("heightBox ");
    Serial.print(passHeightBox);
    Serial.print(" ");
    Serial.println(heightBox);

    Serial.print("lengthBox ");
    Serial.print(passLengthBox);
    Serial.print(" ");
    Serial.println(lengthBox);
   
    Serial.println("///////////////////////////////////////");
    delay(2500);
  }
}

boolean PassedIndication (int &lastIndication, int &indication, int &queueIndications) {
// TODO добавить проверку на максимальное и отрицательные значения
  if (debug) {
    Serial.print("queueIndications ");
    Serial.println(queueIndications);

    Serial.print("lastIndication ");
    Serial.println(lastIndication);

    Serial.print("indication ");
    Serial.println(indication);
  }
  
    if((lastIndication - 1) <= indication && indication <= (lastIndication + 1)) {
      if(queueIndications >= 5) {   
        return true;
      } else {
        queueIndications++;
        return false;
      }
    } else {
      queueIndications = 0;
      lastIndication = indication;

      return false;
    }
}

