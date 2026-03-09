# Conceptual Framework: 26-Dimensional Vector Space
# กรอบแนวคิดเชิงทฤษฎี: ปริภูมิเวกเตอร์ 26 มิติ

## <a name="english-version"></a>English Version
This project treats the filesystem as a high-dimensional data universe. Every file and folder is mapped into a **26-Dimensional (26D) Vector Space** to find hidden relationships and similarities beyond just folder nesting.

### 1. Feature Extraction (Embedding)
Instead of using a pre-trained neural network, we use a deterministic **Hand-Crafted Feature Extractor**. We normalize 26 distinct file attributes into a numerical "profile":
- **Identity (8D):** Size (Log10), Size Tail (Mod 1000), Folder Depth, Name Length, Extension Hash, and Character Patterns (dots, underscores, digits).
- **Time (3D):** Hour of Day, Weekday, and Month of modification.
- **Prefix Signature (5D):** Normalized first 5 characters of the name.
- **Content Hash (10D):** The first 10 bytes of the SHA-256 hash, ensuring content-level similarity is represented.

### 2. Relationship Mapping (Cosine Similarity)
To find how "close" two stars (files) are, we calculate the **Cosine Similarity** between their vectors:
\[ 	ext{Similarity} = \frac{\mathbf{A} \cdot \mathbf{B}}{\|\mathbf{A}\| \|\mathbf{B}\|} \]
This focuses on the "profile" (direction) rather than just raw size or time.

### 3. Force-Directed Visualization
We apply a **Physics-Based Simulation** (Force-Directed Graph) where:
- Nodes with high similarity **attract** each other.
- All nodes **repel** each other to avoid overlap.
- **Damping (Shake & Brake)** is used to cool down the energy and reach a stable state quickly.

---

## <a name="ภาษาไทย"></a>ภาษาไทย
โปรเจกต์นี้มองระบบไฟล์เป็นจักรวาลข้อมูลที่มีมิติสูง โดยไฟล์และโฟลเดอร์แต่ละรายการจะถูกจับคู่เข้ากับ **ปริภูมิเวกเตอร์ 26 มิติ (26D Vector Space)** เพื่อค้นหาความสัมพันธ์และความคล้ายคลึงที่ซ่อนอยู่นอกเหนือจากการจัดวางในโฟลเดอร์แบบปกติ

### 1. การสกัดคุณลักษณะ (Embedding)
เราไม่ได้ใช้โครงข่ายประสาทเทียมที่ฝึกฝนมาล่วงหน้า แต่เราใช้ **ตัวสกัดคุณลักษณะที่ออกแบบด้วยมือ (Hand-Crafted Feature Extractor)** โดยการปรับค่าแอตทริบิวต์ของไฟล์ 26 รายการให้เป็นค่ามาตรฐาน (Normalize):
- **ข้อมูลพื้นฐาน (8 มิติ):** ขนาดไฟล์ (Log10), ส่วนท้ายของขนาดไฟล์ (Mod 1000), ความลึกของโฟลเดอร์, ความยาวชื่อไฟล์, แฮชของนามสกุลไฟล์ และรูปแบบตัวอักษร (จุด, ขีดล่าง, ตัวเลข)
- **เวลา (3 มิติ):** ชั่วโมง, วันในสัปดาห์ และเดือนที่แก้ไขไฟล์
- **ลายเซ็นชื่อไฟล์ (5 มิติ):** ค่ามาตรฐานของตัวอักษร 5 ตัวแรกของชื่อไฟล์
- **แฮชของเนื้อหา (10 มิติ):** 10 ไบต์แรกของ SHA-256 แฮช เพื่อให้ความคล้ายคลึงในระดับเนื้อหาถูกนำมาคำนวณด้วย

### 2. การสร้างความสัมพันธ์ (ความคล้ายคลึงโคไซน์)
เพื่อหาว่า "ดวงดาว" (ไฟล์) สองดวงอยู่ใกล้กันแค่ไหน เราคำนวณ **ความคล้ายคลึงโคไซน์ (Cosine Similarity)** ระหว่างเวกเตอร์:
\[ 	ext{Similarity} = \frac{\mathbf{A} \cdot \mathbf{B}}{\|\mathbf{A}\| \|\mathbf{B}\|} \]
วิธีนี้จะเน้นที่ "โปรไฟล์" (ทิศทางของเวกเตอร์) มากกว่าแค่ขนาดหรือเวลาดิบๆ

### 3. การแสดงผลด้วยแรงฟิสิกส์ (Force-Directed Visualization)
เราใช้ **การจำลองทางฟิสิกส์ (Force-Directed Graph)** โดยที่:
- โหนดที่มีความคล้ายคลึงสูงจะ **ดึงดูด** กัน
- โหนดทั้งหมดจะ **ผลัก** กันเองเพื่อไม่ให้ซ้อนทับกัน
- **การลดแรงสั่นสะเทือน (Shake & Brake)** ถูกนำมาใช้เพื่อลดพลังงานส่วนเกินและทำให้ระบบเข้าสู่สถานะคงที่ (Stable) ได้อย่างรวดเร็ว
