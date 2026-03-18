# Pydio Cells UI 수정 사항 및 안정화 롤백 보고서 (Superset 통합)

## 1. 개요 및 목적
- 좌측 사이드바(Left Rail)에 외부 서비스(Superset)를 내장하기 위한 전용 메뉴 추가
- 관리자 설정(Admin Console) 및 배포 시점(Helm Values)을 통한 유동적인 URL 주입 지원
- **안정성 확보**: v5-dev 버전에서 발생한 Ceph S3 연동 오류(`SignatureDoesNotMatch`) 및 컨테이너 기동 실패 문제를 해결하기 위해, 공식적으로 검증된 안정화 버전으로 롤백 후 UI 커스텀 적용

## 2. 문제 해결 및 롤백 과정 (v4.4.17 선택 배경)

### [A] S3 연동 오류 원인 분석 (`v5-dev` 이슈)
- **증상**: 커스텀 이미지를 K8s에 배포 시 S3(Ceph RGW) 관련 `SignatureDoesNotMatch`, `syncConfig not initialized yet` 에러 발생
- **원인**: 뼈대가 된 소스 코드가 최신 개발 브랜치(`v5-dev`)였음. 해당 브랜치는 `minio-go` SDK가 최신(`v7.0.99`)으로 업데이트되었고, 내부 동기화 엔진 구조가 대대적으로 변경 중인 불안정 버전이었기 때문에 Ceph와 인증 불일치 발생.

### [B] 공식 안정화 버전(v4.4.17) 탐색 및 체크아웃
- **접근**: 기존에 성공적으로 동작했던 공식 이미지(`pydio/cells:latest`)와 동일한 백엔드를 사용하기 위해 롤백 결정
- **탐색**: `git fetch https://github.com/pydio/cells.git 'refs/tags/*:refs/tags/*'` 명령을 통해 공식 리포지토리의 태그 목록 패치
- **선택**: `v4.*` 태그 중 가장 최신이자 안정화된 정식 릴리즈인 **`v4.4.17`** 태그를 찾아 `custom-ui-v4.4.17` 브랜치로 체크아웃

### [C] Dockerfile 쉘(sh) 누락 이슈 해결 (v8 -> v9)
- **증상**: v4.4.17로 빌드한 `v8` 이미지가 K8s에서 시작 직전 `exec: "/bin/sh": stat /bin/sh: no such file or directory` 에러를 내며 `CrashLoopBackOff` 발생
- **원인**: v4.4.17 당시의 `tools/docker/dockerfile` 베이스 이미지가 `FROM scratch`(빈 컨테이너)로 되어 있어, Helm Chart가 요구하는 시작 스크립트(`/bin/sh`)를 실행할 수 없었음. (v5-dev는 `alpine` 베이스를 쓰고 있었음)
- **해결**: `tools/docker/dockerfile` 베이스 이미지를 `FROM alpine`으로 수정하여 쉘 환경 제공. 이후 `v9` 이미지로 재빌드 및 푸시 완료

## 3. UI 주요 수정 내역 (v4.4.17 기반)

### [A] 설정 파라미터 연동
- **파일**: `frontend/front-srv/assets/gui.ajax/manifest.xml`
- **내용**: 
    - `SUPERSET_URL` 글로벌 파라미터 정의 (Main Options 그룹)
    - 프런트엔드 JS에서 즉시 참조 가능하도록 `expose="true"` 설정

### [B] 사이드바 메뉴 구현 (원래 UI 방식 복구)
- **파일**: `frontend/front-srv/assets/gui.ajax/res/js/ui/Workspaces/leftnav/RailPanel.js`
- **수정 상세**:
    1.  **메뉴 위치 및 노출**: Superset 메뉴를 **북마크(Bookmarks) 바로 아래**로 배치하고, `ignore` 조건을 삭제하여 URL 설정 여부와 관계없이 항상 노출되도록 변경.
    2.  **UI 방식 원복**: 화면을 덮는 풀패널(Z-index) 방식이 이질적이라는 피드백에 따라, 기존 Pydio의 부드러운 우측 슬라이딩 바 방식인 `activeBar` 렌더링 로직으로 복구.
    3.  **상태 관리 및 클릭 인터셉트**: Superset 화면이 열린 상태에서 `Home`, `All 신Files` 등 다른 메뉴를 클릭하면 즉시 기존 패널을 닫고 전환되도록 인터셉트 로직 추가.

## 4. 동작 방식 및 결과
- **설정 전**: 좌측 사이드바에 차트 아이콘 상시 노출. 클릭 시 슬라이드 패널로 "Configure `gui.ajax/SUPERSET_URL`..." 안내 문구 표시
- **설정 후**:
    - 아이콘 클릭 시 왼쪽 네비게이션 옆으로 좁은 패널 형태(`activeBar`)의 창이 열리며 Superset `iframe`이 로드됨
    - 다른 파일 메뉴와 시각적 일관성 확보
    - 메인 백엔드가 v4 안정화 버전이므로 S3 스토리지가 오류 없이 즉시 연동됨

## 5. 검증 및 설정 방법

### Helm Values 사용 (권장)
1. `values.yaml` 파일에 아래 설정 추가:
   ```yaml
   service:
     customconfigs:
       "frontend/plugin/gui.ajax/SUPERSET_URL": "https://your-superset-url.com"
   ```
2. Helm 업그레이드 배포 후 동작 확인

## 6. 최종 결과물 정보
- **최신 이미지**: `jinkernel/pydio-cells:v9`
- **빌드 특징**: 검증된 공식 백엔드(v4.4.17) + Alpine 런타임 쉘 추가 + 클릭 인터셉트가 포함된 Superset 슬라이딩 패널 UI 적용
