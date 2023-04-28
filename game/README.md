Mafia: 

Ranks: {"Rookie", "Delivery boy", "Henchman", "Soldier", "Thief", "Hustler", "Advisor", "Lieutenant", "Don", "Capo di tutti capi"}

!config:
	start_cash <int>:		default(10000) #starting cash
	main_channel <int>:		default(null) #channelid if you dont want it to spam other channels
	active <bool>:			default(true) #activate or deactivate bot
	protected <int>:		default(7) #days

!profile: shows information about yourself
	Health 					int 	#Health information in %
	Cash 					int64 	#amount of cash
	Bank 					float64	#amount of cash in bank
	Safe					float64	#amount of cash in safe
	Rank/progress 			float32	#rank progress
	Ammo 					int64	#amount of bullets
	Strength 				int64	#strength 
	Speed 					int64	#speed
	Attack Strength 		int64	#attach strengt combined with gangsters, weapons, bullets and personal weapons
	Defense Strength 		int64	#same as attach + headquarter
	headquarter: 			string	#wich headquarted do you have, and do you have space for your gangsters
		gangsters 			int64	#amount of gangsters
		bed-rooms 			int64	#amount of bedrooms(upgradeable by 100%)

!hospital
	heal <full|%>: cost 1000 pr. 10%

!work [hours] default 1 hour
	Rookie: 			1,200 pr hour
	Delivery boy: 		2,400
	Henchman: 			4,800
	Soldier: 			9,600
	Thief:				19,200
	Hustler:			38,400
	Advisor: 			76,800
	Lieutenant: 		153,600
	Don: 				307,200
	Capo: 				614,400

!bank (need rank of Soldier)
	withdraw <int amount>: take out cash
	deposit <int amount>: put in cash
	amount: checks the balance

!safe (need rank of Thief)
	withdraw <int amount>: take out cash
	deposit <int amount>: put in cash
	amount: checks the balance

!crime
	overview: shows the chances of each robbery
	robchild: steal candy from a child
	robbeggar: steal from a beggar
	pickpocket: steal from pedestrians
	robcar: steal a car
	robatm: steal from a atm
	robstore: steal from a store
	kidnap: kidnap and extortion

!recruit <string location> [int times]
	list: shows where you can recruite gangsters and the price {"street", "bar", "prison", "nightclub", "strip club", "pawnshop", "gun store", "drug den"}

!gym <string type> [int times]
	workout: gives strength
	mill: gives speed
	tennis: gives small strength and speed
	basket: gives small strength and speed

!court
	status: shows if you are in jail or not, if true show price
	pay: pay the fine

!factory
	buy <string weapon|bullet>: buys a weapons or bullet factory
	upgrade <string weapon|bullet>: upgrade the factory
	status: shows how many weapons/bullets are beeing made + price for upgrade

!shop
	buy <string weapon|house> [int amount]: buy a weapon
	weapon [string weapon]: list weapons [weapon stats] {"knife", "coltM911", "uzi", "revolver", "desert eagle", "mp5", "spas-12", "ak47"}
	house [string house]: list houses [house stats] {"shack", "apartment", "condo", "house", "mansion", "villa", "castle"}

!attack
	<string user> <int gangsters> <int weapons> <int bullets>
	#cant attack somone 2 ranks lower
	#cant attack somone on vacation
	#attacker will get the fortune
	#attacker will NOT get cash in the safe
	#only 1 attack pr. day

!death
	#show list of when someone was murdered

!vacation <int days>
	#price 10,000 pr day
	#puts yourself the specified number of days on vacation
	#protected and cant be attacked
	#cant unvacation yourself