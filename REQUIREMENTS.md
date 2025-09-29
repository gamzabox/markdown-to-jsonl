# 1 Requirements
1. makdown 문서 를 RAG 에 임베딩에 최적화된 JSONL 로 변환하는 command line 프로그램
2. command name 은 mtojl(linux) 또는 mtojl.exe(windows)
3. command 파라메터로 markdown 파일을 받으면 jsonl 로 변환 후 파일로 저장
	1. output 파일은 markdown과 동일한 파일명에 .jsonl 확장자로 저장: abc.md -> abc.jsonl
	2. 예를들어 mtojl abc.md 로 실행시 abc.jsonl 파일 생성
	3. 동일한 파일명의 output 파일이 이미 존재 할 경우 파일명의 postfix 로 숫자를 추가함: abc.jsonl -> abc-1.jsonl ->  abc-2.jsonl
4. command 파라메터로 받은 markdown 파일이 존재하지 않으면 파일이 존재하지 않는다는 메시지를 출력하고 종료
5. command 파라메터가 누락되었거나 -h 를 전달 할 경우 다음과 같은 usage 메시지 출력
```
Markdown to JSONL 1.0.0
Usage: mtojl markdown-file
```
6. markdown 에서 jsonl 로 변환시 다음을 고려해서 진행
	1. 리 /헤딩 계층(depth·path 보존)
	2. 코드 블록 최적화
	3. 마크다운 내부 JSONL(워크플로 등) 이중 색인
