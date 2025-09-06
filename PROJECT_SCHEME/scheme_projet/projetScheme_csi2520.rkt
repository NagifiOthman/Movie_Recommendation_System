#lang scheme
; AUTEUR :Othman Nagifi  300347989

; lecture de csv

(define (read-f filename) (call-with-input-file filename
  (lambda (input-port)
   (let loop ((line (read-line input-port)))
   (cond 
    ((eof-object? line) '())
    (#t (begin (cons (string-split line ",") (loop (read-line input-port))))))))))

; conversion du csv en liked/notLiked
(define (convert-rating L) (list (string->number (car L)) (string->number (cadr L)) (< 3.5 (string->number (caddr L)))))

; Permet de définir la liste Ratings

(define Ratings (map convert-rating (read-f "test.csv")))



; Fonction : add-rating
; Paramètres : rating et la liste-des-utilisateurs
; Sortie : une liste d'utilisateurs avec la note ajoutée
; Description : prend une note et une liste d'utilisateurs et ajoute la note à la liste des utilisateurs
(define (add-rating rating list-of-users)
  (if (null? list-of-users) ; cas1 : la liste-des-utilisateurs est vide, besoin de créer un nouvel utilisateur
  (if (caddr rating)
      (list (list (car rating) (list (cadr rating)) '()))
      (list (list (car rating) '() (list (cadr rating)))))
  (let ((user-id (caar list-of-users)) (likedMovies (cadar list-of-users)) (notLikedMovies (caddar list-of-users))) ; définition de variables locales
  (cond
    ((and (= (car rating) user-id) (caddr rating))  ; cas2 : l'utilisateur existe et aime le film
      (cons (list user-id (cons (cadr rating) likedMovies) notLikedMovies) (cdr list-of-users))) ; ajouter le film à la liste des films aimés
      
    ((and (= (car rating) user-id) (not (caddr rating)))  ; cas3 : l'utilisateur existe et n'aime pas le film
      (cons (list user-id likedMovies (cons (cadr rating) notLikedMovies)) (cdr list-of-users))) ; ajouter le film à la liste des films non aimés
    (else
     (cons (car list-of-users) (add-rating rating (cdr list-of-users))))))))

; Fonction : add-ratings
; Paramètres : la liste-des-notes et la liste-des-utilisateurs
; Sortie : une liste d'utilisateurs avec les notes ajoutées
; Description : prend une liste de notes et une liste d'utilisateurs et ajoute les notes à la liste des utilisateurs en utilisant la fonction add-rating
(define (add-ratings list-of-ratings list-of-users)
  (if (null? list-of-ratings) ; cas de base : la liste des notes est vide
  list-of-users
  (add-ratings (cdr list-of-ratings) (add-rating (car list-of-ratings) list-of-users)))) ; continuer à ajouter les notes à la liste des utilisateurs


; Fonction : Users
; Paramètres : la fonction add-ratings et la liste-des-notes
; Sortie : une liste d'utilisateurs avec les notes ajoutées
; Description : sauvegarde tous les utilisateurs dans une liste en utilisant la liste des notes
(define Users (add-ratings Ratings '())) ; créer la base de données des utilisateurs

; Fonction : get-user
; Paramètres : ID et la liste-des-utilisateurs
; Sortie : l'utilisateur ayant l'ID donné
; Description : prend un ID et une liste d'utilisateurs et retourne l'utilisateur avec l'ID donné
(define (get-user ID list-of-users)
  (if (null? list-of-users) '()
  (if (= ID (caar list-of-users)) ; trouvé l'utilisateur correspondant à l'ID
       (car list-of-users)
   (get-user ID (cdr list-of-users))))) ; continuer la recherche

; Fonction : similarity
; Paramètres : deux utilisateurs
; Sortie : le score de similarité entre les deux utilisateurs
; Description : prend deux utilisateurs et retourne la similarité entre eux
(define (similarity user1 user2)
  (let
  ((user1Liked (cadr (get-user user1 Users))) ; variables locales des préférences de l'utilisateur
   (user1Disliked (caddr (get-user user1 Users)))
   (user2Liked (cadr (get-user user2 Users)))
   (user2Disliked (caddr (get-user user2 Users))))

    ; calculer la similarité
    (exact->inexact
     (/ (+ (length (intersection user1Liked user2Liked)) (length (intersection user1Disliked user2Disliked)))
    (length (union (union (union user1Liked user2Liked) user1Disliked) user2Disliked))))))

; Fonction : intersection
; Paramètres : deux listes
; Sortie : la liste d'intersection des deux listes
; Description : prend deux listes et retourne l'intersection des deux listes
(define (intersection user1List user2List)
  (filter (lambda (movie) (member movie user2List)) user1List))

  ; Fonction : difference
  ; Paramètres : deux listes
  ; Sortie : la liste des éléments présents dans la première liste mais pas dans la seconde
  ; Description : prend deux listes et retourne la différence entre elles
  (define (difference list1 list2)
    (filter (lambda (movie) (not (member movie list2))) list1))
; Fonction : union
; Paramètres : deux listes
; Sortie : la liste union des deux listes
; Description : prend deux listes et retourne l'union des deux listes
(define (union user1List user2List)
  (cond ((and (null? user1List) (null? user2List) '())) ; vérifier si les deux listes sont vides
    ((null? user1List) user2List) ; cas de base qui sera probablement atteint
    ((null? user2List) user1List)
    (else
     (if (member (car user1List) user2List)
     (union (cdr user1List) user2List) ; si l'élément de user1 est dans la liste user2, alors le retirer
     (cons (car user1List) (union (cdr user1List) user2List)))))) ; l'ajouter s'il n'est pas dans la liste user2 avec la liste user2
