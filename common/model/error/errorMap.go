package cerror

// Error Code Convention
// aabbccdd
// aa - Service Folder. e.g curriculam
// bb - Sub Service Folder. e.g dbSyncService
// cc - File. e.g. service.go
// dd - error counter. e.g. 22 for 22nd error in perticular file
//
// Directory wise code map
//  /
//  |
//	|---analytics 10
//  |	|---collectorService 01
//	|	|	|---service.go 01 - 100101
//	|	|	|---logging.go 02 - 100102
//	|	|	|---instrumenting.go 03 - 100103
//	|	|	|---endpoint.go 04 - 100104
//	|	|	|---transport_http.go 05 - 100105
//	|	|---config 02
//	|	|	|---config.go 01 - 100201
//	|	|---domain 03
//	|	|	|---dailyAttempt.go 01 - 100301
//	|	|	|---lastActivity.go 02 - 100302
//	|	|	|---weeklyAttempt.go 03 - 100303
//	|---apiGateway 20
//  |	|---authorizationService 01
//	|	|	|---auth.go 01 - 200101
//	|	|	|---transport_http.go 02 - 200102
//  |	|---userSubscriptionService 02
//	|	|	|---auth.go 01 - 200201
//	|	|	|---transport_http.go 02 - 200202
//  |	|---churchService 03
//	|	|	|---auth.go 01 - 200301
//	|	|	|---transport_http.go 02 - 200302
//  |	|---contentDeliveryService 04
//	|	|	|---auth.go 01 - 200401
//	|	|	|---transport_http.go 02 - 200402
//  |	|---contentUploaderService 05
//	|	|	|---auth.go 01 - 200501
//	|	|	|---transport_http.go 02 - 200502
//  |	|---soloSongService 06
//	|	|	|---auth.go 01 - 200601
//	|	|	|---transport_http.go 02 - 200602
//  |	|---songService 07
//	|	|	|---auth.go 01 - 200701
//	|	|	|---transport_http.go 02 - 200702
//  |	|---systemUserService 08
//	|	|	|---auth.go 01 - 200801
//	|	|	|---transport_http.go 02 - 200802
//  |	|---teacherReviewService 09
//	|	|	|---auth.go 01 - 200901
//	|	|	|---transport_http.go 02 - 200902
//  |	|---userChurchService 10
//	|	|	|---auth.go 01 - 201001
//	|	|	|---transport_http.go 02 - 201002
//  |	|---userSettingService 11
//	|	|	|---auth.go 01 - 201101
//	|	|	|---transport_http.go 02 - 201102
//  |	|---warmUpService 12
//	|	|	|---auth.go 01 - 201201
//	|	|	|---transport_http.go 02 - 201202
//  |	|---appExpiryService 13
//	|	|	|---auth.go 01 - 201301
//	|	|	|---transport_http.go 02 - 201302
//	|---appExpiry 30
//  |	|---appExpiryService 01
//	|	|	|---service.go 01 - 300101
//	|	|	|---logging.go 02 - 300102
//	|	|	|---instrumenting.go 03 - 300103
//	|	|	|---endpoint.go 04 - 300104
//	|	|	|---transport_http.go 05 - 300105
//	|	|---config 02
//	|	|	|---config.go 01 - 300201
//	|	|---domain 03
//	|	|	|---appExpiry.go 01 - 300301
//	|---auth 40
//  |	|---authService 01
//	|	|	|---service.go 01 - 400101
//	|	|	|---logging.go 02 - 400102
//	|	|	|---instrumenting.go 03 - 400103
//	|	|	|---endpoint.go 04 - 400104
//	|	|	|---transport_http.go 05 - 400105
//	|	|---config 02
//	|	|	|---config.go 01 - 400201
//	|	|---domain 03
//	|	|	|---role.go 01 - 400301
//	|	|	|---token.go 02 - 400302
//  |	|---authService 04
//	|	|	|---service.go 01 - 400401
//	|	|	|---logging.go 02 - 400402
//	|	|	|---instrumenting.go 03 - 400403
//	|	|	|---endpoint.go 04 - 400404
//	|	|	|---transport_http.go 05 - 400405
//	|---common 50
//  |	|---apiSigner 01
//	|	|	|---signer.go 01 - 500101
//  |	|---archiver 02
//	|	|	|---zip.go 01 - 500201
//  |	|---downloader 03
//	|	|	|---downloader.go 01 - 500301
//  |	|---middleware 04
//	|	|	|---middleware.go 01 - 500401
//  |	|---model 05
//	|	|	|---pan.go 01 - 500501
//  |	|---panDecriptor 06
//	|	|	|---metadataGenerator.go 01 - 500601
//  |	|---remover 07
//	|	|	|---remover.go 01 - 500701
//  |	|---s3UrlHelper 08
//	|	|	|---s3UrlHelper.go 01 - 500801
//  |	|---transport 09
//	|	|	|---notification.go 01 - 500901
//  |	|---uploader 10
//	|	|	|---uploader.go 01 - 501001
//  |	|---uploadVerifier 11
//	|	|	|---s3UploadVerifier.go 01 - 501101
//	|---notifier 60
//  |	|---notificationService 01
//	|	|	|---service.go 01 - 600101
//	|	|	|---logging.go 02 - 600102
//	|	|	|---instrumenting.go 03 - 600103
//	|	|	|---endpoint.go 04 - 600104
//	|	|	|---transport_http.go 05 - 600105
//	|	|---config 02
//	|	|	|---config.go 01 - 600201
//	|	|---domain 03
//	|	|	|---banner.go 01 - 600301
//	|	|	|---notification.go 02 - 600302
//	|	|	|---userSubscription.go 03 - 600303
//	|	|---notifier 04
//	|	|	|---pushNotification.go 01 - 600401
//	|	|---renderer 05
//	|	|	|---renderer.go 01 - 600501
//  |	|---userSubscriptionService 07
//	|	|	|---service.go 01 - 600701
//	|	|	|---logging.go 02 - 600702
//	|	|	|---instrumenting.go 03 - 600703
//	|	|	|---endpoint.go 04 - 600704
//	|	|	|---transport_http.go 05 - 600705
//	|---user 70
//  |	|---userService 01
//	|	|	|---service.go 01 - 700101
//	|	|	|---logging.go 02 - 700102
//	|	|	|---instrumenting.go 03 - 700103
//	|	|	|---endpoint.go 04 - 700104
//	|	|	|---transport_http.go 05 - 700105
//	|	|---config 02
//	|	|	|---config.go 01 - 700201
//	|	|---domain 03
//	|	|	|---attemptHistory.go 01 - 700301
//	|	|	|---attemptUpload.go 02 - 700302
//	|	|	|---church.go 03 - 700303
//	|	|	|---churchSong.go 04 - 700304
//	|	|	|---invite.go 05 - 700305
//	|	|	|---passwordReset.go 06 - 700306
//	|	|	|---settings.go 07 - 700307
//	|	|	|---soloSong.go 08 - 700308
//	|	|	|---song.go 09 - 700309
//	|	|	|---systemUser.go 10 - 700310
//	|	|	|---teacherReview.go 11 - 700311
//	|	|	|---teacherReviewVersion.go 12 - 700312
//	|	|	|---userChurch.go 13 - 700313
//	|	|	|---warmup.go 14 - 700314
//  |	|---userService 04
//	|	|	|---service.go 01 - 700401
//	|	|	|---logging.go 02 - 700402
//	|	|	|---instrumenting.go 03 - 700403
//	|	|	|---endpoint.go 04 - 700404
//	|	|	|---transport_http.go 05 - 700405
//  |	|---contentDeliveryService 05
//	|	|	|---service.go 01 - 700501
//	|	|	|---logging.go 02 - 700502
//	|	|	|---instrumenting.go 03 - 700503
//	|	|	|---endpoint.go 04 - 700504
//	|	|	|---transport_http.go 05 - 700505
//  |	|---contentUploaderService 06
//	|	|	|---service.go 01 - 700601
//	|	|	|---logging.go 02 - 700602
//	|	|	|---instrumenting.go 03 - 700603
//	|	|	|---endpoint.go 04 - 700604
//	|	|	|---transport_http.go 05 - 700605
//  |	|---filestore 07
//	|	|	|---store.go 01 - 700701
//  |	|---soloSongService 08
//	|	|	|---service.go 01 - 700801
//	|	|	|---logging.go 02 - 700802
//	|	|	|---instrumenting.go 03 - 700803
//	|	|	|---endpoint.go 04 - 700804
//	|	|	|---transport_http.go 05 - 700805
//  |	|---songService 09
//	|	|	|---service.go 01 - 700901
//	|	|	|---logging.go 02 - 700902
//	|	|	|---instrumenting.go 03 - 700903
//	|	|	|---endpoint.go 04 - 700904
//	|	|	|---transport_http.go 05 - 700905
//  |	|---systemUserService 10
//	|	|	|---service.go 01 - 701001
//	|	|	|---logging.go 02 - 701002
//	|	|	|---instrumenting.go 03 - 701003
//	|	|	|---endpoint.go 04 - 701004
//	|	|	|---transport_http.go 05 - 701005
//  |	|---teacherReviewService 11
//	|	|	|---service.go 01 - 701101
//	|	|	|---logging.go 02 - 701102
//	|	|	|---instrumenting.go 03 - 701103
//	|	|	|---endpoint.go 04 - 701104
//	|	|	|---transport_http.go 05 - 701105
//  |	|---userChurchService 12
//	|	|	|---service.go 01 - 701201
//	|	|	|---logging.go 02 - 701202
//	|	|	|---instrumenting.go 03 - 701203
//	|	|	|---endpoint.go 04 - 701204
//	|	|	|---transport_http.go 05 - 701205
//  |	|---userSettingService 13
//	|	|	|---service.go 01 - 701301
//	|	|	|---logging.go 02 - 701302
//	|	|	|---instrumenting.go 03 - 701303
//	|	|	|---endpoint.go 04 - 701304
//	|	|	|---transport_http.go 05 - 701305
//  |	|---warmupService 14
//	|	|	|---service.go 01 - 701401
//	|	|	|---logging.go 02 - 701402
//	|	|	|---instrumenting.go 03 - 701403
//	|	|	|---endpoint.go 04 - 701404
//	|	|	|---transport_http.go 05 - 701405
