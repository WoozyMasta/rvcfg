class CfgPatches
{
	class UtesPlusSurfaces
	{
		requiredAddons[]=
		{
			"DZ_Data"
		};
	};
};
class CfgSurfaces
{
	class DZ_SurfacesInt;
	class DZ_SurfacesExt;
	class utes_concrete: DZ_SurfacesExt
	{
		files="utes_concrete*";
		rough=0.0099999998;
		dust=0.050000001;
		friction=0.98000002;
		restitution=0.55000001;
		vpSurface="Asphalt";
		soundEnviron="road";
		soundHit="hard_ground";
		character="cp_concrete_grass";
		footDamage=0.117;
		audibility=1;
		impact="Hit_Concrete";
		deflection=0.1;
	};
};
class CfgSoundTables
{
	class CfgStepSoundTables
	{
		class BirdWalk_LookupTable
		{
			class BirdWalkstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"BirdWalk_Concrete_SoundSet"
				};
			};
		};
		class BirdGrazing_LookupTable
		{
			class BirdGrazingstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"BirdGrazing_Concrete_SoundSet"
				};
			};
		};
		class BirdBodyfall_LookupTable
		{
			class BirdBodyfallstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"BirdBodyfall_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumWalk_LookupTable
		{
			class HoofMediumWalkstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumWalk_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumRun_LookupTable
		{
			class HoofMediumRunstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumRun_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumGrazing_LookupTable
		{
			class HoofMediumGrazingstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumGrazing_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumBodyfall_LookupTable
		{
			class HoofMediumBodyfallstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumBodyfall_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumSettle_LookupTable
		{
			class HoofMediumSettlestepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumSettle_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumRest2standA_LookupTable
		{
			class HoofMediumRest2standAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumRest2standA_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumRest2standB_LookupTable
		{
			class HoofMediumRest2standBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumRest2standB_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumStand2restA_LookupTable
		{
			class HoofMediumStand2restAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumStand2restA_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumStand2restB_LookupTable
		{
			class HoofMediumStand2restBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumStand2restB_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumStand2restC_LookupTable
		{
			class HoofMediumStand2restCstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumStand2restC_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumRub1_LookupTable
		{
			class HoofMediumRub1stepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumRub1_Concrete_SoundSet"
				};
			};
		};
		class HoofMediumRub2_LookupTable
		{
			class HoofMediumRub2stepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofMediumRub2_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallWalk_LookupTable
		{
			class HoofSmallWalkstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallWalk_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallRun_LookupTable
		{
			class HoofSmallRunstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallRun_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallGrazing_LookupTable
		{
			class HoofSmallGrazingstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallGrazing_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallGrazingHard_LookupTable
		{
			class HoofSmallGrazingHardstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallGrazingHard_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallGrazingLeave_LookupTable
		{
			class HoofSmallGrazingLeavestepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallGrazingLeave_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallBodyfall_LookupTable
		{
			class HoofSmallBodyfallstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallBodyfall_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallSettle_LookupTable
		{
			class HoofSmallSettlestepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallSettle_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallRest2standA_LookupTable
		{
			class HoofSmallRest2standAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallRest2standA_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallRest2standB_LookupTable
		{
			class HoofSmallRest2standBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallRest2standB_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallStand2restA_LookupTable
		{
			class HoofSmallStand2restAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallStand2restA_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallStand2restB_LookupTable
		{
			class HoofSmallStand2restBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallStand2restB_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallStand2restC_LookupTable
		{
			class HoofSmallStand2restCstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallStand2restC_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallRub1_LookupTable
		{
			class HoofSmallRub1stepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallRub1_Concrete_SoundSet"
				};
			};
		};
		class HoofSmallRub2_LookupTable
		{
			class HoofSmallRub2stepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"HoofSmallRub2_Concrete_SoundSet"
				};
			};
		};
		class PawBigWalk_LookupTable
		{
			class PawBigWalkstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigWalk_Concrete_SoundSet"
				};
			};
		};
		class PawBigRun_LookupTable
		{
			class PawBigRunstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigRun_Concrete_SoundSet"
				};
			};
		};
		class PawBigGrazing_LookupTable
		{
			class PawBigGrazingstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigGrazing_Concrete_SoundSet"
				};
			};
		};
		class PawBigBodyfall_LookupTable
		{
			class PawBigBodyfallstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigBodyfall_Concrete_SoundSet"
				};
			};
		};
		class PawBigSettle_LookupTable
		{
			class PawBigSettlestepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigSettle_Concrete_SoundSet"
				};
			};
		};
		class PawBigRest2standA_LookupTable
		{
			class PawBigRest2standAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigRest2standA_Concrete_SoundSet"
				};
			};
		};
		class PawBigRest2standB_LookupTable
		{
			class PawBigRest2standBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigRest2standB_Concrete_SoundSet"
				};
			};
		};
		class PawBigStand2restA_LookupTable
		{
			class PawBigStand2restAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigStand2restA_Concrete_SoundSet"
				};
			};
		};
		class PawBigStand2restB_LookupTable
		{
			class PawBigStand2restBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigStand2restB_Concrete_SoundSet"
				};
			};
		};
		class PawBigStand2restC_LookupTable
		{
			class PawBigStand2restCstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigStand2restC_Concrete_SoundSet"
				};
			};
		};
		class PawBigJump_LookupTable
		{
			class PawBigJumpstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigJump_Concrete_SoundSet"
				};
			};
		};
		class PawBigImpact_LookupTable
		{
			class PawBigImpactstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawBigImpact_Concrete_SoundSet"
				};
			};
		};
		class PawMediumWalk_LookupTable
		{
			class PawMediumWalkstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumWalk_Concrete_SoundSet"
				};
			};
		};
		class PawMediumRun_LookupTable
		{
			class PawMediumRunstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumRun_Concrete_SoundSet"
				};
			};
		};
		class PawMediumGrazing_LookupTable
		{
			class PawMediumGrazingstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumGrazing_Concrete_SoundSet"
				};
			};
		};
		class PawMediumBodyfall_LookupTable
		{
			class PawMediumBodyfallstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumBodyfall_Concrete_SoundSet"
				};
			};
		};
		class PawMediumSettle_LookupTable
		{
			class PawMediumSettlestepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumSettle_Concrete_SoundSet"
				};
			};
		};
		class PawMediumRest2standA_LookupTable
		{
			class PawMediumRest2standAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumRest2standA_Concrete_SoundSet"
				};
			};
		};
		class PawMediumRest2standB_LookupTable
		{
			class PawMediumRest2standBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumRest2standB_Concrete_SoundSet"
				};
			};
		};
		class PawMediumStand2restA_LookupTable
		{
			class PawMediumStand2restAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumStand2restA_Concrete_SoundSet"
				};
			};
		};
		class PawMediumStand2restB_LookupTable
		{
			class PawMediumStand2restBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumStand2restB_Concrete_SoundSet"
				};
			};
		};
		class PawMediumStand2restC_LookupTable
		{
			class PawMediumStand2restCstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumStand2restC_Concrete_SoundSet"
				};
			};
		};
		class PawMediumJump_LookupTable
		{
			class PawMediumJumpstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumJump_Concrete_SoundSet"
				};
			};
		};
		class PawMediumImpact_LookupTable
		{
			class PawMediumImpactstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawMediumImpact_Concrete_SoundSet"
				};
			};
		};
		class PawSmallWalk_LookupTable
		{
			class PawSmallWalkstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallWalk_Concrete_SoundSet"
				};
			};
		};
		class PawSmallRun_LookupTable
		{
			class PawSmallRunstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallRun_Concrete_SoundSet"
				};
			};
		};
		class PawSmallGrazing_LookupTable
		{
			class PawSmallGrazingstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallGrazing_Concrete_SoundSet"
				};
			};
		};
		class PawSmallBodyfall_LookupTable
		{
			class PawSmallBodyfallstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallBodyfall_Concrete_SoundSet"
				};
			};
		};
		class PawSmallSettle_LookupTable
		{
			class PawSmallSettlestepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallSettle_Concrete_SoundSet"
				};
			};
		};
		class PawSmallRest2standA_LookupTable
		{
			class PawSmallRest2standAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallRest2standA_Concrete_SoundSet"
				};
			};
		};
		class PawSmallRest2standB_LookupTable
		{
			class PawSmallRest2standBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallRest2standB_Concrete_SoundSet"
				};
			};
		};
		class PawSmallStand2restA_LookupTable
		{
			class PawSmallStand2restAstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallStand2restA_Concrete_SoundSet"
				};
			};
		};
		class PawSmallStand2restB_LookupTable
		{
			class PawSmallStand2restBstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallStand2restB_Concrete_SoundSet"
				};
			};
		};
		class PawSmallStand2restC_LookupTable
		{
			class PawSmallStand2restCstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallStand2restC_Concrete_SoundSet"
				};
			};
		};
		class PawSmallJump_LookupTable
		{
			class PawSmallJumpstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallJump_Concrete_SoundSet"
				};
			};
		};
		class PawSmallImpact_LookupTable
		{
			class PawSmallImpactstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"PawSmallImpact_Concrete_SoundSet"
				};
			};
		};
		class bodyfall_Zmb_LookupTable
		{
			class bodyfall_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_asphalt_ext_Zmb_SoundSet"
				};
			};
		};
		class bodyfall_hand_Zmb_LookupTable
		{
			class bodyfall_hand_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_hand_asphalt_ext_Zmb_SoundSet"
				};
			};
		};
		class bodyfall_slide_Zmb_LookupTable
		{
			class bodyfall_slide_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_slide_asphalt_ext_Zmb_SoundSet"
				};
			};
		};
		class walkErc_Bare_Zmb_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class runErc_Bare_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class sprintErc_Bare_Zmb_LookupTable
		{
			class sprintErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"sprintErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class landFeetErc_Bare_Zmb_LookupTable
		{
			class landFeetErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"landFeetErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class scuffErc_Bare_Zmb_LookupTable
		{
			class scuffErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"scuffErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class walkRasErc_Bare_Zmb_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class runRasErc_Bare_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class landFootErc_Bare_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class walkCro_Bare_Zmb_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class runCro_Bare_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class jumpErc_Bare_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_bare_Zmb_Soundset"
				};
			};
		};
		class walkErc_Boots_Zmb_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class runErc_Boots_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class sprintErc_Boots_Zmb_LookupTable
		{
			class sprintErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"sprintErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class landFeetErc_Boots_Zmb_LookupTable
		{
			class landFeetErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"landFeetErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class scuffErc_Boots_Zmb_LookupTable
		{
			class scuffErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"scuffErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class walkRasErc_Boots_Zmb_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class runRasErc_Boots_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class landFootErc_Boots_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class walkCro_Boots_Zmb_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class runCro_Boots_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class jumpErc_Boots_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_boots_Zmb_Soundset"
				};
			};
		};
		class walkErc_Sneakers_Zmb_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class runErc_Sneakers_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class sprintErc_Sneakers_Zmb_LookupTable
		{
			class sprintErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"sprintErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class landFeetErc_Sneakers_Zmb_LookupTable
		{
			class landFeetErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"landFeetErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class scuffErc_Sneakers_Zmb_LookupTable
		{
			class scuffErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"scuffErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class walkRasErc_Sneakers_Zmb_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class runRasErc_Sneakers_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class landFootErc_Sneakers_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class walkCro_Sneakers_Zmb_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class runCro_Sneakers_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class jumpErc_Sneakers_Zmb_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_sneakers_Zmb_Soundset"
				};
			};
		};
		class walkProne_Zmb_LookupTable
		{
			class walkProne_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkProne_concrete_ext_Zmb_Soundset"
				};
			};
		};
		class runProne_Zmb_LookupTable
		{
			class runProne_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runProne_concrete_ext_Zmb_Soundset"
				};
			};
		};
		class walkErc_Char_LookupTable
		{
			class walkErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkErc_concrete_ext_Char_Soundset"
				};
			};
		};
		class walkRasErc_Char_LookupTable
		{
			class walkRasErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkRasErc_concrete_ext_Char_Soundset"
				};
			};
		};
		class runErc_Char_LookupTable
		{
			class runErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runErc_concrete_ext_Char_Soundset"
				};
			};
		};
		class runRasErc_Char_LookupTable
		{
			class runRasErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runRasErc_concrete_ext_Char_Soundset"
				};
			};
		};
		class sprintErc_Char_LookupTable
		{
			class sprintErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"sprintErc_concrete_ext_Char_Soundset"
				};
			};
		};
		class landFootErc_Char_LookupTable
		{
			class landFootErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"landFootErc_concrete_ext_Char_Soundset"
				};
			};
		};
		class landFeetErc_Char_LookupTable
		{
			class landFeetErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"landFeetErc_concrete_ext_Char_Soundset"
				};
			};
		};
		class scuffErc_Char_LookupTable
		{
			class scuffErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"scuffErc_concrete_ext_Char_Soundset"
				};
			};
		};
		class walkCro_Char_LookupTable
		{
			class walkCro_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkCro_concrete_ext_Char_Soundset"
				};
			};
		};
		class runCro_Char_LookupTable
		{
			class runCro_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runCro_concrete_ext_Char_Soundset"
				};
			};
		};
		class jumpErc_Char_LookupTable
		{
			class jumpErc_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"jumpErc_concrete_ext_Char_Soundset"
				};
			};
		};
		class walkProne_Char_LookupTable
		{
			class walkProne_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkProne_concrete_ext_Char_Soundset"
				};
			};
		};
		class runProne_Char_LookupTable
		{
			class runProne_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"runProne_concrete_ext_Char_Soundset"
				};
			};
		};
		class walkProne_noHS_Char_LookupTable
		{
			class walkProne_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkProne_noHS_asphalt_ext_Char_Soundset"
				};
			};
		};
		class walkProneLong_noHS_Char_LookupTable
		{
			class walkProneLong_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"walkProneLong_noHS_asphalt_ext_Char_Soundset"
				};
			};
		};
		class handstepSound_Char_LookupTable
		{
			class handstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"Handstep_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class handstepSound_Hard_Char_LookupTable
		{
			class handstepSound_Hard_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"Handstep_Hard_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class handsstepSound_Char_LookupTable
		{
			class handsstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"Handsstep_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class bodyfallSound_Char_LookupTable
		{
			class bodyfallSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class bodyfall_handSound_Char_LookupTable
		{
			class bodyfall_handSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_hand_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class bodyfall_rollSound_Char_LookupTable
		{
			class bodyfall_rollSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_roll_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class bodyfall_rollHardSound_Char_LookupTable
		{
			class bodyfall_rollHardSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_rollHard_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class bodyfall_slideSound_Char_LookupTable
		{
			class bodyfall_slideSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_slide_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class bodyfall_slide_lightSound_Char_LookupTable
		{
			class bodyfall_slide_lightSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_slide_light_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class bodyfall_hand_lightSound_Char_LookupTable
		{
			class bodyfall_hand_lightSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"bodyfall_hand_light_asphalt_ext_Char_SoundSet"
				};
			};
		};
		class step_ladder_Char_LookupTable
		{
			class handstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"step_ladder_Char_Soundset"
				};
			};
		};
		class step_ladder_run_Char_LookupTable
		{
			class handstepSound_utes_concrete
			{
				surface="utes_concrete";
				soundSets[]=
				{
					"step_ladder_run_Char_Soundset"
				};
			};
		};
	};
};
