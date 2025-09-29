# 1 Requirements
1. makdown 문서 를 RAG 에 임베딩에 최적화된 JSON 로 변환하는 command line 프로그램
2. command name 은 mtoj(linux) 또는 mtoj.exe(windows)
3. command 파라메터로 markdown 파일을 받으면 json 로 변환 후 파일로 저장
   1. output 파일은 markdown과 동일한 파일명에 .json 확장자로 저장: abc.md -> abc.json
   2. 예를들어 mtojl abc.md 로 실행시 abc.json 파일 생성
   3. 동일한 파일명의 output 파일이 이미 존재 할 경우 파일명의 postfix 로 숫자를 추가함: abc.json -> abc-1.json ->  abc-2.json
4. command 파라메터로 받은 markdown 파일이 존재하지 않으면 파일이 존재하지 않는다는 메시지를 출력하고 종료
5. command 파라메터가 누락되었거나 -h 를 전달 할 경우 다음과 같은 usage 메시지 출력
```
Markdown to JSON 1.0.0
Usage: mtoj markdown-file
```
1. markdown 에서 JSON 로 변환시 다음을 고려해서 진행
   1. markdown 의 각 라인을 JSONL 로 변환하고 JSONL ARRAY 로 구성된 JSON 으로 변환
   2. 리스트/헤딩 계층(depth·path 보존)
   3. 코드 블록을 하나의 JSONL 로 변환
