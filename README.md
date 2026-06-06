# Dokumentacja realizacji zadań

## 1. Konfiguracja środowiska i uwierzytelnianie

Poniżej przedstawiono proces generowania tokena PAT na platformie GitHub oraz uwierzytelnianie lokalnego środowiska za pomocą narzędzia `gh` CLI.
<img width="940" height="294" alt="image" src="https://github.com/user-attachments/assets/c59541f0-9f26-4f73-844e-5b3924118ece" />
<img width="940" height="157" alt="image" src="https://github.com/user-attachments/assets/2428f4d7-f927-4f51-a8c7-444246a542f6" />

## 2. Pliki konfiguracyjne i struktura repozytorium

W celu wykluczenia zbędnych plików tymczasowych oraz zabezpieczenia danych wrażliwych, utworzono dedykowany plik `.gitignore`. Utworzono również plik definicji łańcucha CI.
<img width="940" height="346" alt="image" src="https://github.com/user-attachments/assets/6b9f77b3-f695-47cd-bb47-c24c3b7713db" />

<img width="940" height="823" alt="image" src="https://github.com/user-attachments/assets/f5599dc6-cc2f-4ac5-96a4-b9b46e49c0d9" />

## 3. Utworzenie repozytorium zdalnego i konfiguracja sekretów

Repozytorium zdalne zostało utworzone bezpośrednio z terminala. Następnie skonfigurowano zmienne i sekrety wymagane do logowania w usłudze DockerHub.
<img width="940" height="124" alt="image" src="https://github.com/user-attachments/assets/fd3a9011-08d2-4276-9f1d-988a965e5a3c" />
<img width="940" height="50" alt="image" src="https://github.com/user-attachments/assets/08bf09f9-5208-4eb5-9aaf-8d608566ffc1" />
<img width="854" height="152" alt="image" src="https://github.com/user-attachments/assets/d13f91cc-9813-4fe7-8a61-8e2814806d8a" />

## 4. Konfiguracja kluczy SSH dla repozytorium z kodem źródłowym

Aby umożliwić agentowi SSH wewnątrz potoku GitHub Actions bezpieczne pobranie kodu aplikacji, wygenerowano dedykowaną parę kluczy.
<img width="940" height="98" alt="image" src="https://github.com/user-attachments/assets/69490b94-2506-48f7-b6b2-df16414d0d07" />
<img width="940" height="251" alt="image" src="https://github.com/user-attachments/assets/6fc45dec-59cd-40dd-85fe-8feed9c986a1" />

## 5. Synchronizacja kodu i wyzwolenie potoku CI/CD

Pliki zostały dodane do indeksu, zatwierdzone i wypchnięte na serwer. Potok został wyzwolony poprzez nadanie tagu wersji.
<img width="940" height="186" alt="image" src="https://github.com/user-attachments/assets/f742670e-ccb2-482e-98cd-2d7e18116634" />
<img width="916" height="327" alt="image" src="https://github.com/user-attachments/assets/73298d8d-d281-4416-9ceb-5df301c573fa" />
<img width="940" height="202" alt="image" src="https://github.com/user-attachments/assets/447afa85-7a40-42e3-9dd1-e0141d11dbd0" />

## 6. Weryfikacja działania potoku GitHub Actions

Monitorowanie uruchomionego łańcucha CI potwierdziło poprawne przejście wszystkich zaplanowanych kroków: logowania, konfiguracji SSH, budowania, testu Trivy oraz publikacji do rejestru GHCR.
<img width="940" height="704" alt="image" src="https://github.com/user-attachments/assets/fe7d50f1-ba27-4208-8087-482779961ce9" />
<img width="940" height="289" alt="image" src="https://github.com/user-attachments/assets/f1c69618-4833-4208-88bf-c6ab0609956f" />

### Uzasadnienie rozwiązań technicznych

#### Przyjęty sposób tagowania obrazów i danych cache
1. **Obrazy (GHCR):** Tagowanie jest realizowane przez moduł `metadata-action` zgodnie ze schematami `sha` oraz `semver` . Zastosowano hierarchię priorytetów (SemVer wyższy) oraz parametr `latest=false`, wyłączając automatyczne generowanie tagu `latest`. Gwarantuje to jednoznaczność wersji wdrożeniowych.
2. **Pamięć cache (DockerHub):** Zastosowano mechanizm `type=registry`. Zdecydowano się na użycie oddzielnego, stałego tagu `:cache` (`pawcho-cache:cache`). Tagi pamięci podręcznej wskazują na odrębne manifesty w architekturze rejestru OCI.

#### Analiza zachowania cache
**Co się dzieje z cachem, kiedy usuwamy stare obrazy z repozytorium?**
Usunięcie docelowego obrazu aplikacyjnego z GHCR nie wpływa na pamięć podręczną. Dane cache są przechowywane na niezależnym koncie (DockerHub) pod odrębnym adresem (`pawcho-cache:cache`). Rejestry OCI traktują manifesty cache i warstwy obrazów wynikowych jako oddzielne encje, dzięki czemu kolejne budowanie aplikacji wykorzysta zachowany na DockerHubie stan cache nawet po wyczyszczeniu GitHub Packages.

#### Testy CVE (Trivy)
Do weryfikacji obrazu użyto skanera Trivy. Jest to optymalne rozwiązanie w łańcuchach CI, ponieważ pozwala na zatrzymanie operacji (`exit-code: '1'`) dla zdefiniowanych poziomów zagrożenia (`CRITICAL, HIGH`). Skan odbywa się na tymczasowo zbudowanym obrazie z parametrem `load: true` (przed krokiem `push`). Warunkowa dyrektywa `if: success()` w kolejnym kroku gwarantuje, że do publicznego repozytorium nie trafi obraz ze zidentyfikowanymi lukami.
