# Local student project of lab 1 (dev/ is not committed).
LAB1_DIR = dev/student

# Brings up the student project, prepares its data and runs the lab 1 autograder against it.
.PHONY: lab1
lab1:
	@test -d $(LAB1_DIR) || { echo "❌ $(LAB1_DIR) not found: create the local student project first"; exit 1; }
	@$(MAKE) -C $(LAB1_DIR) run
	@$(MAKE) -C $(LAB1_DIR) setup
	@$(MAKE) -C $(LAB1_DIR) check
