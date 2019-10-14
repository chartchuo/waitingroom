# waitingroom
ปัญหาของการพัฒนาระบบของบ้านเราที่เจอกันบ่อย ๆ คือความสามารถในการรับโหลดหนักๆ เปิดจองตั๋ว ลงทะเบียน หรือประกาศผล เห็นเปิดปุ๊บก็ล้นปั๊บ ผมทำ project ตัวนี้ขึ้นมาตั้งแต่ปีที่แล้ว เสร็จ phase 1 สามารถเอาไปใช้งานจริงได้แต่ไม่ค่อยมีเวลาทำต่อ และวันนี้เป็นวันดีที่จะ Open source ให้ไปต่อยอดกัน ตัวนี้ชื่อว่า waiting room เป็นโปรแกรมที่ใช้ดักอยู่หน้า web site เพื่อจัดคิวการเข้า ถ้าคนเข้ามาเยอะเกินที่ระบบจะรับไหว ก็จะให้รออยู่ในห้องรับรองรอถึงคิวก็ได้เข้าระบบ ใครนิสัยไม่ดีกด F5 refresh รัวๆ จะโดนโยนไปต่อท้ายคิว 

โปรแกรม รันได้ทั้ง standalone และ บน kubernetes ลองระดมยิงจำลองว่ามีคนเข้าแสนคนพร้อมกันยังรับไหวสบาย ๆ ยังไม่เคยยิงมากกว่านี้เพราะเครื่องต้นทางยิงไม่ไหว เครื่องที่เป็นคนยิงล่มเองซะก่อน 555

ลองเอาไปต่อยอดกันได้ที่ GitHub repo ด้านล่าง จะใส่รูปใส่เพลงให้นั่งรอใจเย็น ๆ ก็เป็นไอเดียที่ดี แต่พอดีทำไม่เป็น ไม่ถนัดเรื่อง web UI ;)

หวังว่าจะเป็นประโยชน์ ครับโผม


# Install 
install minikube
install docker
local registry

https://blog.hasura.io/sharing-a-local-registry-for-minikube-37c7240d0615/

kubectl create -f kube-registry.yml
kubectl port-forward --namespace kube-system $(kubectl get po -n kube-system | grep kube-registry-v0 | \awk '{print $1;}') 5000:5000

minikube addons enable registry
git clone https://github.com/kameshsampath/minikube-helpers


