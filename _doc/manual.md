# User Manual & Program Guide
# คู่มือผู้ใช้และคำอธิบายโปรแกรม

## English Version
### Features
1.  **Scanner & Real-time Progress**: Recursively scans folders and shows the current file being read at the bottom of the screen.
2.  **2D Interaction (UX/UI)**:
    - **Zoom**: Use Mouse Wheel.
    - **Pan**: Click and drag the background.
    - **Drag & Bounce**: Click and drag a "star" (node). Release to watch it bounce back into the colony.
3.  **Visualization Controls**:
    - **Spacing Slider**: Adjust the separation between file clusters.
    - **Rotate Slider**: Rotate the entire file galaxy around the center.
4.  **Colony Selection**: Clicking a node highlights its "Star Colony" (related files) and fades out unrelated stars.
5.  **Search**: Auto-jump to any file by typing at least 2 characters.
6.  **Server Health**: Bottom-right indicator shows if the Go backend is alive and the current server time.

### How to Run
- Run `RUN.bat` to start the server.
- Open `http://localhost:8080`.
- Use `stop.bat` to completely kill the server and the loop.

---

## ภาษาไทย
### ฟีเจอร์หลัก
1.  **ตัวสแกนและความคืบหน้าแบบเรียลไทม์**: สแกนโฟลเดอร์แบบเรียกซ้ำและแสดงไฟล์ที่กำลังอ่านอยู่ในแถบด้านล่างของหน้าจอ
2.  **การโต้ตอบแบบ 2D (UX/UI)**:
    - **การซูม (Zoom)**: ใช้ปุ่มลูกกลิ้งเมาส์
    - **การเลื่อนหน้าจอ (Pan)**: คลิกเมาส์ค้างไว้บนพื้นหลังแล้วลาก
    - **การลากและเด้ง (Drag & Bounce)**: คลิกและลาก "ดวงดาว" (โหนด) แล้วปล่อยเพื่อดูมันเด้งกลับเข้าสู่กลุ่มดวงดาวของมันเอง
3.  **ส่วนควบคุมการแสดงภาพ**:
    - **แถบเลื่อนระยะห่าง (Spacing Slider)**: ปรับระยะห่างระหว่างกลุ่มไฟล์
    - **แถบเลื่อนการหมุน (Rotate Slider)**: หมุนกาแล็กซี่ของไฟล์ทั้งหมดรอบจุดศูนย์กลาง
4.  **การเลือกกลุ่มดวงดาว**: เมื่อคลิกที่โหนด จะเน้นเฉพาะ "กลุ่มดวงดาว" (ไฟล์ที่เกี่ยวข้องกัน) และจางดาวดวงอื่นๆ ที่ไม่เกี่ยวข้องออกไป
5.  **การค้นหา**: ระบบจะกระโดดไปยังไฟล์เป้าหมายโดยอัตโนมัติเมื่อพิมพ์ชื่อไฟล์อย่างน้อย 2 ตัวอักษร
6.  **สุขภาพของเซิร์ฟเวอร์**: ตัวบ่งชี้ที่มุมล่างขวาแสดงสถานะว่าเซิร์ฟเวอร์ Go ยังทำงานอยู่หรือไม่ พร้อมระบุเวลาของเซิร์ฟเวอร์

### วิธีการรันโปรแกรม
- รันไฟล์ `RUN.bat` เพื่อเริ่มต้นเซิร์ฟเวอร์
- เปิดบราวเซอร์ไปที่ `http://localhost:8080`
- ใช้ไฟล์ `stop.bat` เพื่อปิดเซิร์ฟเวอร์และหยุดการทำงานทั้งหมด
