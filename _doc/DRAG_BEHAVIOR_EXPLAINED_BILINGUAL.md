# Drag Behavior & Link Relationships Explained / การอธิบายพฤติกรรมการลากและการเชื่อมโยงไฟล์

## English Version

### Current Dragging Behavior
When you drag a file node, **all semantically related files move together** as a cohesive group.

### How Link Detection Works
**Direct Connections**: Files with **>75% similarity** are linked together.

**Example**:
```
File A --similar--> File B --similar--> File C
```
If you drag **File A**:
- File A moves ✅ (dragged node)
- File B moves ✅ (directly connected)
- File C does NOT move ❌ (not directly connected to A)

**Why "Far Away" Nodes Might Move**: Files sharing **high similarity scores** form **clusters** that move together.

### Semantic Similarity Based On:
- File extensions (.js, .html, .css)
- Content patterns
- Functionality relationships
- Directory structures

### Visual Indicators
- 🔵 **Dragged node**: Blue highlight
- 🟢 **Connected nodes**: Green highlight
- 📝 **Labels**: Connected node names visible

---

## ภาษาไทย

### พฤติกรรมการลากไฟล์ในปัจจุบัน
เมื่อคุณลากไฟล์ใดไฟล์หนึ่ง **ไฟล์ทั้งหมดที่เกี่ยวข้องกันทางความหมายจะเคลื่อนที่ไปด้วยกัน** เป็นกลุ่มหนึ่งเดียว

### วิธีการตรวจจับการเชื่อมโยง
**การเชื่อมโยงโดยตรง**: ไฟล์ที่มี **ความคล้ายคลึงกันมากกว่า 75%** จะถูกเชื่อมโยงเข้าด้วยกัน

**ตัวอย่าง**:
```
ไฟล์ A --คล้าย--> ไฟล์ B --คล้าย--> ไฟล์ C
```
หากคุณลาก **ไฟล์ A**:
- ไฟล์ A เคลื่อนที่ ✅ (ไฟล์ที่ถูกลาก)
- ไฟล์ B เคลื่อนที่ ✅ (เชื่อมโยงโดยตรง)
- ไฟล์ C ไม่เคลื่อนที่ ❌ (ไม่เชื่อมโยงโดยตรงกับ A)

**เหตุผลที่ไฟล์ที่อยู่ไกลกันเคลื่อนที่ด้วย**: ไฟล์ที่มี **คะแนนความคล้ายคลึงสูง** จะสร้าง **กลุ่มคลัสเตอร์** ที่เคลื่อนที่ไปพร้อมกัน

### ความคล้ายคลึงทางความหมายพิจารณาจาก:
- นามสกุลไฟล์ (.js, .html, .css)
- รูปแบบเนื้อหา
- ความสัมพันธ์ทางฟังก์ชันการทำงาน
- โครงสร้างโฟลเดอร์

### สัญลักษณ์ภาพ
- 🔵 **ไฟล์ที่ลาก**: ไฮไลต์สีฟ้า
- 🟢 **ไฟล์ที่เชื่อมโยง**: ไฮไลต์สีเขียว
- 📝 **ชื่อไฟล์**: แสดงชื่อไฟล์ที่เชื่อมโยงทั้งหมด

### คำถามที่พบบ่อย / Frequently Asked Questions

**Q: ทำไมไฟล์จากโฟลเดอร์ต่างกันถึงเคลื่อนที่ด้วยกัน?**
A: เพราะไฟล์เหล่านั้นมีความคล้ายคลึงทางความหมายสูง แม้จะอยู่ในโฟลเดอร์คนละที่

**Q: Why do files from different folders move together?**
A: They share high semantic similarity despite folder separation.

**Q: สามารถลากไฟล์เดี่ยวๆ ได้ไหม?**
A: ไม่ได้ - การออกแบบมีจุดประสงค์เพื่อสำรวจความสัมพันธ์ ใช้ฟังก์ชันค้นหา/กรองสำหรับไฟล์เดี่ยว

**Q: Can I move individual files only?**
A: No - the intent is to explore relationships. Use search/filter for individual files.

---

**Tip / เคล็ดลับ**: Use the search function (`Ctrl+F`) to find specific files without dragging clusters.
**เคล็ดลับ**: ใช้ฟังก์ชันค้นหา (`Ctrl+F`) เพื่อหาไฟล์เฉพาะโดยไม่ต้องลากกลุ่มไฟล์