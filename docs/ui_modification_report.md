# Pydio Cells UI 수정 사항 보고서

## 1. 수정 목적
- 좌측 사이드바(Left Rail)에 외부 서비스(Superset)를 내장하기 위한 메뉴 추가
- 관리자 설정을 통해 URL을 유동적으로 변경 가능하도록 구현

## 2. 주요 수정 내역

### [A] 설정 파라미터 추가
- **파일**: `frontend/assets/gui.ajax/manifest.xml`
- **내용**: 
    - `SUPERSET_URL` 글로벌 파라미터 정의
    - 관리자 콘솔(Main Options)에서 입력 가능하도록 노출
    - 프런트엔드 JS에서 접근 가능하도록 `expose="true"` 설정

### [B] 사이드바 메뉴 및 패널 구현
- **파일**: `frontend/assets/gui.ajax/res/js/ui/Workspaces/leftnav/RailPanel.js`
- **수정 상세**:
    1.  **`ExternalPanel` 컴포넌트 추가**: 
        - `iframe`을 사용하여 외부 URL 로드
        - 별도 창에서 열기(`Open`) 버튼 포함
        - URL 미설정 시 안내 문구 출력
    2.  **`toolbars` 배열 확장**:
        - `id: 'superset'` 항목 추가
        - `chart-box-outline` 아이콘 적용
        - `SUPERSET_URL` 값이 없을 경우 메뉴 숨김 처리 (`ignore` 필드)
    3.  **상태 관리**:
        - 클릭 시 `activePanel`을 'superset'으로 변경하여 화면 전환

## 3. 예상 동작 및 결과
- **설정 전**: 좌측 사이드바에 변화 없음 (기존과 동일)
- **설정 후** (관리자 콘솔에서 URL 입력 시):
    - 좌측 메뉴 상단(Home 아래)에 차트 아이콘 메뉴 출현
    - 메뉴 클릭 시 메인 콘텐츠 영역 대신 Superset 페이지가 포함된 `iframe` 패널이 우측에 열림
    - 상단 `Open` 버튼 클릭 시 새 탭으로 해당 URL 연결

## 4. 검증 방법
1.  **관리자 로그인** -> `Settings` -> `Application Parameters` -> `Cells Main Options` 이동
2.  **Superset URL** 항목에 `https://superset.example.com` 입력 후 저장
3.  **새침한 로드** 후 좌측 사이드바에 메뉴 아이콘 생성 확인
4.  아이콘 클릭 시 내부 `iframe`으로 페이지가 정상적으로 임베딩되는지 확인
